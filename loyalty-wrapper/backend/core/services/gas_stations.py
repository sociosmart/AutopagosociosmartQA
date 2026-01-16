from typing import Optional

from pymongo.errors import DuplicateKeyError

from core.dtos.gas_stations import (
    AddGasStationMargin,
    UpdateGasStationMargin,
    UpdateGasStationMarginInDb,
)
from core.dtos.pagination import CursorPage, PaginationParams
from core.enums.margins import MarginType
from core.exceptions.common import DuplicatedKeyException, NotFoundException
from core.models.gas_stations import GasStation, GasStationMargin
from core.repositories.gas_stations import gas_station_repository


class GasStationService:
    def __init__(self):
        self.__repo = gas_station_repository

    async def get_gas_stations_paginated(
        self,
        pagination: Optional[PaginationParams],
        search: Optional[dict] = None,
        **filters
    ) -> CursorPage[GasStation]:
        return await self.__repo.get_gas_stations_paginated(pagination, search)

    async def get_gas_station_by_id(self, id: str) -> GasStation:
        return await self.__repo.get_gas_station_by_id(id)

    async def get_gas_stations_margin_paginated(
        self,
        pagination: Optional[PaginationParams],
        search: Optional[dict] = None,
        **filters
    ) -> CursorPage[GasStationMargin]:
        return await self.__repo.get_gas_stations_margin_paginated(pagination, search)

    async def get_gas_station_margin_by_id(self, id: str) -> GasStationMargin:
        return await self.__repo.get_gas_station_margin_by_id(id)

    async def create_gas_station_margin(
        self, data: AddGasStationMargin
    ) -> GasStationMargin:
        if data.gas_station:
            try:
                await self.get_gas_station_by_id(data.gas_station)
            except NotFoundException:
                raise NotFoundException(message="Gas Station not found")
            except Exception as e:
                raise e

        try:
            return await self.__repo.create_gas_station_margin(
                GasStationMargin(**data.model_dump())
            )
        except DuplicateKeyError:
            raise DuplicatedKeyException
        except Exception as e:
            raise e

    async def update_gas_station_margin(self, id: str, data: UpdateGasStationMargin):
        if data.gas_station:
            try:
                await self.get_gas_station_by_id(data.gas_station)
            except NotFoundException:
                raise NotFoundException(message="Gas Station not found")
            except Exception as e:
                raise e
        try:
            return await self.__repo.update_gas_station_margin(
                id, UpdateGasStationMarginInDb(**data.model_dump())
            )
        except DuplicateKeyError:
            raise DuplicatedKeyException
        except Exception as e:
            raise e

    async def get_gas_station_by_external_id(self, id: str):
        return await self.__repo.get_gas_station_by_external_id(id)

    async def get_applicable_margin_by_product(
        self,
        product_codename: str,
        external_gas_station_id: Optional[str],
    ):
        return await self.__repo.get_applicable_margin_by_product(
            product_codename, external_gas_station_id
        )


gas_station_sevice = GasStationService()
