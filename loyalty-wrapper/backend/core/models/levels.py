import datetime

from beanie import Document, Link
from pydantic import Field

from core.models.customers import Customer


class Level(Document):
    name: str = Field(max_length=100)
    min_points: float = Field(ge=0)
    is_active: bool = True

    class Settings:
        name = "levels"


class CustomerLevel(Document):
    customer: Link[Customer]
    level: Link[Level]
    start_date: datetime.datetime
    end_date: datetime.datetime
    is_active: bool = True

    class Settings:
        name = "customer_levels"
