from typing import Annotated, List, Optional

import strawberry

from core.dtos.gas_stations import AddGasStationMargin as PydanticAddGasStationMargin
from core.dtos.gas_stations import (
    UpdateGasStationMargin as PydanticUpdateGasStationMargin,
)
from core.enums.margins import MarginType
from gql.types.errors import GeneralError, InputValidationError
from gql.types.pagination import CursorPage
from gql.types.products import Product


@strawberry.type
class GasStation:
    id: str
    name: str
    external_id: str
    cre_permission: str
    latitude: str
    longitude: str
    city: str
    regular_price: float
    premium_price: float
    diesel_price: float


@strawberry.type
class GasStationMargin:
    id: str
    margin_type: MarginType = MarginType.by_margin
    margin: float = strawberry.field(description="Margin (%) to apply in accumulation")
    points: float = strawberry.field(description="Points per unit to apply")
    product: Product = strawberry.field(description="Product")
    gas_station: Optional[GasStation] = strawberry.field(description="Gas station")


@strawberry.type
class GasStationPagination(CursorPage):
    items: List[GasStation]


@strawberry.type
class GasStationMarginPagination(CursorPage):
    items: List[GasStationMargin]


@strawberry.experimental.pydantic.input(PydanticAddGasStationMargin, all_fields=True)
class AddGasStationMargin:
    pass


@strawberry.experimental.pydantic.input(PydanticUpdateGasStationMargin, all_fields=True)
class UpdateGasStationMargin:
    pass


UpdateGasStationMarginResponse = Annotated[
    GasStationMargin | InputValidationError | GeneralError | None,
    strawberry.union("UpdateGasStationMarginResponse"),
]
AddGasStationMarginResponse = Annotated[
    GasStationMargin | InputValidationError | GeneralError | None,
    strawberry.union("AddGasStationMarginResponse"),
]


GetGasStationMarginByIdResponse = Annotated[
    GasStationMargin | GeneralError | None,
    strawberry.union("GetGasStationMarginByIdResponse"),
]

GasStationMarginPaginationResponse = Annotated[
    GasStationMarginPagination | InputValidationError | GeneralError | None,
    strawberry.union("GasStationMarginPaginationResponse"),
]

GasStationPaginationResponse = Annotated[
    GasStationPagination | InputValidationError | GeneralError | None,
    strawberry.union("GasStationPaginationResponse"),
]


GetByIdGasStationResponse = Annotated[
    GasStation | GeneralError | None,
    strawberry.union("GetByIdGasStationResponse"),
]
