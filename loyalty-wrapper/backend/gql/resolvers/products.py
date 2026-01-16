from typing import Optional

import strawberry
from beanie import PydanticObjectId
from strawberry.permission import PermissionExtension

from core.exceptions.common import NotFoundException
from core.logging.logger import logger
from core.services.products import products_service
from core.utils.input_validation import validate_body
from gql.permissions import IsAuthenticated
from gql.types.errors import GeneralError
from gql.types.pagination import PaginationParams
from gql.types.products import (
    AddProductBody,
    AddProductResponse,
    Product,
    ProductsPagination,
    ProductsPaginationResponse,
    UpdateProductBody,
    UpdateProductResponse,
    GetByIdProductsResponse,
)


@strawberry.type
class ProductResolverQuery:
    @strawberry.field(
        extensions=[
            PermissionExtension(permissions=[IsAuthenticated()], fail_silently=True)
        ]
    )
    async def products(
        self, pagination: Optional[PaginationParams] = None,
        search: Optional[str] = None 
    ) -> ProductsPaginationResponse:
        errors = validate_body(pagination)
        if errors:
            return errors

        try:
            data = await products_service.get_products_paginated(
                pagination.to_pydantic() if pagination else None,
                search
            )
        except Exception as e:
            logger.error(f"Something went wrong while getting products paginated - {e}")
            return GeneralError(code=500, message="Internal server error")

        return ProductsPagination(**dict(data))
    

    
    @strawberry.field(
        extensions=[
            PermissionExtension(permissions=[IsAuthenticated()], fail_silently=True)
        ]
    )
    async def get_product_by_id(self, id: str) -> GetByIdProductsResponse:
        try:
            product = await products_service.get_product_by_id(id)  

        except NotFoundException as e:
            return GeneralError(code=404, message=e.message)
        except Exception as e:
            logger.error(f"Unexpected error while getting product by id - {e}")
            return GeneralError(code=500, message="Internal Server Error")
        
        return Product(**product.model_dump())  


@strawberry.type
class ProductResolverMutation:
    @strawberry.field(
        extensions=[
            PermissionExtension(permissions=[IsAuthenticated()], fail_silently=True)
        ]
    )
    async def add_product(self, data: AddProductBody) -> AddProductResponse:
        errors = validate_body(data)

        if errors:
            return errors

        try:
            product = await products_service.create_product(data.to_pydantic())
        except Exception as e:
            logger.error(f"Something went wrong while adding product - {e}")
            return GeneralError(code=500, message="Internal server error")

        return Product(**product.model_dump())

    @strawberry.field(
        extensions=[
            PermissionExtension(permissions=[IsAuthenticated()], fail_silently=True)
        ]
    )
    async def update_product(
        self, id: str, data: UpdateProductBody
    ) -> UpdateProductResponse:
        errors = validate_body(data)

        if errors:
            return errors

        try:
            product = await products_service.update_product(id, data.to_pydantic())
        except NotFoundException as e:
            return GeneralError(code=404, message=e.message)
        except Exception as e:
            logger.error(f"Something went wrong while updating product - {e}")
            return GeneralError(code=500, message="Internal server error")

        return Product(**product.model_dump())




