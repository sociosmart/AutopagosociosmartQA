from typing import Optional

import strawberry
from strawberry.permission import PermissionExtension

from core.exceptions.common import NotFoundException
from core.logging.logger import logger
from core.services.levels import level_service
from core.services.smart_gas import smart_gas_sevice
from core.utils.input_validation import validate_body
from gql.permissions import (
    IsAuthenticated,
    IsAuthorizedAppAuthenticated,
    IsCustomerAuthenticated,
)
from gql.types.errors import GeneralError
from gql.types.levels import (
    AddlevelBody,
    AddLevelResponse,
    CustomerLevel,
    CustomerLevelPagination,
    CustomerLevelPaginationResponse,
    GetCustomerLevelByPhoneResponse,
    GetCustomerLevelResponse,
    GetLevelByIdResponse,
    Level,
    LevelPagination,
    LevelPaginationResponse,
    UpdateLevelBody,
    UpdateLevelResponse,
)
from gql.types.pagination import PaginationParams


@strawberry.type
class LevelResolverQuery:
    @strawberry.field(
        extensions=[
            PermissionExtension(permissions=[IsAuthenticated()], fail_silently=True)
        ]
    )
    async def levels(
        self, pagination: Optional[PaginationParams], search: Optional[str] = None
    ) -> LevelPaginationResponse:
        errors = validate_body(pagination)
        if errors:
            return errors

        try:
            data = await level_service.get_levels_paginated(
                pagination.to_pydantic() if pagination else None, search
            )
        except Exception as e:
            logger.error(f"Unexpected error while paginating users - {e}")
            return GeneralError(code=500, message="Internal Server Error")

        return LevelPagination(**dict(data))

    @strawberry.field(
        extensions=[
            PermissionExtension(permissions=[IsAuthenticated()], fail_silently=True)
        ]
    )
    async def get_level_by_id(self, id: str) -> GetLevelByIdResponse:
        try:
            level = await level_service.get_level_by_id(id)

        except NotFoundException as e:
            return GeneralError(code=404, message=e.message)
        except Exception as e:
            logger.error(f"Unexpected error while getting level by id - {e}")
            return GeneralError(code=500, message="Internal Server Error")

        return Level(**level.model_dump())

    @strawberry.field(
        extensions=[
            PermissionExtension(permissions=[IsAuthenticated()], fail_silently=True)
        ]
    )
    async def customer_levels(
        self, pagination: Optional[PaginationParams]
    ) -> CustomerLevelPaginationResponse:
        errors = validate_body(pagination)
        if errors:
            return errors

        try:
            data = await level_service.get_customer_levels_paginated(
                pagination.to_pydantic() if pagination else None
            )
        except Exception as e:
            logger.error(f"Unexpected error while paginating users - {e}")
            return GeneralError(code=500, message="Internal Server Error")

        return CustomerLevelPagination(**dict(data))

    @strawberry.field(
        extensions=[
            PermissionExtension(
                permissions=[IsCustomerAuthenticated()], fail_silently=True
            )
        ]
    )
    async def customer_level(self, info: strawberry.Info) -> GetCustomerLevelResponse:

        customer = info.context.customer_data
        try:
            customer_level = await level_service.get_customer_level_by_phone(
                customer.phone_number
            )
        except NotFoundException:
            return GeneralError(code=404, message="No active level for customer")
        except Exception as e:
            logger.error(
                f"Unexpected error while getting customer level by phone - {e}"
            )
            return GeneralError(code=500, message="Internal Server Error")

        to_response = customer_level.__dict__
        del to_response["revision_id"]
        return CustomerLevel(**dict(to_response))

    @strawberry.field(
        extensions=[
            PermissionExtension(
                permissions=[IsAuthorizedAppAuthenticated()], fail_silently=True
            )
        ]
    )
    async def customer_level_by_phone(
        self, phone: str
    ) -> GetCustomerLevelByPhoneResponse:
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
            customer_level = await level_service.get_customer_level_by_phone(phone)
        except NotFoundException:
            return GeneralError(code=404, message="No active level for customer")
        except Exception as e:
            logger.error(
                f"Unexpected error while getting customer level by phone - {e}"
            )
            return GeneralError(code=500, message="Internal Server Error")

        to_response = customer_level.__dict__
        del to_response["revision_id"]
        return CustomerLevel(**dict(to_response))


@strawberry.type
class LevelResolverMutation:
    @strawberry.field(
        extensions=[
            PermissionExtension(permissions=[IsAuthenticated()], fail_silently=True)
        ]
    )
    async def add_level(self, body: AddlevelBody) -> AddLevelResponse:
        errors = validate_body(body)
        if errors:
            return errors

        try:
            level = await level_service.add_level(body.to_pydantic())
        except Exception as e:
            logger.error(f"Unexpected error while creating level - {e}")
            return GeneralError(code=500, message="Internal Server Error")

        return Level(**level.model_dump())

    @strawberry.field(
        extensions=[
            PermissionExtension(permissions=[IsAuthenticated()], fail_silently=True)
        ]
    )
    async def update_level(self, id: str, body: UpdateLevelBody) -> UpdateLevelResponse:
        errors = validate_body(body)
        if errors:
            return errors
        try:
            level = await level_service.update_level(id, body.to_pydantic())

        except NotFoundException as e:
            return GeneralError(code=404, message=e.message)
        except Exception as e:
            logger.error(f"Unexpected error while getting level by id - {e}")
            return GeneralError(code=500, message="Internal Server Error")

        return Level(**level.model_dump())
