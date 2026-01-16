from typing import Annotated, List

import strawberry

from gql.types.errors import GeneralError, InputValidationError
from gql.types.pagination import CursorPage
from core.dtos.users import AddUserBody as PydanticAddUserBody
from core.dtos.users import UpdateUserBody as PydanticUpdateUserBody

@strawberry.type
class User:
    first_name: str
    last_name: str
    id: str
    email: str
    is_active: bool


@strawberry.type
class UserPagination(CursorPage):
    items: List[User]

@strawberry.experimental.pydantic.input(PydanticAddUserBody, all_fields=True)
class AddUserBody:
    pass


UsersPaginationResponse = Annotated[
    UserPagination | InputValidationError | GeneralError | None,
    strawberry.union("UsersPaginationResponse"),
]


GetByIdUsersResponse = Annotated[
    User | GeneralError | None,
    strawberry.union("GetByIdUsersResponse"),
]


AddUserResponse = Annotated[
    User | InputValidationError | GeneralError | None,
    strawberry.union("AddUserResponse"),
]


@strawberry.experimental.pydantic.input(PydanticUpdateUserBody, all_fields=True)
class UpdateUserBody:
    pass


UpdateUserResponse = Annotated[
    User | InputValidationError | GeneralError | None,
    strawberry.union("UpdateUserResponse"),
]
