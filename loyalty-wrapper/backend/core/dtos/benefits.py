from typing import Optional

from pydantic import BaseModel, Field

from core.enums.benefits import BenefitFrequency, BenefitType


class AddBenefit(BaseModel):
    level: str = Field(description="Level id")
    name: str = Field(description="Benefit name")
    type: BenefitType = Field(BenefitType.digital, description="Benefit type")
    num_times: int = Field(0, description="Times")
    external_product_id: str = Field("", description="External product id")
    frequency: BenefitFrequency = Field(
        BenefitFrequency.n_times, description="Benefit frequency"
    )
    discount: float = Field(
        0, description="Applicable only for type gas and periferics"
    )
    stock: int = Field(-1, description="Benefit stock")
    is_active: bool = Field(True, description="Benefit status")
    dependency: bool = Field(False, description="Benefit over benefit")
    min_amount: float = Field(0.0, description="Minimum amount needed")


class UpdateGeneratedBenefit(BaseModel):
    stock: Optional[int] = Field(None, description="Benefit stock")
    is_active: Optional[bool] = Field(None, description="Benefit status")


class UpdateBenefit(BaseModel):
    level: Optional[str] = Field(None, description="Level id")
    name: Optional[str] = Field(None, description="Benefit name")
    num_times: Optional[int] = Field(None, description="Times")
    external_product_id: Optional[str] = Field(None, description="External product id")
    type: Optional[BenefitType] = Field(None, description="Benefit type")
    frequency: Optional[BenefitFrequency] = Field(None, description="Benefit frequency")
    discount: Optional[float] = Field(
        None, description="Applicable only for type gas and periferics"
    )
    stock: Optional[int] = Field(None, description="Benefit stock")
    is_active: Optional[bool] = Field(None, description="Benefit status")
    dependency: Optional[bool] = Field(None, description="Benefit over benefit")
    min_amount: Optional[float] = Field(None, description="Minimum amount needed")


class UpdateBenefitInDb(UpdateBenefit):
    pass
