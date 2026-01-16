import datetime
from typing import Optional

from core.dtos.pagination import CursorPage, PaginationParams
from core.enums.margins import MarginType
from core.models.accumulations import Accumulation, AccumulationReportView
from core.models.customers import Customer
from core.models.gas_stations import GasStationMargin
from core.repositories.accumulations import accumulation_repository


class AccumulationService:
    def __init__(self):
        self.__repo = accumulation_repository

    async def get_accumulations_paginated(
        self, pagination: Optional[PaginationParams], search: Optional[dict] = None
    ) -> CursorPage[Accumulation]:
        return await self.__repo.get_accumulations_paginated(pagination, search)

    async def get_accumulations_report_paginated(
        self, pagination: Optional[PaginationParams]
    ) -> CursorPage[AccumulationReportView]:
        return await self.__repo.get_accumulations_report_paginated(pagination)

    async def get_accumulations_report_in_period(
        self,
        customer_id: str,
        start_date: datetime.datetime,
        end_date: datetime.datetime,
    ):
        return await self.__repo.get_accumulations_report_in_period(
            customer_id, start_date, end_date
        )

    async def generate_accumulation(
        self, margin: GasStationMargin, customer: Customer, amount: float
    ):
        # rounding to two decimals
        gas_price = 0
        if margin.margin_type == MarginType.by_margin:
            generated_points = round(amount * (margin.margin / 100) * margin.points, 2)
        else:
            price_key = "regular_price"
            if margin.product.codename == "gas_premium":
                price_key = "premium_price"
            elif margin.product.codename == "gas_diesel":
                price_key = "diesel_price"
            gas_price = getattr(margin.gas_station, price_key)
            total_liters = round(amount / gas_price, 2)
            generated_points = round(total_liters * margin.points, 2)

        accum = Accumulation(
            **margin.model_dump(
                include=["points", "margin", "product", "gas_station", "margin_type"]
            ),
            customer=customer,
            amount=amount,
            gas_price=gas_price,
            generated_points=generated_points
        )
        return await self.__repo.create_accumulation(accum)


accumulation_service = AccumulationService()
