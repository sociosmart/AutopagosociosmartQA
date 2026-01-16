from typing import Optional

from pydantic import BaseModel, Field


class AccumulationBody(BaseModel):
    product_codename: str = Field(description="Product codename")
    external_gas_station_id: Optional[str] = Field(
        None, description="Gas station id in SmartGas"
    )
    customer_phone: str = Field(description="Customer Phone")
    amount: float = Field(description="Amount")
