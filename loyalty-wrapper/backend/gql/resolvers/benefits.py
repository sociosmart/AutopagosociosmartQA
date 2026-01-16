from typing import Optional

import strawberry
from strawberry.permission import PermissionExtension

from core.enums.benefits import BenefitType
from core.exceptions.common import (
    BenefitAlreadyRedeemedException,
    BenefitNotInDateRangeException,
    BenefitNotStockLeftException,
    NotFoundException,
)
from core.logging.logger import logger
from core.services.benefits import benefit_service
from core.services.smart_gas import smart_gas_sevice
from core.utils.input_validation import validate_body
from gql.permissions import (
    IsAuthenticated,
    IsAuthorizedAppAuthenticated,
    IsCustomerAuthenticated,
)
from gql.types.benefits import (
    ActiveBenefits,
    ActiveBenefitsByPhone,
    AddBenefit,
    AddBenefitResponse,
    Benefit,
    BenefitGenerated,
    BenefitGeneratedPagination,
    BenefitGeneratedPaginationResponse,
    BenefitPagination,
    BenefitPaginationResponse,
    BenefitTicket,
    BenefitTicketPagination,
    BenefitTicketPaginationResponse,
    GasDiscount,
    GasDiscountByPhoneResponse,
    GasDiscountResponse,
    GetActiveBenefitsByPhoneResponse,
    GetActiveBenefitsResponse,
    GetBenefitByIdResponse,
    GetBenefitGeneratedByIdResponse,
    Redemption,
    RedemptionByPhoneResponse,
    RedemptionResponse,
    UpdateBenefit,
    UpdateBenefitResponse,
    UpdateGeneratedBenefit,
    UpdateGeneratedBenefitResponse,
)
from gql.types.errors import GeneralError
from gql.types.pagination import PaginationParams


@strawberry.type
class BenefitResolverQuery:
    @strawberry.field(
        extensions=[
            PermissionExtension(permissions=[IsAuthenticated()], fail_silently=True)
        ]
    )
    async def benefits(
        self, pagination: Optional[PaginationParams], search: Optional[str] = None
    ) -> BenefitPaginationResponse:
        errors = validate_body(pagination)
        if errors:
            return errors

        try:
            data = await benefit_service.get_benefits_paginated(
                pagination.to_pydantic() if pagination else None, search
            )
        except Exception as e:
            logger.error(f"Unexpected error while paginating gas stations margin - {e}")
            return GeneralError(code=500, message="Internal Server Error")

        return BenefitPagination(**dict(data))

    @strawberry.field(
        extensions=[
            PermissionExtension(permissions=[IsAuthenticated()], fail_silently=True)
        ]
    )
    async def benefits_generated(
        self, pagination: Optional[PaginationParams], search: Optional[str] = None
    ) -> BenefitGeneratedPaginationResponse:
        errors = validate_body(pagination)
        if errors:
            return errors

        try:
            data = await benefit_service.get_benefits_generated_paginated(
                pagination.to_pydantic() if pagination else None, search
            )
        except Exception as e:
            logger.error(f"Unexpected error while paginating gas stations margin - {e}")
            return GeneralError(code=500, message="Internal Server Error")

        return BenefitGeneratedPagination(**dict(data))

    @strawberry.field(
        extensions=[
            PermissionExtension(permissions=[IsAuthenticated()], fail_silently=True)
        ]
    )
    async def benefits_tickets(
        self, pagination: Optional[PaginationParams], search: Optional[str] = None
    ) -> BenefitTicketPaginationResponse:
        errors = validate_body(pagination)
        if errors:
            return errors

        try:
            data = await benefit_service.get_benefits_ticket_paginated(
                pagination.to_pydantic() if pagination else None, search
            )
        except Exception as e:
            logger.error(f"Unexpected error while paginating benefits tickets - {e}")
            return GeneralError(code=500, message="Internal Server Error")

        return BenefitTicketPagination(**dict(data))

    @strawberry.field(
        extensions=[
            PermissionExtension(permissions=[IsAuthenticated()], fail_silently=True)
        ]
    )
    async def get_benefit_by_id(self, id: str) -> GetBenefitByIdResponse:
        try:
            benefit = await benefit_service.get_benefit_by_id(id)

        except NotFoundException as e:
            return GeneralError(code=404, message=e.message)
        except Exception as e:
            logger.error(f"Unexpected error while getting benefit by id - {e}")
            return GeneralError(code=500, message="Internal Server Error")

        to_response = benefit.__dict__
        del to_response["revision_id"]
        return Benefit(**dict(benefit))

    @strawberry.field(
        extensions=[
            PermissionExtension(permissions=[IsAuthenticated()], fail_silently=True)
        ]
    )
    async def get_benefit_generated_by_id(
        self, id: str
    ) -> GetBenefitGeneratedByIdResponse:
        try:
            benefit = await benefit_service.get_benefit_generated_by_id(id)

        except NotFoundException as e:
            return GeneralError(code=404, message=e.message)
        except Exception as e:
            logger.error(
                f"Unexpected error while getting benefit generated by id - {e}"
            )
            return GeneralError(code=500, message="Internal Server Error")

        to_response = benefit.__dict__
        del to_response["revision_id"]
        return BenefitGenerated(**dict(benefit))

    @strawberry.field(
        extensions=[
            PermissionExtension(
                permissions=[IsAuthorizedAppAuthenticated()], fail_silently=True
            )
        ]
    )
    async def get_active_benefits_by_phone(
        self, phone: str
    ) -> GetActiveBenefitsByPhoneResponse:
        try:
            customer = await smart_gas_sevice.update_or_create_customer(
                key="Celular", value=phone
            )
        except NotFoundException:
            return GeneralError(
                code=404,
                message="Not customer found",
            )
        except Exception as e:
            logger.error(f"Unexpected error while creating/updating customer - {e}")
            return GeneralError(code=500, message="Internal Server Error")

        try:
            active_benefits = await benefit_service.get_active_benefits_by_phone(
                phone, BenefitType.physical
            )
        except Exception as e:
            logger.error(f"Unexpected error while getting active benefits - {e}")
            return GeneralError(code=500, message="Internal Server Error")

        benefits = []
        for benefit_ticket in active_benefits:
            to_add = benefit_ticket.__dict__
            del to_add["revision_id"]
            benefits.append(BenefitTicket(**to_add))

        return ActiveBenefitsByPhone(items=benefits)

    @strawberry.field(
        extensions=[
            PermissionExtension(
                permissions=[IsCustomerAuthenticated()], fail_silently=True
            )
        ]
    )
    async def get_active_benefits(
        self, info: strawberry.Info
    ) -> GetActiveBenefitsResponse:

        customer = info.context.customer_data
        try:
            active_benefits = await benefit_service.get_active_benefits_by_phone(
                customer.phone_number
            )
        except Exception as e:
            logger.error(f"Unexpected error while getting active benefits - {e}")
            return GeneralError(code=500, message="Internal Server Error")

        benefits = []
        for benefit_ticket in active_benefits:
            to_add = benefit_ticket.__dict__
            del to_add["revision_id"]
            benefits.append(BenefitTicket(**to_add))

        return ActiveBenefits(items=benefits)

    @strawberry.field(
        extensions=[
            PermissionExtension(
                permissions=[IsCustomerAuthenticated()], fail_silently=True
            )
        ]
    )
    async def gas_discount(self, info: strawberry.Info) -> GasDiscountResponse:

        customer = info.context.customer_data

        discount = 0

        try:
            benefit_ticket = await benefit_service.get_gas_discount_by_phone(
                customer.phone_number
            )
        except NotFoundException:
            discount = 0
        except Exception as e:
            logger.error(f"Unexpected error while getting discount ticket - {e}")
            return GeneralError(code=500, message="Internal Server Error")
        else:
            discount = benefit_ticket.benefit_generated.discount

        return GasDiscount(discount=discount)

    @strawberry.field(
        extensions=[
            PermissionExtension(
                permissions=[IsAuthorizedAppAuthenticated()], fail_silently=True
            )
        ]
    )
    async def gas_discount_by_phone(self, phone: str) -> GasDiscountByPhoneResponse:
        try:
            customer = await smart_gas_sevice.update_or_create_customer(
                key="Celular", value=phone
            )
        except NotFoundException:
            return GeneralError(
                code=404,
                message="Not customer found",
            )
        except Exception as e:
            logger.error(f"Unexpected error while creating/updating customer - {e}")
            return GeneralError(code=500, message="Internal Server Error")

        discount = 0

        try:
            benefit_ticket = await benefit_service.get_gas_discount_by_phone(phone)
        except NotFoundException:
            discount = 0
        except Exception as e:
            logger.error(f"Unexpected error while getting discount ticket - {e}")
            return GeneralError(code=500, message="Internal Server Error")
        else:
            discount = benefit_ticket.benefit_generated.discount

        return GasDiscount(discount=discount)


@strawberry.type
class BenefitResolverMutation:
    @strawberry.field(
        extensions=[
            PermissionExtension(permissions=[IsAuthenticated()], fail_silently=True)
        ]
    )
    async def add_benefit(self, body: AddBenefit) -> AddBenefitResponse:
        errors = validate_body(body)
        if errors:
            return errors

        try:
            benefit = await benefit_service.create_benefit(body.to_pydantic())
        except NotFoundException as e:
            return GeneralError(code=404, message=e.message)
        except Exception as e:
            logger.error(f"Unexpected error while creating benefit - {e}")
            return GeneralError(code=500, message="Internal Server error")

        to_response = benefit.__dict__
        del to_response["revision_id"]
        return Benefit(**dict(benefit))

    @strawberry.field(
        extensions=[
            PermissionExtension(permissions=[IsAuthenticated()], fail_silently=True)
        ]
    )
    async def update_benefit(
        self, id: str, body: UpdateBenefit
    ) -> UpdateBenefitResponse:
        errors = validate_body(body)
        if errors:
            return errors
        try:
            benefit = await benefit_service.update_benefit(id, body.to_pydantic())
        except NotFoundException as e:
            return GeneralError(code=404, message=e.message)
        except Exception as e:
            logger.error(f"Unexpected error while updating benefit - {e}")
            return GeneralError(code=500, message="Internal Server error")

        to_response = benefit.__dict__
        del to_response["revision_id"]
        return Benefit(**dict(benefit))

    @strawberry.field(
        extensions=[
            PermissionExtension(permissions=[IsAuthenticated()], fail_silently=True)
        ]
    )
    async def update_generated_benefit(
        self, id: str, body: UpdateGeneratedBenefit
    ) -> UpdateGeneratedBenefitResponse:
        errors = validate_body(body)
        if errors:
            return errors
        try:
            benefit = await benefit_service.update_generated_benefit(
                id, body.to_pydantic()
            )
        except NotFoundException as e:
            return GeneralError(code=404, message=e.message)
        except Exception as e:
            logger.error(f"Unexpected error while updating benefit - {e}")
            return GeneralError(code=500, message="Internal Server error")

        to_response = benefit.__dict__
        del to_response["revision_id"]
        return BenefitGenerated(**dict(benefit))

    @strawberry.field(
        description="Mutation used to redeem a benefit assigned. GAS and Periferics are not applicable",
        extensions=[
            PermissionExtension(
                permissions=[IsAuthorizedAppAuthenticated()], fail_silently=True
            )
        ],
    )
    async def redemption_by_phone(
        self, phone: str, benefit_ticket_id: str
    ) -> RedemptionByPhoneResponse:
        try:
            customer = await smart_gas_sevice.update_or_create_customer(
                key="Celular", value=phone
            )
        except NotFoundException:
            return GeneralError(
                code=404,
                message="Not customer found",
            )
        except Exception as e:
            logger.error(f"Unexpected error while creating/updating customer - {e}")
            return GeneralError(code=500, message="Internal Server Error")

        try:
            await benefit_service.redeem_benefit(benefit_ticket_id, customer.id)
        except NotFoundException:
            return GeneralError(
                code=404,
                message="Ticket not found",
            )
        except BenefitAlreadyRedeemedException as e:
            return GeneralError(code=412, message=e.message)
        except BenefitNotInDateRangeException as e:
            return GeneralError(code=412, message=e.message)
        except BenefitNotStockLeftException as e:
            return GeneralError(code=400, message=e.message)
        except Exception as e:
            logger.error(f"Something went wrong while redeeming benefit - {e}")
            return GeneralError(code=500, message="Internal Server Error")

        return Redemption(message="Redeemed")

    @strawberry.field(
        description="Mutation used to redeem a benefit assigned. GAS and Periferics are not applicable",
        extensions=[
            PermissionExtension(
                permissions=[IsCustomerAuthenticated()], fail_silently=True
            )
        ],
    )
    async def redemption_by_customer(
        self, benefit_ticket_id: str, info: strawberry.Info
    ) -> RedemptionResponse:
        customer = info.context.customer_data

        try:
            await benefit_service.redeem_benefit(benefit_ticket_id, customer.id)
        except NotFoundException:
            return GeneralError(
                code=404,
                message="Ticket not found",
            )
        except BenefitAlreadyRedeemedException as e:
            return GeneralError(code=412, message=e.message)
        except BenefitNotInDateRangeException as e:
            return GeneralError(code=412, message=e.message)
        except BenefitNotStockLeftException as e:
            return GeneralError(code=400, message=e.message)
        except Exception as e:
            logger.error(f"Something went wrong while redeeming benefit - {e}")
            return GeneralError(code=500, message="Internal Server Error")

        return Redemption(message="Redeemed")
