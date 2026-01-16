import datetime
from typing import List, Optional

import httpx
from pymongo.errors import DuplicateKeyError

from core.config import settings
from core.dtos.customer import ExternalCostumerResponse
from core.dtos.gas_stations import ExternalGasStationResponse
from core.dtos.pagination import CursorPage, PaginationParams
from core.exceptions.auth import UnauthorizationException
from core.exceptions.common import NotFoundException
from core.logging.logger import logger
from core.models.customers import Customer
from core.models.gas_stations import GasStation
from core.repositories.customers import customer_repository
from core.repositories.gas_stations import gas_station_repository
from core.services.levels import level_service


class SmartGasService:
    def __init__(self):
        self.__repo = gas_station_repository
        self.__base_url = settings.external_services.smartgas_api_url
        self.__customer_repo = customer_repository

    async def fetch_external_gas_stations(self) -> List[ExternalGasStationResponse]:
        try:
            url = f"{self.__base_url}/rest/operacion"
            async with httpx.AsyncClient() as client:
                response = await client.get(url)
                response.raise_for_status()

            external_data = response.json()
            return [ExternalGasStationResponse(**station) for station in external_data]
        except httpx.RequestError as e:
            raise e
        except Exception as e:
            raise e

    async def save_gas_stations(
        self, external_data: List[ExternalGasStationResponse]
    ) -> List[GasStation]:
        try:
            current_stations = await self.__repo.get_all_gas_stations()
            current_ids = {station.external_id: station for station in current_stations}

            updated_stations = []
            new_stations = []

            for station in external_data:
                if station.external_id in current_ids:
                    updated_station = await self.update_gas_station(
                        current_ids[station.external_id], station
                    )
                    updated_stations.append(updated_station)
                else:
                    new_station = await self.create_gas_station(station)
                    new_stations.append(new_station)

            return new_stations + updated_stations
        except Exception as e:
            raise e

    async def create_gas_station(
        self, station: ExternalGasStationResponse
    ) -> GasStation:
        new_station = GasStation(
            name=station.name,
            external_id=station.external_id,
            cre_permission=station.cre_permission,
            latitude=station.latitude,
            longitude=station.longitude,
            city=station.city,
            regular_price=float(station.regular_price),
            premium_price=float(station.premium_price),
            diesel_price=float(station.diesel_price),
        )
        await self.__repo.create_gas_station(new_station)
        return new_station

    async def update_gas_station(
        self, existing_station: GasStation, station: ExternalGasStationResponse
    ) -> GasStation:
        existing_station.name = station.name
        existing_station.cre_permission = station.cre_permission
        existing_station.latitude = station.latitude
        existing_station.longitude = station.longitude
        existing_station.city = station.city
        existing_station.regular_price = float(station.regular_price)
        existing_station.premium_price = float(station.premium_price)
        existing_station.diesel_price = float(station.diesel_price)
        await self.__repo.update_gas_station(existing_station)
        return existing_station

    async def verify_customer(self, key: str, value: str) -> ExternalCostumerResponse:
        try:
            url = f"{self.__base_url}/rest/clientes?Verifica"

            async with httpx.AsyncClient() as client:
                response = await client.post(url, json={key: value})
                response.raise_for_status()

            data = response.json()

            if not data:
                raise NotFoundException
            if data[0].get("status") == "error":
                raise UnauthorizationException

            customer_data = ExternalCostumerResponse(**data[0])

            return customer_data

        except httpx.RequestError as e:
            raise Exception(f"Request error while verifying customer: {e}")
        except Exception as e:
            raise e

    async def update_or_create_customer(self, key: str, value) -> Customer | None:
        customer_data = await self.verify_customer(key, value)
        try:
            data = await self.__customer_repo.upsert_customer(
                customer_data.external_id,
                dict(
                    external_id=customer_data.external_id,
                    name=customer_data.name,
                    last_name=customer_data.last_name
                    + " "
                    + customer_data.second_last_name,
                    status=customer_data.status,
                    email=customer_data.email,
                    phone_number=customer_data.phone_number,
                    push_token=(
                        customer_data.push_token if customer_data.push_token else ""
                    ),
                    token=customer_data.token,
                ),
            )

            # This means that data has just been created
            if isinstance(data, Customer):
                try:
                    now = datetime.datetime.now()
                    await level_service.generate_customer_level(
                        data, date=now, is_new_user=True
                    )
                except Exception as e:
                    logger.error(
                        f"Unable to create level for customer {data.phone_number} - {e}"
                    )
                return data

            customer = await self.__customer_repo.get_customer_by_external_id(
                customer_data.external_id
            )

            return customer

        except DuplicateKeyError:
            # If duplicated key because multi tasking, we are forcing to get the customer btw
            customer = await self.__customer_repo.get_customer_by_external_id(
                customer_data.external_id
            )

            return customer

        except Exception as e:
            raise e

    async def get_customers_paginated(
        self, pagination: Optional[PaginationParams], **filters
    ) -> CursorPage[Customer]:
        try:
            return await self.__customer_repo.get_customers_paginated(pagination)
        except Exception as e:
            raise e

    async def post_notification(self, customer: Customer):
        uri = f"{self.__base_url}/rest/notifications"
        try:
            response = httpx.post(
                uri,
                json={
                    "Token": customer.token,
                    "identify": "InsertMessagge",
                    "Celular": customer.phone_number,
                },
            )

            if response.status_code != 200:
                logger.error(
                    f"Something went wrong while posting notification in smartgas - {response.status_code}"
                )
                raise Exception("Unable to post notification")
        except Exception as e:
            logger.error(f"Something went wrong while posting notification - {e}")
            raise e


smart_gas_sevice = SmartGasService()
