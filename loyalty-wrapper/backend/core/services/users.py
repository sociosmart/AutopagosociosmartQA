from typing import Optional

from core.dtos.pagination import CursorPage, PaginationParams
from core.models.users import User
from core.repositories.users import user_repository
from core.exceptions.common import NotFoundException
from beanie.exceptions import DocumentNotFound
from core.dtos.users import UpdateUserBody

class UserService:
    def __init__(self) -> None:
        self.__repo = user_repository

    async def get_users_paginated(
        self, pagination: Optional[PaginationParams],
        search: Optional[dict] = None, **filters
    ) -> CursorPage[User]:
        try:
            return await self.__repo.get_users_paginated(pagination, search)
        except Exception as e:
            raise e
        

    async def get_user_by_id(self, id: str) -> User:
        return await self.__repo.get_user_by_id(id)




user_service = UserService()
