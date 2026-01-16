from beanie import init_beanie
from motor.motor_asyncio import AsyncIOMotorClient

from core.config import settings
from core.models.accumulations import Accumulation, AccumulationReportView
from core.models.authorized_apps import AuthorizedApp
from core.models.benefits import Benefit, BenefitGenerated, BenefitTicket, PeriodCovered
from core.models.customers import Customer
from core.models.gas_stations import GasStation, GasStationMargin

# Registering models
from core.models.levels import CustomerLevel, Level
from core.models.products import Product
from core.models.settings import Setting
from core.models.users import User


async def init_db():
    client = AsyncIOMotorClient(str(settings.db.connection_uri))

    # Init beanie with the Product document class
    await init_beanie(
        database=client[settings.db.collection],
        document_models=[
            User,
            GasStation,
            Customer,
            Product,
            GasStationMargin,
            Accumulation,
            AccumulationReportView,
            Level,
            AuthorizedApp,
            CustomerLevel,
            Benefit,
            BenefitGenerated,
            Setting,
            BenefitTicket,
            PeriodCovered,
        ],
        recreate_views=True,
    )
