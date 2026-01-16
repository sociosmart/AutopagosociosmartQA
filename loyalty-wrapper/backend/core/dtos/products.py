from typing import Optional

from pydantic import BaseModel, Field


class ProductBase(BaseModel):
    name: str = Field(description="Product name", max_length=100)
    codename: str = Field(description="Codename", max_length=100)
    is_active: bool = Field(description="Product status")


class AddProductBody(ProductBase):
    pass


class UpdateProductBody(BaseModel):
    name: Optional[str] = Field(None, description="Product name", max_length=100)
    codename: Optional[str] = Field(None, description="Codename", max_length=100)
    is_active: Optional[bool] = Field(None, description="Product status")


# useful if we want to pass extra data
# TODO: Add user that is updating here
class UpdateProductInDb(UpdateProductBody):
    pass
