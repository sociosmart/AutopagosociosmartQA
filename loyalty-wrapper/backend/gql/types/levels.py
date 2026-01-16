import datetime
from typing import Annotated, List

import strawberry

from core.dtos.levels import AddLevelBody as PydanticAddLevelBody
from core.dtos.levels import UpdateLevelBody as PydanticUpdateLevelBody
from gql.types.customers import Customer
from gql.types.errors import GeneralError, InputValidationError
from gql.types.pagination import CursorPage


@strawberry.type
class Level:
    id: str
    name: str = strawberry.field(description="Level name")
    min_points: float = strawberry.field(description="Min points to get this level")
    is_active: bool = strawberry.field(description="Level status")


@strawberry.type
class CustomerLevel:
    id: str
    customer: Customer
    level: Level
    start_date: datetime.datetime
    end_date: datetime.datetime
    is_active: bool


@strawberry.experimental.pydantic.input(PydanticAddLevelBody, all_fields=True)
class AddlevelBody:
    pass


@strawberry.experimental.pydantic.input(PydanticUpdateLevelBody, all_fields=True)
class UpdateLevelBody:
    pass


@strawberry.type
class LevelPagination(CursorPage):
    items: List[Level]


@strawberry.type
class CustomerLevelPagination(CursorPage):
    items: List[CustomerLevel]


GetLevelByIdResponse = Annotated[
    Level | GeneralError | None,
    strawberry.union("GetLevelByIdResponse"),
]

GetCustomerLevelByPhoneResponse = Annotated[
    CustomerLevel | GeneralError | None,
    strawberry.union("GetCustomerLevelByPhoneResponse"),
]
GetCustomerLevelResponse = Annotated[
    CustomerLevel | GeneralError | None,
    strawberry.union("GetCustomerLevelResponse"),
]

AddLevelResponse = Annotated[
    Level | InputValidationError | GeneralError | None,
    strawberry.union("AddLevelResponse"),
]

UpdateLevelResponse = Annotated[
    Level | InputValidationError | GeneralError | None,
    strawberry.union("UpdateLevelResponse"),
]


LevelPaginationResponse = Annotated[
    LevelPagination | InputValidationError | GeneralError | None,
    strawberry.union("LevelPaginationResponse"),
]

CustomerLevelPaginationResponse = Annotated[
    CustomerLevelPagination | InputValidationError | GeneralError | None,
    strawberry.union("CustomerLevelPaginationResponse"),
]
