import datetime
from typing import Any, Dict

from beanie import Document, Link
from pydantic import BaseModel, Field, model_validator

from core.enums.benefits import BenefitFrequency, BenefitType
from core.models.customers import Customer
from core.models.levels import Level
from core.models.mixins import CreatedAtMixin, UpdatedAtMixin


class BenefitBase(BaseModel):
    level: Link[Level]
    name: str = Field(max_length=100)
    type: BenefitType = BenefitType.digital
    external_product_id: str = Field("", description="SmartGas product id")
    frequency: BenefitFrequency = BenefitFrequency.n_times
    discount: float = Field(0, ge=0, description="Discount value")
    num_times: int = Field(0, ge=0, description="Times")
    stock: int = Field(-1, description="Stock for benefit, for example. 1000 termos.")
    is_active: bool = True
    dependency: bool = Field(False, description="Benefit over benefit")
    min_amount: float = Field(0.0, description="Minimum amount needed")

    @model_validator(mode="before")
    @classmethod
    def validator(cls, data: Dict[str, Any]) -> Dict[str, Any]:
        if (
            data.get("type") == BenefitType.gas
            or data.get("type") == BenefitType.periferics
        ):
            data["frequency"] = BenefitFrequency.always
            data["num_times"] = 0
            data["stock"] = 0
            data["external_product_id"] = ""
            data["dependency"] = False
            data["min_amount"] = 0.0
        elif data.get("frequency") != BenefitFrequency.n_times:
            data["num_times"] = 0
            data["value"] = 0
        return data


class Benefit(Document, BenefitBase, CreatedAtMixin, UpdatedAtMixin):
    class Settings:
        name = "benefits"


# After a benefit is generated in a period, we should make a copy, just in case a benefit changes in the period
class BenefitGenerated(Document, BenefitBase, CreatedAtMixin, UpdatedAtMixin):
    benefit: Link[Benefit]
    stock_used: int = Field(0, description="Used stock.")
    start_date: datetime.datetime
    end_date: datetime.datetime

    class Settings:
        name = "benefits_generated"


class BenefitTicket(Document, CreatedAtMixin, UpdatedAtMixin):
    customer: Link[Customer]
    benefit_generated: Link[BenefitGenerated]
    start_date: datetime.datetime
    end_date: datetime.datetime
    redeemed: bool = Field(False, description="Has been redeemed")

    class Settings:
        name = "benefit_tickets"


class PeriodCovered(Document):
    start_date: datetime.datetime
    end_date: datetime.datetime
    is_active: bool = True

    class Settings:
        name = "period_covered"
