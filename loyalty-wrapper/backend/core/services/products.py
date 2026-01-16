from typing import Optional

from beanie.exceptions import DocumentNotFound

from core.dtos.pagination import PaginationParams
from core.dtos.products import AddProductBody, UpdateProductBody, UpdateProductInDb
from core.exceptions.common import NotFoundException
from core.models.products import Product
from core.repositories.products import product_repository


class ProductService:
    def __init__(self):
        self.__repo = product_repository

    async def get_products_paginated(
        self, pagination: Optional[PaginationParams] = None,
        search: Optional[dict] = None
    ):
        return await self.__repo.get_products_paginated(pagination, search)

    async def create_product(self, data: AddProductBody) -> Product:
        return await self.__repo.create_product(Product(**data.model_dump()))

    async def update_product(self, id: str, data: UpdateProductBody) -> Product:
        try:
            return await self.__repo.update_product(
                id, UpdateProductInDb(**data.model_dump())
            )
        except DocumentNotFound:
            raise NotFoundException
        except Exception as e:
            raise e

    async def get_product_by_id(self, id: str) -> Product:
        return await self.__repo.get_product_by_id(id)


products_service = ProductService()
