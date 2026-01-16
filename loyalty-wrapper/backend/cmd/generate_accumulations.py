import asyncio
import datetime
import os
import random
import sys
from calendar import monthrange

sys.path.append(os.path.join(os.path.dirname(__file__), ".."))
from core.db import init_db
from core.models.accumulations import Accumulation
from core.models.customers import Customer
from core.models.gas_stations import GasStation, GasStationMargin


async def generate_samples(n: int):
    await init_db()

    gas_stations_margin = await GasStationMargin.find().to_list()
    customers = await Customer.find().to_list()
    accumulations = []
    for _ in range(n):
        year = random.choice([2024, 2025])
        month = random.randint(1, 12)
        day = random.choice([1, monthrange(year, month)[1]])
        hour = random.randint(0, 23)
        minute = random.randint(0, 59)
        date = datetime.datetime(year, month, day, hour, minute)

        amount = round(random.uniform(1.0, 1000), 2)
        gas_station_margin = random.choice(gas_stations_margin)
        generated_points = (
            amount * (gas_station_margin.margin / 100) * gas_station_margin.points
        )
        accumulations.append(
            Accumulation(
                **gas_station_margin.model_dump(
                    include=["points", "margin", "product", "gas_station"]
                ),
                customer=random.choice(customers),
                amount=amount,
                generated_points=generated_points,
                created_at=date,
            )
        )

    await Accumulation.insert_many(accumulations)


if __name__ == "__main__":
    asyncio.run(generate_samples(1000))
