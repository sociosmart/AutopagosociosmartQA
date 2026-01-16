from typing import Annotated, List

import strawberry

from gql.types.errors import GeneralError, InputValidationError
from gql.types.pagination import CursorPage


@strawberry.type
class Customer:
    id: str
    external_id: str
    name: str
    last_name: str
    status: str
    phone_number: str
    email: str


@strawberry.type
class CustomerPagination(CursorPage):
    items: List[Customer]


@strawberry.type
class MeCustomer:
    external_id: str
    name: str
    last_name: str
    phone_number: str
    email: str


CustomersPaginationResponse = Annotated[
    CustomerPagination | InputValidationError | GeneralError | None,
    strawberry.union("CustomerPaginationResponse"),
]

GetByIdCustomersResponse = Annotated[
    Customer | GeneralError | None,
    strawberry.union("GetByIdCustomersResponse"),
]
