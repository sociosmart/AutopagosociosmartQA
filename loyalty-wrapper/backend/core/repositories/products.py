from typing import Optional, Set

from beanie import PydanticObjectId
from bson import ObjectId

from core.dtos.pagination import CursorPage, PaginationParams
from core.dtos.products import UpdateProductInDb
from core.exceptions.common import NotFoundException
from core.models.products import Product
from core.utils.pagination import Paginator


class ProductRepository:
    async def get_products_paginated(
        self, pagination: Optional[PaginationParams] = None,
        search: Optional[str] = None
    ) -> CursorPage[Product]:
        search_query = {}
        if search:
            search_query = {
                "$or": [
                    {"name": {"$regex": search, "$options": "i"}},
                ]
            }

        return await Paginator(Product, pagination).paginate(search_query)



    async def create_product(self, product: Product) -> Product:
        return await product.create()

    async def update_product(self, id: str, data: UpdateProductInDb) -> Product:
        product = await Product.find_one({"_id": ObjectId(id)})
        if not product:
            raise NotFoundException
        await product.set(data.model_dump(exclude_none=True))
        return product
    
    async def get_product_by_id(self, id: str) -> Product:
        product = await Product.get(id)
        if not product:
            raise NotFoundException
        return product


product_repository = ProductRepository()
