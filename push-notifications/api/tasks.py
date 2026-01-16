import datetime

import pytz
from sqlalchemy import Date, and_, cast, func, select

from celery import shared_task

from .internal.db import db
from .internal.models.campaigns import ScheduledCampaigns
from .internal.services.gorush import GorushException, create_push_notification
from .internal.services.smart_gas import (
    change_status_campaign,
    get_campaigns,
    get_devices_to_notify,
)

timezone = "America/Mazatlan"


def get_localized_date():
    local_tz = pytz.timezone(timezone)
    return datetime.datetime.now(local_tz)


def local_to_utc(local_date: datetime.datetime):
    local_tz = pytz.timezone(timezone)
    ma = local_tz.localize(local_date)

    return ma.astimezone(pytz.UTC)


@shared_task(ignore_result=False)
def notify(external_campaign_id, title, message, extra_data):
    result = get_devices_to_notify(external_campaign_id)
    if result.status_code != 200:
        return {"status_code": result.status_code}

    data = result.json()
    tokens = []
    devices = {}
    for device in data:
        device_token = device.get("TokenM")
        notification_id = device.get("Cve_Id")
        tokens.append(device_token)
        devices[device_token] = notification_id

    try:
        r = create_push_notification(
            data={
                "notifications": [
                    {
                        "tokens": tokens,
                        "title": title,
                        "message": message,
                        "data": extra_data,
                        "platform": 2,
                    }
                ]
            }
        )
        return r
    except GorushException:
        return {"error": "There is an exception when notifying"}


@shared_task(ignore_result=False)
def check_notifications():
    campaings_scheduled = 0

    result = get_campaigns()

    if result.status_code != 200:
        return {"campaings": 0, "status_code": result.status_code}

    campaings = result.json()

    campaing_ids = [c["Cve_Campana"] for c in campaings]

    localized_now = get_localized_date()

    stmt = select(ScheduledCampaigns).where(
        and_(
            ScheduledCampaigns.external_campaign_id.in_(campaing_ids),
            cast(ScheduledCampaigns.scheduled_at, Date) == datetime.date.today(),
        )
    )

    scheduled_campaigns = db.session.scalars(stmt)
    existing_ids = [sc.external_campaign_id for sc in scheduled_campaigns]

    campaigns_to_add = []
    for c in campaings:
        external_id = c["Cve_Campana"]
        if external_id in existing_ids:
            continue

        # Notifying
        final_date_str = f"{c['F_Final']} {c['Hora_Inicial']}"
        initial_date_str = f"{c['F_Inicial']} {c['Hora_Inicial']}"
        final_date = datetime.datetime.strptime(final_date_str, "%Y-%m-%d %H:%M:%S")
        initial_date = datetime.datetime.strptime(initial_date_str, "%Y-%m-%d %H:%M:%S")

        # checking between
        if (
            localized_now.date() >= initial_date.date()
            and localized_now.date() <= final_date.date()
        ):
            eta = localized_now.replace(
                hour=initial_date.hour,
                minute=initial_date.minute,
                second=initial_date.second,
            ).astimezone(pytz.UTC)
            if localized_now.date() >= final_date.date():
                response = change_status_campaign(external_campaign_id=external_id)
                if response.status_code != 200:
                    continue

            notify.apply_async(
                eta=eta, args=(external_id, c["Titulo"], c["Mensaje"], {})
            )
            campaings_scheduled += 1
            campaigns_to_add.append(
                ScheduledCampaigns(external_campaign_id=external_id, scheduled_at=eta),
            )

    db.session.add_all(campaigns_to_add)
    db.session.commit()

    return {"campaigns": campaings_scheduled, "error": ""}
