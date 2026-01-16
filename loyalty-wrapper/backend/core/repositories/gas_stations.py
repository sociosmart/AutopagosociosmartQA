from typing import List, Optional

from core.dtos.gas_stations import UpdateGasStationMarginInDb
from core.dtos.pagination import CursorPage, PaginationParams
from core.enums.margins import MarginType
from core.exceptions.common import NotFoundException
from core.models.gas_stations import GasStation, GasStationMargin
from core.utils.pagination import Paginator


class GasStationRepository:
    async def create_gas_station(self, gas_station: GasStation) -> GasStation:
        return await gas_station.create()

    async def get_all_gas_stations(self) -> List[GasStation]:
        return await GasStation.find_all().to_list()

    async def update_gas_station(self, station: GasStation):
        return await station.save()

    async def get_gas_stations_paginated(
        self,
        pagination: Optional[PaginationParams] = None,
        search: Optional[str] = None,
    ) -> CursorPage[GasStation]:
        search_query = {}

        if search:
            search_query = {
                "$or": [
                    {"name": {"$regex": search, "$options": "i"}},
                    {"cre_permission": {"$regex": search, "$options": "i"}},
                    {"city": {"$regex": search, "$options": "i"}},
                    {"external_id": search},
                ]
            }
        return await Paginator(GasStation, pagination).paginate(search_query)

    async def get_gas_station_by_id(self, id: str) -> GasStation:
        gas_station = await GasStation.get(id)
        if not gas_station:
            raise NotFoundException
        return gas_station

    async def get_gas_stations_margin_paginated(
        self,
        pagination: Optional[PaginationParams] = None,
        search: Optional[str] = None,
    ) -> CursorPage[GasStationMargin]:
        search_query = {}
        if search:
            search_query = {
                "$or": [
                    {"gas_station.name": {"$regex": search, "$options": "i"}},
                    {"product.name": {"$regex": search, "$options": "i"}},
                ]
            }

        return await Paginator(GasStationMargin, pagination).paginate(
            search_query, fetch_links=True, on_demand=False
        )

    async def get_gas_station_margin_by_id(self, id: str) -> GasStationMargin:
        gas_station_margin = await GasStationMargin.get(id, fetch_links=True)
        if not gas_station_margin:
            raise NotFoundException
        return gas_station_margin

    async def create_gas_station_margin(
        self, gas_station_margin: GasStationMargin
    ) -> GasStationMargin:
        gas_station_margin = await gas_station_margin.create()
        await gas_station_margin.fetch_all_links()
        return gas_station_margin

    async def update_gas_station_margin(
        self, id: str, data: UpdateGasStationMarginInDb
    ) -> GasStationMargin:
        gas_station_margin = await GasStationMargin.get(id)
        if not gas_station_margin:
            raise NotFoundException
        new_data = gas_station_margin.model_dump() | data.model_dump(
            exclude_defaults=True
        )
        new_model = GasStationMargin(**new_data)
        await new_model.replace()
        await new_model.fetch_all_links()
        return new_model

    async def get_gas_station_by_external_id(self, id: str) -> GasStation:
        gas_station = await GasStation.find_one(GasStation.external_id == id)

        if not gas_station:
            raise NotFoundException

        return gas_station

    async def get_applicable_margin_by_product(
        self,
        product_codename: str,
        exteral_gas_station_id: Optional[str],
    ) -> GasStationMargin:
        extra = {"gas_station.external_id": exteral_gas_station_id}

        gas_station_margin = await GasStationMargin.find_one(
            GasStationMargin.product.codename == product_codename,
            extra,
            fetch_links=True,
        )

        if not gas_station_margin:
            raise NotFoundException

        return gas_station_margin


gas_station_repository = GasStationRepository()
