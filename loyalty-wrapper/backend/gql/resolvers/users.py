from typing import Optional

import strawberry
from strawberry.permission import PermissionExtension

from core.logging.logger import logger
from core.services.users import user_service
from core.utils.input_validation import validate_body
from gql.permissions import IsAuthenticated
from gql.types.errors import GeneralError
from gql.types.pagination import PaginationParams
from gql.types.users import UserPagination, UsersPaginationResponse, GetByIdUsersResponse, AddUserResponse, AddUserBody, UpdateUserBody, UpdateUserResponse
from gql.types.users import User
from core.exceptions.common import NotFoundException
from core.services.auth import auth_service


@strawberry.type
class UsersResolverQuery:
    @strawberry.field(
        extensions=[
            PermissionExtension(permissions=[IsAuthenticated()], fail_silently=True)
        ]
    )
    async def users(
        self, pagination: Optional[PaginationParams],
        search: Optional[str] = None 
    ) -> UsersPaginationResponse:
        errors = validate_body(pagination)
        if errors:
            return errors

        try:
            data = await user_service.get_users_paginated(
                pagination.to_pydantic() if pagination else None,
                search
            )
        except Exception as e:
            logger.error(f"Unexpected error while paginating users - {e}")
            return GeneralError(code=500, message="Internal Server Error")

        return UserPagination(**dict(data))
    

    @strawberry.field(
        extensions=[
            PermissionExtension(permissions=[IsAuthenticated()], fail_silently=True)
        ]
    )
    async def get_user_by_id(self, id: str) -> GetByIdUsersResponse:
        try:
            user = await user_service.get_user_by_id(id)  

        except NotFoundException as e:
            return GeneralError(code=404, message=e.message)
        except Exception as e:
            logger.error(f"Unexpected error while getting user by id - {e}")
            return GeneralError(code=500, message="Internal Server Error")
        
        return User(**user.model_dump())  


    @strawberry.field(
        extensions=[
            PermissionExtension(
                permissions=[IsAuthenticated()], fail_silently=True
            )
        ]
    )
    async def me_user(self, info: strawberry.Info) -> Optional[User]:
        user = info.context.user_data

        User(
            first_name=user.first_name,
            last_name=user.last_name,
            id=user.id,
            email=user.email,
            is_active=user.is_active,
        )

        return user


@strawberry.type
class UserResolverMutation:
    @strawberry.field(
        extensions=[
            PermissionExtension(permissions=[IsAuthenticated()], fail_silently=True)
        ]
    )
    async def add_user(self, data: AddUserBody) -> AddUserResponse:
        errors = validate_body(data)

        if errors:
            return errors

        try:
            user = await auth_service.create_user(data.to_pydantic())
        except Exception as e:
            logger.error(f"Something went wrong while adding user - {e}")
            return GeneralError(code=500, message="Internal server error")

        return User(**user.model_dump())
    

    @strawberry.field(
        extensions=[
            PermissionExtension(permissions=[IsAuthenticated()], fail_silently=True)
        ]
    )
    async def update_user(self, id: str, body: UpdateUserBody) -> UpdateUserResponse:
        errors = validate_body(body)
        if errors:
            return errors
        try:
            user = await auth_service.update_user(id, body.to_pydantic())

        except NotFoundException as e:
            return GeneralError(code=404, message=e.message)
        except Exception as e:
            logger.error(f"Unexpected error while getting user by id - {e}")
            return GeneralError(code=500, message="Internal Server Error")

        return User(**user.model_dump())
