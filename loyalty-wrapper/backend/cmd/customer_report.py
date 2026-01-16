import asyncio
import calendar
from datetime import datetime

from beanie import PydanticObjectId
from pydantic import BaseModel, Field

from core.db import init_db
from core.models.accumulations import Accumulation
from core.models.customers import Customer
from core.services.levels import level_service


async def generate_samples():
    await init_db()

    customer = await Customer.find_one(Customer.phone_number == "0512678970")

    start_date = datetime(2024, 7, 1)
    end_date = datetime(2024, 7, 31, 23, 59)

    accumulations = (
        await Accumulation.find(
            Accumulation.customer.id == customer.id,
            Accumulation.created_at >= start_date,
            Accumulation.created_at <= end_date,
        )
        .aggregate(
            [
                {
                    "$group": {
                        "_id": "$customer.$id",
                        "total": {"$sum": "$generated_points"},
                        "total_transactions": {"$sum": 1},
                    }
                }
            ],
            projection_model=AccumulationsInPeriod,
        )
        .to_list()
    )

    print(accumulations)


async def print_accumulations():
    await init_db()

    customer = await Customer.find_one(Customer.phone_number == "0512678970")

    accumulations = await Accumulation.find(
        Accumulation.customer.id == customer.id
    ).to_list()

    for a in accumulations:
        await a.fetch_link(Accumulation.customer)
        print(a.customer.id, a.created_at)


async def generate_customers_level():
    await init_db()
    now = datetime.now()
    range_ = calendar.monthrange(now.year, now.month)
    start_date = datetime(now.year, now.month, 1)
    end_date = datetime(now.year, now.month, range_[1], 23, 59, 59)
    await level_service.generate_customers_level(start_date, end_date)


if __name__ == "__main__":
    # asyncio.run(print_accumulations())
    # asyncio.run(generate_samples())
    asyncio.run(generate_customers_level())
