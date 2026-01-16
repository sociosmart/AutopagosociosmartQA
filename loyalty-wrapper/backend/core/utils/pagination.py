from typing import Generic, List, Optional, TypeVar, Union

from beanie import Document, View
from bson import ObjectId

from core.dtos.pagination import CursorPage, PaginationParams
from gql.types.pagination import CursorPage as GraphqlCursorPage

T = TypeVar("T", bound=Union[Document, View])
T_Cast = TypeVar("T_Cast")


class Paginator(Generic[T]):
    items: List[T]
    limit: int = 20
    _prev_cursor: Optional[str] = None
    _next_cursor: Optional[str] = None

    def __init__(self, model: T, pagination: Optional[PaginationParams] = None):
        self.model = model
        if pagination:
            self.limit = pagination.limit
            self._next_cursor = pagination.next_cursor
            self._prev_cursor = pagination.prev_cursor

    async def _get_data(
        self, *filters, fetch_links=False, on_demand=True
    ) -> CursorPage[T]:
        conditions = {}
        sort = -1
        if self._prev_cursor:
            conditions["_id"] = {"$gte": ObjectId(self._prev_cursor)}
            sort = 1
        elif self._next_cursor:
            conditions["_id"] = {"$lte": ObjectId(self._next_cursor)}
        self.items = (
            await self.model.find(conditions)
            .find(*filters, fetch_links=fetch_links and not on_demand)
            .limit(self.limit + 1)
            .sort(("_id", sort))
            .to_list()
        )

        total = await self.model.find(
            *filters, fetch_links=fetch_links and not on_demand
        ).count()

        new_prev_cursor = None
        new_next_cursor = None

        if len(self.items) > self.limit:
            if self._prev_cursor:
                self.items = sorted(
                    self.items, key=lambda item: str(item.id), reverse=True
                )
            new_next_cursor = str(self.items.pop(-1).id)
            if self.items and (self._prev_cursor or self._next_cursor):
                new_prev_cursor = str(self.items[0].id)
        elif len(self.items) == 1 and self._prev_cursor:
            # No items in prev cursor page
            self.items = []

        # If not more cursors we get last page with data here
        if not new_prev_cursor and not new_next_cursor:
            new_prev_cursor = self._next_cursor
            new_next_cursor = self._prev_cursor

        # On damnd fetching
        if fetch_links and on_demand:
            for v in self.items:
                await v.fetch_all_links()
        return CursorPage[T](
            items=self.items,
            next_cursor=new_next_cursor,
            prev_cursor=new_prev_cursor,
            total=total,
        )

    #
    async def paginate(
        self, *filters, fetch_links=False, on_demand=True
    ) -> CursorPage[T]:
        return await self._get_data(
            *filters, fetch_links=fetch_links, on_demand=on_demand
        )
