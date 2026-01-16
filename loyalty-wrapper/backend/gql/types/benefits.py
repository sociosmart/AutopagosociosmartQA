import datetime
from typing import Annotated, List, Optional

import strawberry

from core.dtos.benefits import AddBenefit as PydanticAddBenefit
from core.dtos.benefits import UpdateBenefit as PydanticUpdateBenefit
from core.dtos.benefits import UpdateGeneratedBenefit as PydanticUpdateGeneratedBenefit
from core.enums.benefits import BenefitFrequency, BenefitType
from gql.types.customers import Customer
from gql.types.errors import GeneralError, InputValidationError
from gql.types.levels import Level
from gql.types.pagination import CursorPage


@strawberry.type
class Benefit:
    id: str
    level: Level
    name: str = strawberry.field(description="Benefit name")
    type: BenefitType = strawberry.field(description="Benefit type")
    frequency: BenefitFrequency = strawberry.field(description="Benefit frequency")
    discount: float = strawberry.field(
        default=0,
        description="Applicable for gas and periferics",
    )
    num_times: int = strawberry.field(
        description="Hoy many times an event is going to occur"
    )
    external_product_id: str = strawberry.field(description="External smartgas id")
    stock: int = strawberry.field(default=-1, description="Stock for items")
    is_active: bool
    dependency: bool = strawberry.field(description="Benefit over benefit")
    min_amount: float = strawberry.field(
        description="Min amount needed for benefit over benefit"
    )

    created_at: datetime.datetime
    updated_at: Optional[datetime.datetime]


@strawberry.type
class BenefitGenerated(Benefit):
    benefit: Benefit = strawberry.field(description="parent benefit")
    stock_used: int = strawberry.field(description="Used stock.")
    start_date: datetime.datetime
    end_date: datetime.datetime


@strawberry.type
class BenefitPagination(CursorPage):
    items: List[Benefit]


@strawberry.type
class BenefitGeneratedPagination(CursorPage):
    items: List[BenefitGenerated]


@strawberry.experimental.pydantic.input(PydanticAddBenefit, all_fields=True)
class AddBenefit:
    pass


@strawberry.experimental.pydantic.input(PydanticUpdateBenefit, all_fields=True)
class UpdateBenefit:
    pass


@strawberry.experimental.pydantic.input(PydanticUpdateGeneratedBenefit, all_fields=True)
class UpdateGeneratedBenefit:
    pass


@strawberry.type
class BenefitTicket:
    id: str
    customer: Customer
    benefit_generated: BenefitGenerated
    start_date: datetime.datetime
    end_date: datetime.datetime
    redeemed: bool = False
    created_at: datetime.datetime
    updated_at: Optional[datetime.datetime]


@strawberry.type
class BenefitTicketPagination(CursorPage):
    items: List[BenefitTicket]


@strawberry.type
class ActiveBenefitsByPhone:
    items: List[BenefitTicket]


@strawberry.type
class ActiveBenefits:
    items: List[BenefitTicket]


@strawberry.type
class GasDiscount:
    discount: float = strawberry.field(description="Discount")


@strawberry.type
class Redemption:
    message: str


RedemptionByPhoneResponse = Annotated[
    Redemption | GeneralError | None, strawberry.union("RedemptionByPhoneResponse")
]
RedemptionResponse = Annotated[
    Redemption | GeneralError | None, strawberry.union("RedemptionResponse")
]

GasDiscountByPhoneResponse = Annotated[
    GasDiscount | GeneralError | None,
    strawberry.union("GasDiscountByPhoneResponse"),
]
GasDiscountResponse = Annotated[
    GasDiscount | GeneralError | None,
    strawberry.union("GasDiscountResponse"),
]

GetActiveBenefitsByPhoneResponse = Annotated[
    ActiveBenefitsByPhone | GeneralError | None,
    strawberry.union("GetActiveBenefitsByPhoneResponse"),
]
GetActiveBenefitsResponse = Annotated[
    ActiveBenefits | GeneralError | None,
    strawberry.union("GetActiveBenefitsResponse"),
]
BenefitTicketPaginationResponse = Annotated[
    BenefitTicketPagination | InputValidationError | GeneralError | None,
    strawberry.union("BenefitTicketPaginationResponse"),
]

UpdateGeneratedBenefitResponse = Annotated[
    BenefitGenerated | InputValidationError | GeneralError | None,
    strawberry.union("UpdateGeneratedBenefitResponse"),
]


AddBenefitResponse = Annotated[
    Benefit | InputValidationError | GeneralError | None,
    strawberry.union("AddBenefitResponse"),
]

UpdateBenefitResponse = Annotated[
    Benefit | InputValidationError | GeneralError | None,
    strawberry.union("UpdateBenefitResponse"),
]


BenefitPaginationResponse = Annotated[
    BenefitPagination | InputValidationError | GeneralError | None,
    strawberry.union("BenefitPaginationResponse"),
]

BenefitGeneratedPaginationResponse = Annotated[
    BenefitGeneratedPagination | InputValidationError | GeneralError | None,
    strawberry.union("BenefitGeneratedPaginationResponse"),
]

GetBenefitByIdResponse = Annotated[
    Benefit | GeneralError | None, strawberry.union("GetBenefitByIdResponse")
]

GetBenefitGeneratedByIdResponse = Annotated[
    BenefitGenerated | GeneralError | None,
    strawberry.union("GetBenefitGeneratedByIdResponse"),
]
