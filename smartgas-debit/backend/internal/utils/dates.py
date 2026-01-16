import datetime
import pytz


def localize_datetime(date: datetime.datetime, tz=None):
    if not tz:
        tz = pytz.timezone("America/Mazatlan")

    return date.astimezone(tz)


def localized_date(date: datetime.datetime, tz=None):
    if not tz:
        tz = pytz.timezone("America/Mazatlan")

    return tz.localize(date)


def localized_now(tz=None) -> datetime.datetime:
    if not tz:
        tz = pytz.timezone("America/Mazatlan")

    now = datetime.datetime.now(tz)

    return now
