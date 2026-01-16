import asyncio

from core.db import init_db
from core.models.accumulations import Accumulation
from core.models.benefits import BenefitGenerated, BenefitTicket, PeriodCovered
from core.models.customers import Customer
from core.models.levels import CustomerLevel


async def cleanup():
    await init_db()
    await Customer.delete_all()
    await Accumulation.delete_all()
    await CustomerLevel.delete_all()
    await BenefitTicket.delete_all()
    await BenefitGenerated.delete_all()
    await PeriodCovered.delete_all()


if __name__ == "__main__":
    asyncio.run(cleanup())
