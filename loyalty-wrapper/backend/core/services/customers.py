from core.repositories.customers import customer_repository
from core.models.customers import Customer
from core.dtos.pagination import CursorPage, PaginationParams
from typing import Optional



class CustomerService:
    def __init__(self):
        self.__customer_repo = customer_repository
        
    async def get_customers_paginated(
        self, pagination: Optional[PaginationParams], 
        search: Optional[dict] = None
    ) -> CursorPage[Customer]:
        try:
            return await self.__customer_repo.get_customers_paginated(pagination, search)
        except Exception as e:
            raise e
        

    async def get_customer_by_id(self, id: str) -> Customer:
        return await self.__customer_repo.get_customer_by_id(id)


        




customer_sevice = CustomerService()