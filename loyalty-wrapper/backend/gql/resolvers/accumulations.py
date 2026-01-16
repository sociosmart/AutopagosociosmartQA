from typing import Optional

import strawberry
from strawberry.permission import PermissionExtension

from core.exceptions.common import NotFoundException
from core.logging.logger import logger
from core.services.accumulations import accumulation_service
from core.services.benefits import benefit_service
from core.services.gas_stations import gas_station_sevice
from core.services.push_notifications import push_notifications_service
from core.services.smart_gas import smart_gas_sevice
from core.utils.input_validation import validate_body
from gql.permissions import IsAuthenticated, IsAuthorizedAppAuthenticated
from gql.types.accumulations import (
    Accumulation,
    AccumulationBody,
    AccumulationPagination,
    AccumulationPaginationResponse,
    AccumulationReportPagination,
    AccumulationReportPaginationResponse,
    AccumulationWithBenefits,
    AddAccumulationResponse,
)
from gql.types.benefits import BenefitTicket
from gql.types.errors import GeneralError
from gql.types.pagination import PaginationParams


@strawberry.type
class AccumulationResolverQuery:
    @strawberry.field(
        extensions=[
            PermissionExtension(permissions=[IsAuthenticated()], fail_silently=True)
        ]
    )
    async def accumulations(
        self, pagination: Optional[PaginationParams], search: Optional[str] = None
    ) -> AccumulationPaginationResponse:
        errors = validate_body(pagination)
        if errors:
            return errors

        try:
            data = await accumulation_service.get_accumulations_paginated(
                pagination.to_pydantic() if pagination else None, search
            )
        except Exception as e:
            logger.error(f"Unexpected error while paginating gas stations margin - {e}")
            return GeneralError(code=500, message="Internal Server Error")

        return AccumulationPagination(**dict(data))

    @strawberry.field(
        extensions=[
            PermissionExtension(permissions=[IsAuthenticated()], fail_silently=True)
        ]
    )
    async def accumulations_report(
        self, pagination: Optional[PaginationParams]
    ) -> AccumulationReportPaginationResponse:
        errors = validate_body(pagination)
        if errors:
            return errors

        try:
            data = await accumulation_service.get_accumulations_report_paginated(
                pagination.to_pydantic() if pagination else None
            )
        except Exception as e:
            logger.error(
                f"Unexpected error while paginating gas stations report margin - {e}"
            )
            return GeneralError(code=500, message="Internal Server Error")

        return AccumulationReportPagination(**dict(data))


@strawberry.type
class AccumulationResolverMutation:
    @strawberry.field(
        extensions=[
            PermissionExtension(
                permissions=[IsAuthorizedAppAuthenticated()], fail_silently=True
            )
        ]
    )
    async def accumulate(self, body: AccumulationBody) -> AddAccumulationResponse:
        errors = validate_body(body)

        if errors:
            return errors

        data = body.to_pydantic()

        try:
            margin = await gas_station_sevice.get_applicable_margin_by_product(
                data.product_codename,
                data.external_gas_station_id,
            )
        except NotFoundException:
            logger.info(f"Not applicable margin for this product")
            return GeneralError(
                code=412,
                message="Not applicable margin for this product and/or gas station",
            )
        except Exception as e:
            logger.error(f"Unexpected error while getting margin - {e}")
            return GeneralError(code=500, message="Internal Server Error")

        try:
            customer = await smart_gas_sevice.update_or_create_customer(
                key="Celular", value=data.customer_phone
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
            accumulation = await accumulation_service.generate_accumulation(
                margin=margin, customer=customer, amount=data.amount
            )
        except Exception as e:
            logger.error(f"Unexpected error while creating accumulation - {e}")
            return GeneralError(code=500, message="Internal Server Error")

        try:
            benefits = await benefit_service.get_active_dependent_benefits(
                customer.phone_number, data.amount
            )
        except Exception as e:
            logger.error(f"Something went wrong while bringing benefits {e}")
            return GeneralError(code=500, message="Internal Server Error")

        to_response = accumulation.__dict__
        del to_response["revision_id"]

        benefit_tickets = []
        for b in benefits:
            b_dict = b.__dict__
            del b_dict["revision_id"]
            benefit_tickets.append(BenefitTicket(**b_dict))

        accumulation_response = AccumulationWithBenefits(
            accumulation=Accumulation(**to_response), benefits=benefit_tickets
        )

        # Push notification
        try:
            if customer.push_token:
                push_notifications_service.send_notification(
                    title="Nueva acumulacion",
                    message="Se ha generado una nueva acumulacion en su cuenta.",
                    tokens=[customer.push_token],
                )
        except Exception as e:
            pass

        # Post notification in smart gas
        try:
            # customer
            if customer.token:
                await smart_gas_sevice.post_notification(customer)
        except Exception as e:
            pass

        return accumulation_response
