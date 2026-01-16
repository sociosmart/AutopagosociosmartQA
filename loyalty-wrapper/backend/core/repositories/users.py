from typing import List, Optional

from core.dtos.pagination import CursorPage, PaginationParams
from core.models.users import User
from core.utils.pagination import Paginator
from core.exceptions.common import NotFoundException
from core.dtos.users import UpdateUserBody

class UserRepository:
    async def get_user_by_email(self, email) -> User | None:
        try:
            return await User.find_one(User.email == email)
        except Exception as e:
            raise e

    async def create_user(self, user: User) -> User:
        try:
            return await user.create()
        except Exception as e:
            raise e

    async def get_user_by_criterias(self, **filters) -> User | None:
        try:
            return await User.find_one(filters)
        except Exception as e:
            raise e

    async def get_users_by_criterias(self, **filters) -> List[User]:
        try:
            return await User.find(filters).to_list()
        except Exception as e:
            raise e

    async def get_users_paginated(
        self, pagination: Optional[PaginationParams] = None,
        search: Optional[str] = None
    ) -> CursorPage[User]:
        search_query = {}
        if search:
            search_query = {
                "$or": [
                    {"name": {"$regex": search, "$options": "i"}},
                    {"last_name": {"$regex": search, "$options": "i"}},
                    {"email": {"$regex": search, "$options": "i"}},
                    {
                        "$expr": {
                            "$regexMatch": {
                                "input": { "$concat": ["$first_name", " ", "$last_name"] },
                                "regex": search,
                                "options": "i"
                            }
                        }
                    }
                ]
            }

        return await Paginator(User, pagination).paginate(search_query)

        
    async def get_user_by_id(self, id: str) -> User:
        user = await User.get(id)
        if not user:
            raise NotFoundException
        return user
    

    async def update_user(self, id: str, data: UpdateUserBody) -> User:
        user = await self.get_user_by_id(id)

        await user.set(data.model_dump(exclude_none=True))

        return user

user_repository = UserRepository()
