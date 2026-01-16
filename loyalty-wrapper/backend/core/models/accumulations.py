from typing import Optional

from beanie import Document, Link, PydanticObjectId, View
from pydantic import BaseModel, Field

from core.enums.margins import MarginType
from core.models.customers import Customer
from core.models.gas_stations import GasStation
from core.models.mixins import CreatedAtMixin
from core.models.products import Product


class Accumulation(Document, CreatedAtMixin):
    # Redudant data with GasStationMargin
    # since data could change we have to keep a local version of what was applied
    customer: Link[Customer]
    margin: float = Field(ge=0, le=100)
    margin_type: MarginType = MarginType.by_margin
    gas_price: float = 0
    points: float = Field(ge=0, le=100)
    amount: float = Field(ge=0)
    gas_station: Optional[Link[GasStation]]
    product: Link[Product]
    generated_points: float = Field(ge=0)
    used_points: float = Field(0, ge=0)
    is_active: bool = True

    class Settings:
        name = "accumulations"


class AccumulationReportView(View):
    id: PydanticObjectId = Field(alias="_id")
    customer: Link[Customer]
    total_generated_points: float
    total_used_points: float
    total_amount: float
    avg_amount: float
    total_points: float
    total_transactions: int

    class Settings:
        source = Accumulation
        pipeline = [
            {
                "$group": {
                    "_id": "$customer.$id",
                    "total_generated_points": {"$sum": "$generated_points"},
                    "total_used_points": {"$sum": "$used_points"},
                    "customer": {"$first": "$customer"},
                    "total_transactions": {"$sum": 1},
                    "total_amount": {"$sum": "$amount"},
                    "avg_amount": {"$avg": "$amount"},
                }
            },
            {
                "$set": {
                    "total_points": {
                        "$subtract": ["$total_generated_points", "$total_used_points"]
                    }
                }
            },
        ]


class AccumulationsInPeriod(BaseModel):
    id: PydanticObjectId = Field(alias="_id")
    total: float
    total_transactions: int
