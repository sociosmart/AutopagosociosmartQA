import datetime
from typing import Annotated, List, Optional

import strawberry

from core.dtos.accumulations import AccumulationBody as PydanticAccumulationBody
from core.enums.margins import MarginType
from gql.types.benefits import BenefitTicket
from gql.types.customers import Customer
from gql.types.errors import GeneralError, InputValidationError
from gql.types.gas_stations import GasStation
from gql.types.pagination import CursorPage
from gql.types.products import Product


@strawberry.type
class Accumulation:
    id: str
    margin: float = strawberry.field(description="Applied Margin")
    margin_type: MarginType = MarginType.by_margin
    customer: Customer = strawberry.field(
        description="Customer who made the transaction"
    )
    points: float = strawberry.field(description="Applied points per mxn peso")
    product: Product = strawberry.field(description="Product")
    gas_station: Optional[GasStation] = strawberry.field(description="Gas station")
    amount: float = strawberry.field(description="Amount")
    generated_points: float = strawberry.field(description="Points Generated")
    gas_price: float = 0
    is_active: bool = strawberry.field(description="Is active")
    used_points: float = strawberry.field(description="Points used in this transaction")
    created_at: datetime.datetime = strawberry.field(description="Transaction date")


@strawberry.type
class AccumulationWithBenefits:
    accumulation: Accumulation
    benefits: List[BenefitTicket]


@strawberry.type
class AccumulationReport:
    id: str = strawberry.field(description="Customer id")
    customer: Customer = strawberry.field(
        description="Customer who made the transaction"
    )
    total_amount: float = strawberry.field(description="Total amount spent and used")
    avg_amount: float = strawberry.field(description="Avg amount")
    total_points: float = strawberry.field(description="Total points generated")
    total_transactions: int = strawberry.field(description="Total # of transactions")
    total_generated_points: float = strawberry.field(
        description="Total points generated over time"
    )
    total_used_points: float = strawberry.field(description="Total used points")


@strawberry.experimental.pydantic.input(PydanticAccumulationBody, all_fields=True)
class AccumulationBody:
    pass


@strawberry.type
class AccumulationPagination(CursorPage):
    items: List[Accumulation]


@strawberry.type
class AccumulationReportPagination(CursorPage):
    items: List[AccumulationReport]


AddAccumulationResponse = Annotated[
    AccumulationWithBenefits | GeneralError | InputValidationError | None,
    strawberry.union("AddAccumulationResponse"),
]

AccumulationPaginationResponse = Annotated[
    AccumulationPagination | GeneralError | InputValidationError | None,
    strawberry.union("AccumulationPaginationResponse"),
]

AccumulationReportPaginationResponse = Annotated[
    AccumulationReportPagination | GeneralError | InputValidationError | None,
    strawberry.union("AccumulationReportPaginationResponse"),
]
