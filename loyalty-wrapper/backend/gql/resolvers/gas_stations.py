from typing import Optional

import strawberry
from strawberry.permission import PermissionExtension

from core.exceptions.common import DuplicatedKeyException, NotFoundException
from core.logging.logger import logger
from core.services.gas_stations import gas_station_sevice
from core.utils.input_validation import validate_body
from gql.permissions import IsAuthenticated, IsAuthorizedAppAuthenticated
from gql.types.errors import GeneralError
from gql.types.gas_stations import (
    AddGasStationMargin,
    AddGasStationMarginResponse,
    GasStation,
    GasStationMargin,
    GasStationMarginPagination,
    GasStationMarginPaginationResponse,
    GasStationPagination,
    GasStationPaginationResponse,
    GetByIdGasStationResponse,
    GetGasStationMarginByIdResponse,
    UpdateGasStationMargin,
    UpdateGasStationMarginResponse,
)
from gql.types.pagination import PaginationParams
from sync import populate_gas_stations


@strawberry.type
class GasStationResolverQuery:
    @strawberry.field(
        extensions=[
            PermissionExtension(permissions=[IsAuthenticated()], fail_silently=True)
        ]
    )
    async def gas_stations_margin(
        self, pagination: Optional[PaginationParams], search: Optional[str] = None
    ) -> GasStationMarginPaginationResponse:
        errors = validate_body(pagination)
        if errors:
            return errors

        try:
            data = await gas_station_sevice.get_gas_stations_margin_paginated(
                pagination.to_pydantic() if pagination else None, search
            )
        except Exception as e:
            logger.error(f"Unexpected error while paginating gas stations margin - {e}")
            return GeneralError(code=500, message="Internal Server Error")

        return GasStationMarginPagination(**dict(data))

    @strawberry.field(
        extensions=[
            PermissionExtension(permissions=[IsAuthenticated()], fail_silently=True)
        ]
    )
    async def gas_stations(
        self, pagination: Optional[PaginationParams], search: Optional[str] = None
    ) -> GasStationPaginationResponse:
        errors = validate_body(pagination)
        if errors:
            return errors

        try:
            data = await gas_station_sevice.get_gas_stations_paginated(
                pagination.to_pydantic() if pagination else None, search
            )
        except Exception as e:
            logger.error(f"Unexpected error while paginating gas stations - {e}")
            return GeneralError(code=500, message="Internal Server Error")
        return GasStationPagination(**dict(data))

    @strawberry.field(
        extensions=[
            PermissionExtension(permissions=[IsAuthenticated()], fail_silently=True)
        ]
    )
    async def get_gas_station_by_id(self, id: str) -> GetByIdGasStationResponse:
        try:
            gas_station = await gas_station_sevice.get_gas_station_by_id(id)

        except NotFoundException as e:
            return GeneralError(code=404, message=e.message)
        except Exception as e:
            logger.error(f"Unexpected error while getting gas station by id - {e}")
            return GeneralError(code=500, message="Internal Server Error")

        return GasStation(**gas_station.model_dump())

    @strawberry.field(
        extensions=[
            PermissionExtension(permissions=[IsAuthenticated()], fail_silently=True)
        ]
    )
    async def get_gas_station_margin_by_id(
        self, id: str
    ) -> GetGasStationMarginByIdResponse:
        try:
            gas_station_margin = await gas_station_sevice.get_gas_station_margin_by_id(
                id
            )

        except NotFoundException as e:
            return GeneralError(code=404, message=e.message)
        except Exception as e:
            logger.error(
                f"Unexpected error while getting gas station margin by id - {e}"
            )
            return GeneralError(code=500, message="Internal Server Error")

        # TODO: Find a better way to use revision_id
        to_response = gas_station_margin.__dict__
        del to_response["revision_id"]
        return GasStationMargin(**dict(gas_station_margin))


@strawberry.type
class GasStationResolverMutation:
    @strawberry.field(
        extensions=[
            PermissionExtension(
                permissions=[IsAuthorizedAppAuthenticated()], fail_silently=True
            )
        ]
    )
    async def sync_gas_stations(self) -> Optional[str]:
        populate_gas_stations.delay()
        return "Synced"

    @strawberry.field(
        extensions=[
            PermissionExtension(permissions=[IsAuthenticated()], fail_silently=True)
        ]
    )
    async def add_gas_station_margin(
        self, data: AddGasStationMargin
    ) -> AddGasStationMarginResponse:
        errors = validate_body(data)
        if errors:
            return errors

        try:
            gas_station_margin = await gas_station_sevice.create_gas_station_margin(
                data.to_pydantic()
            )
        except DuplicatedKeyException:
            return GeneralError(
                code=412,
                message="Given gas station and product already present",
            )
        except Exception as e:
            logger.error(f"Unexpected error while creating gas station margin - {e}")
            return GeneralError(code=500, message="Internal Server error")

        to_response = gas_station_margin.__dict__
        del to_response["revision_id"]
        return GasStationMargin(**dict(gas_station_margin))

    @strawberry.field(
        extensions=[
            PermissionExtension(permissions=[IsAuthenticated()], fail_silently=True)
        ]
    )
    async def update_gas_station_margin(
        self, id: str, data: UpdateGasStationMargin
    ) -> UpdateGasStationMarginResponse:
        errors = validate_body(data)
        if errors:
            return errors

        try:
            gas_station_margin = await gas_station_sevice.update_gas_station_margin(
                id, data.to_pydantic()
            )
        except NotFoundException as e:
            return GeneralError(code=404, message=e.message)
        except DuplicatedKeyException:
            return GeneralError(
                code=412,
                message="Given gas station and product already present",
            )
        except Exception as e:
            logger.error(f"Unexpected error while creating gas station margin - {e}")
            return GeneralError(code=500, message="Internal Server error")

        to_response = gas_station_margin.__dict__
        del to_response["revision_id"]
        return GasStationMargin(**dict(gas_station_margin))
