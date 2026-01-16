import datetime
from typing import Optional

from bson import ObjectId

from core.dtos.levels import UpdateLevelBody
from core.dtos.pagination import CursorPage, PaginationParams
from core.exceptions.common import NotFoundException
from core.models.levels import CustomerLevel, Level
from core.utils.pagination import Paginator


class LevelRepository:
    async def get_levels_paginated(
        self, pagination: Optional[PaginationParams], search: Optional[str] = None
    ) -> CursorPage[Level]:
        search_query = {}
        if search:
            search_query = {
                "$or": [
                    {"name": {"$regex": search, "$options": "i"}},
                ]
            }

        return await Paginator(Level, pagination).paginate(search_query)

    async def add_level(self, level: Level) -> Level:
        return await level.create()

    async def update_level(self, id: str, data: UpdateLevelBody) -> Level:
        level = await self.get_level_by_id(id)

        await level.set(data.model_dump(exclude_none=True))

        return level

    async def get_level_by_id(self, id) -> Level:
        level = await Level.get(id)
        if not level:
            raise NotFoundException
        return level

    async def get_customer_levels_paginated(
        self, pagination: Optional[PaginationParams]
    ) -> CursorPage[CustomerLevel]:
        return await Paginator(CustomerLevel, pagination).paginate(fetch_links=True)

    async def create_customer_level(
        self, customer_level: CustomerLevel
    ) -> CustomerLevel:
        cl = await customer_level.create()
        await cl.fetch_all_links()
        return cl

    async def get_customer_level(
        self, customer_id: str, date: datetime.datetime
    ) -> CustomerLevel:
        customer_level = await CustomerLevel.find_one(
            CustomerLevel.customer.id == ObjectId(customer_id),
            CustomerLevel.start_date <= date,
            CustomerLevel.end_date >= date,
        )

        if not customer_level:
            raise NotFoundException

        return customer_level

    async def get_customer_level_by_phone(
        self, phone: str, date: datetime.datetime
    ) -> CustomerLevel:
        customer_level = await CustomerLevel.find_one(
            CustomerLevel.customer.phone_number == phone,
            CustomerLevel.start_date <= date,
            CustomerLevel.end_date >= date,
            fetch_links=True,
        )

        if not customer_level:
            raise NotFoundException

        return customer_level

    async def get_suitable_level(self, points: float) -> Level | None:
        levels = (
            await Level.find(Level.min_points <= points, Level.is_active == True)
            .sort(-Level.min_points)
            .limit(1)
            .to_list()
        )
        return levels[0] if levels else None


level_repository = LevelRepository()
