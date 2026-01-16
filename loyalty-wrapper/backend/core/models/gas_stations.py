from typing import Optional

import pymongo
from beanie import Document, Indexed, Link
from pydantic import BaseModel, EmailStr, Field
from pymongo import IndexModel

from core.enums.margins import MarginType
from core.models.products import Product


class GasStation(Document):
    name: str
    external_id: Indexed(str, unique=True)
    cre_permission: str
    latitude: str
    longitude: str
    city: str
    regular_price: float = 0
    premium_price: float = 0
    diesel_price: float = 0

    class Settings:
        name = "gas_stations"


class GasStationMargin(Document):
    margin_type: MarginType = MarginType.by_margin
    margin: float = Field(ge=0, le=100)
    points: float = Field(ge=0, le=100)
    gas_station: Optional[Link[GasStation]]
    product: Link[Product]

    class Settings:
        name = "gas_station_margins"
        indexes = [
            IndexModel(
                [
                    ("gas_station", pymongo.DESCENDING),
                    ("product", pymongo.DESCENDING),
                ],
                unique=True,
            )
        ]
