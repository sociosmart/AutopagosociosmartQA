import asyncio
import datetime
import logging
from typing import Tuple

from core.db import init_db
from core.enums.benefits import BenefitGeneratorKeys
from core.models.benefits import BenefitGenerated, PeriodCovered
from core.services.benefits import benefit_service
from core.services.levels import level_service
from core.services.settings import settings_service

logger = logging.getLogger(__name__)


async def get_active_period(date: datetime.datetime) -> PeriodCovered | None:
    return await PeriodCovered.find_one(
        PeriodCovered.start_date <= date,
        PeriodCovered.end_date >= date,
        PeriodCovered.is_active == True,
    )


async def get_period_based_on_frequency() -> (
    Tuple[datetime.datetime, datetime.datetime]
):

    now = datetime.datetime.now()

    setting = await settings_service.get_by_key(BenefitGeneratorKeys.levels_duration)

    days = 30

    if setting:
        days = int(setting.value)

    return now, now + datetime.timedelta(days=days)


async def async_generate_levels():
    await init_db()

    now = datetime.datetime.now()
    active_period = await get_active_period(now)
    if active_period:
        logger.info(
            f"Period covered from - {active_period.start_date} - {active_period.end_date}"
        )
        return

    start_date, end_date = await get_period_based_on_frequency()

    # saving period
    try:
        await PeriodCovered(start_date=start_date, end_date=end_date).insert()
    except Exception as e:
        logger.error(f"Unable to save period covered")
        return

    # Getting active benefits
    benefits = await benefit_service.get_benefits(is_active=True)

    # Getting active benefits
    benefit_ids = [benefit.id for benefit in benefits]

    # Active benefits in period

    active_benefits = await benefit_service.get_generated_benefits_in_period(
        benefit_ids,
        now,
    )

    active_benefits_id = [b.benefit.id for b in active_benefits]

    benefits_to_activate = []
    for benefit in benefits:
        if benefit.id in active_benefits_id:
            continue

        generated_benefit = BenefitGenerated(
            benefit=benefit.id,
            start_date=start_date,
            end_date=end_date,
            **benefit.model_dump(
                exclude=["revision_id", "id", "created_at", "updated_at"]
            ),
        )
        benefits_to_activate.append(generated_benefit)

    if benefits_to_activate:
        try:
            await benefit_service.insert_generated_benefits_many(benefits_to_activate)
        except Exception as e:
            logger.error(f"Unable to activate benefits - {e}")
            return

    # Generating user level and benefits
    await level_service.generate_customers_level()


def generate_levels():
    asyncio.run(async_generate_levels())


if __name__ == "__main__":
    generate_levels()
