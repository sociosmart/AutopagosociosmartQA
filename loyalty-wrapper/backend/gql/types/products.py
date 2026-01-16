from typing import Annotated, List

import strawberry

from core.dtos.products import AddProductBody as PydanticAddProductBody
from core.dtos.products import UpdateProductBody as PydanticUpdateProductBody
from gql.types.errors import GeneralError, InputValidationError
from gql.types.pagination import CursorPage


@strawberry.type
class Product:
    id: str
    name: str
    codename: str
    is_active: bool


@strawberry.type
class ProductsPagination(CursorPage):
    items: List[Product]


@strawberry.experimental.pydantic.input(PydanticAddProductBody, all_fields=True)
class AddProductBody:
    pass


@strawberry.experimental.pydantic.input(PydanticUpdateProductBody, all_fields=True)
class UpdateProductBody:
    pass


AddProductResponse = Annotated[
    Product | InputValidationError | GeneralError | None,
    strawberry.union("AddProductResponse"),
]

UpdateProductResponse = Annotated[
    Product | InputValidationError | GeneralError | None,
    strawberry.union("UpdateProductResponse"),
]


ProductsPaginationResponse = Annotated[
    ProductsPagination | InputValidationError | GeneralError | None,
    strawberry.union("ProductsPaginationResponse"),
]


GetByIdProductsResponse = Annotated[
    Product | InputValidationError | GeneralError | None,
    strawberry.union("GetByIdProductsResponse"),
]
