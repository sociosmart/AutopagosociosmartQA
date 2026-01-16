import asyncio
import random

from core.db import init_db
from core.models.gas_stations import GasStation, GasStationMargin
from core.models.products import Product


async def main():
    await init_db()

    stations = await GasStation.find().to_list()
    products = await Product.find().to_list()

    margins = []
    for _ in range(10_000):
        points = round(random.uniform(0.1, 1), 2)
        margin = round(random.uniform(0.1, 30), 2)

        margins.append(GasStationMargin(
            gas_station=random.choice(stations),
            product=random.choice(products),
            points=points,
            margin=margin
        ))


    try:
        await GasStationMargin.insert_many(margins)
    except Exception as e:
        print("Something went wrong", e)







if __name__ == "__main__":
    asyncio.run(main())
