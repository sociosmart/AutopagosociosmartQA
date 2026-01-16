from typing import List, Optional

from beanie.odm.operators.update.general import Set

from core.dtos.pagination import CursorPage, PaginationParams
from core.exceptions.common import NotFoundException
from core.models.customers import Customer
from core.utils.pagination import Paginator


class CustomerRepository:
    async def create_customer(self, customer: Customer) -> Customer:

        # DOUBLE CHECK
        c = await self.get_customer_by_external_id(customer.external_id)

        if c:
            return c

        return await customer.create()

    async def update_customer(self, customer: Customer):
        try:
            return await customer.save()
        except Exception as e:
            raise e

    async def upsert_customer(self, external_id: str, data: dict):
        return await Customer.find_one(Customer.external_id == external_id).upsert(
            Set(data),
            on_insert=Customer(**data),
        )

    async def get_customer_by_external_id(self, external_id: str) -> Customer | None:
        return await Customer.find_one(Customer.external_id == external_id)

    async def get_customers_paginated(
        self,
        pagination: Optional[PaginationParams] = None,
        search: Optional[str] = None,
    ) -> CursorPage[Customer]:
        search_query = {}
        if search:
            search_query = {
                "$or": [
                    {"name": {"$regex": search, "$options": "i"}},
                    {"last_name": {"$regex": search, "$options": "i"}},
                    {"phone_number": {"$regex": search, "$options": "i"}},
                    {"email": {"$regex": search, "$options": "i"}},
                    {
                        "$expr": {
                            "$regexMatch": {
                                "input": {"$concat": ["$name", " ", "$last_name"]},
                                "regex": search,
                                "options": "i",
                            }
                        }
                    },
                ]
            }
        return await Paginator(Customer, pagination).paginate(search_query)

    async def get_customer_by_id(self, id: str) -> Customer:
        customer = await Customer.get(id)
        if not customer:
            raise NotFoundException
        return customer

    async def get_customers(self, **filters) -> List[Customer]:
        return await Customer.find(filters).to_list()


customer_repository = CustomerRepository()
