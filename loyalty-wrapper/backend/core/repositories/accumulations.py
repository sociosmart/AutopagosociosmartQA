import datetime
from typing import List, Optional

from bson import ObjectId

from core.dtos.pagination import CursorPage, PaginationParams
from core.models.accumulations import (
    Accumulation,
    AccumulationReportView,
    AccumulationsInPeriod,
)
from core.utils.pagination import Paginator


class AccumulationRepository:
    async def get_accumulations_paginated(
        self, pagination: Optional[PaginationParams], search: Optional[str] = None
    ) -> CursorPage[Accumulation]:
        search_query = {}
        if search:
            search_query = {
                "$or": [
                    {"customer.phone_number": {"$regex": search, "$options": "i"}},
                    {"customer.name": {"$regex": search, "$options": "i"}},
                    {"customer.last_name": {"$regex": search, "$options": "i"}},
                    {
                        "$expr": {
                            "$regexMatch": {
                                "input": {
                                    "$concat": [
                                        "$customer.name",
                                        " ",
                                        "$customer.last_name",
                                    ]
                                },
                                "regex": search,
                                "options": "i",
                            }
                        }
                    },
                ]
            }
        return await Paginator(Accumulation, pagination).paginate(  search_query, fetch_links=True, on_demand=False)

    async def get_accumulations_report_paginated(
        self, pagination: Optional[PaginationParams]
    ) -> CursorPage[AccumulationReportView]:
        return await Paginator(AccumulationReportView, pagination).paginate(
            fetch_links=True
        )

    async def get_accumulations_report_in_period(
        self,
        customer_id: str,
        start_date: datetime.datetime,
        end_date: datetime.datetime,
    ) -> AccumulationsInPeriod | None:
        accumulation_report = (
            await Accumulation.find(
                Accumulation.customer.id == ObjectId(customer_id),
                Accumulation.created_at >= start_date,
                Accumulation.created_at <= end_date,
            )
            .aggregate(
                [
                    {
                        "$group": {
                            "_id": "$customer.$id",
                            "total": {"$sum": "$generated_points"},
                            "total_transactions": {"$sum": 1},
                        }
                    }
                ],
                projection_model=AccumulationsInPeriod,
            )
            .to_list()
        )

        return accumulation_report[0] if accumulation_report else None

    async def create_accumulation(self, accum: Accumulation) -> Accumulation:
        accum = await accum.create()
        await accum.fetch_all_links()
        return accum


accumulation_repository = AccumulationRepository()
