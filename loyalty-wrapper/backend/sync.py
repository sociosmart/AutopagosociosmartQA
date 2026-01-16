import asyncio

from celery import Celery
from celery.schedules import crontab

from core.config import settings
from core.db import init_db
from core.services.smart_gas import smart_gas_sevice
from tasks import async_generate_levels

app = Celery(__name__, broker=str(settings.redis.redis_connection_uri))


app.conf.beat_schedule = {
    "level-generations": {
        "task": "sync.gen_levels",
        "schedule": crontab(minute=0, hour="0,12"),
        "args": (),
    },
    "populate-gas-stations": {
        "task": "sync.populate_gas_stations",
        "schedule": crontab(minute="*/15"),
        "args": (),
    },
}


async def _popolute_gas_stations():

    await init_db()

    external_data = await smart_gas_sevice.fetch_external_gas_stations()

    new_stations = await smart_gas_sevice.save_gas_stations(external_data)

    if new_stations:
        print("Nuevas estaciones guardadas")
    else:
        print("No se encontraron estaciones nuevas.")


@app.task
def populate_gas_stations():
    asyncio.run(_popolute_gas_stations())


@app.task
def gen_levels():
    asyncio.run(async_generate_levels())


app.conf.timezone = settings.tz
