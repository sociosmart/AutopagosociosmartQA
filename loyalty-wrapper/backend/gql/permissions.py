from typing import Any

from strawberry import Info
from strawberry.permission import BasePermission


class IsAuthenticated(BasePermission):
    message = "User is not authenticated"
    error_extensions = {"code": "unauthorized"}

    async def has_permission(self, source: Any, info: Info, **kwargs: Any) -> bool:
        return await info.context.user


class IsCustomerAuthenticated(BasePermission):
    message = "Customer is not authenticated"
    error_extensions = {"code": "unauthorized"}

    async def has_permission(self, source: Any, info: Info, **kwargs: Any) -> bool:
        return await info.context.customer


class IsAuthorizedAppAuthenticated(BasePermission):
    message = "Unauthorized app"
    error_extensions = {"code": "unauthorized"}

    async def has_permission(self, source: Any, info: Info, **kwargs: Any) -> bool:
        return await info.context.authorized_app

