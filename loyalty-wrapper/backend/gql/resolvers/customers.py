from typing import Optional

import strawberry
from strawberry.permission import PermissionExtension

from core.exceptions.common import NotFoundException
from core.logging.logger import logger
from core.services.customers import customer_sevice
from core.utils.input_validation import validate_body
from gql.permissions import IsAuthenticated, IsCustomerAuthenticated
from gql.types.customers import (
    Customer,
    CustomerPagination,
    CustomersPaginationResponse,
    GetByIdCustomersResponse,
    MeCustomer,
)
from gql.types.errors import GeneralError
from gql.types.pagination import PaginationParams


@strawberry.type
class CustomersResolverQuery:
    @strawberry.field(
        extensions=[
            PermissionExtension(permissions=[IsAuthenticated()], fail_silently=True)
        ]
    )
    async def customers(
        self, pagination: Optional[PaginationParams], search: Optional[str] = None
    ) -> CustomersPaginationResponse:
        errors = validate_body(pagination)
        if errors:
            return errors

        try:
            data = await customer_sevice.get_customers_paginated(
                pagination.to_pydantic() if pagination else None, search
            )
        except Exception as e:
            logger.error(f"Unexpected error while paginating users - {e}")
            return GeneralError(code=500, message="Internal Server Error")
        return CustomerPagination(**dict(data))

    @strawberry.field(
        extensions=[
            PermissionExtension(
                permissions=[IsCustomerAuthenticated()], fail_silently=True
            )
        ]
    )
    async def me_customer(self, info: strawberry.Info) -> Optional[MeCustomer]:
        customer = info.context.customer_data

        MeCustomer(
            external_id=customer.external_id,
            name=customer.name,
            last_name=customer.last_name,
            phone_number=customer.phone_number,
            email=customer.email,
        )

        return customer

    @strawberry.field(
        extensions=[
            PermissionExtension(permissions=[IsAuthenticated()], fail_silently=True)
        ]
    )
    async def get_customer_by_id(self, id: str) -> GetByIdCustomersResponse:
        try:
            customer = await customer_sevice.get_customer_by_id(id)

        except NotFoundException as e:
            return GeneralError(code=404, message=e.message)
        except Exception as e:
            logger.error(f"Unexpected error while getting user by id - {e}")
            return GeneralError(code=500, message="Internal Server Error")

        return Customer(**customer.model_dump())
