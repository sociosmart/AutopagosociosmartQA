import datetime
from typing import List, Optional

from beanie import PydanticObjectId
from beanie.operators import In
from bson import ObjectId

from core.dtos.benefits import UpdateBenefitInDb
from core.dtos.pagination import CursorPage, PaginationParams
from core.enums.benefits import BenefitType
from core.exceptions.common import NotFoundException
from core.models.benefits import Benefit, BenefitGenerated, BenefitTicket
from core.utils.pagination import Paginator


class BenefitRepository:
    async def get_benefits_paginated(
        self, pagination: Optional[PaginationParams], search: Optional[str] = None
    ) -> CursorPage[Benefit]:
        search_query = {}
        if search:
            search_query = {"$or": [{"name": {"$regex": search, "$options": "i"}}]}
        return await Paginator(Benefit, pagination).paginate(
            search_query, fetch_links=True
        )

    async def get_benefits_generated_paginated(
        self, pagination: Optional[PaginationParams], search: Optional[str] = None
    ) -> CursorPage[Benefit]:
        search_query = {}
        if search:
            search_query = {
                "$or": [
                    {
                        "name": {
                            "$regex": search,
                            "$options": "i",
                        }
                    }
                ]
            }
        return await Paginator(BenefitGenerated, pagination).paginate(
            search_query, fetch_links=True
        )

    async def get_benefits_tickets_paginated(
        self, pagination: Optional[PaginationParams], search: Optional[str] = None
    ) -> CursorPage[BenefitTicket]:
        search_query = {}
        if search:
            search_query = {
                "$or": [
                    {"customer.name": {"$regex": search, "$options": "i"}},
                    {"customer.last_name": {"$regex": search, "$options": "i"}},
                    {"benefit_generated.name": {"$regex": search, "$options": "i"}},
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
        return await Paginator(BenefitTicket, pagination).paginate(
            search_query, fetch_links=True, on_demand=False
        )

    async def get_benefit_by_id(self, id: str) -> Benefit:
        benefit = await Benefit.get(id, fetch_links=True)

        if not benefit:
            raise NotFoundException

        return benefit

    async def get_benefit_generated_by_id(self, id: str) -> BenefitGenerated:
        benefit = await BenefitGenerated.get(id, fetch_links=True)

        if not benefit:
            raise NotFoundException

        return benefit

    async def get_benefits(self, **filters):
        return await Benefit.find(filters, fetch_links=True).to_list()

    async def create_benefit(self, benefit: Benefit) -> Benefit:
        benefit = await benefit.create()
        await benefit.fetch_all_links()
        return benefit

    async def update_benefit(self, id: str, data: UpdateBenefitInDb) -> Benefit:
        benefit = await Benefit.get(id)
        if not benefit:
            raise NotFoundException
        new_data = benefit.model_dump() | data.model_dump(exclude_none=True)
        new_model = Benefit(**new_data)
        await new_model.replace()
        await new_model.fetch_all_links()
        return new_model

    async def insert_generated_benefits_many(
        self, generated_benefits: List[BenefitGenerated]
    ):
        await BenefitGenerated.insert_many(generated_benefits)

    async def get_generated_benefits_in_period(
        self, ids: List[PydanticObjectId], date: datetime.datetime
    ):
        return await BenefitGenerated.find(
            In(BenefitGenerated.benefit.id, ids),
            BenefitGenerated.start_date <= date,
            BenefitGenerated.end_date >= date,
            fetch_links=True,
        ).to_list()

    async def get_generated_benefits_in_period_for_level(
        self, level: PydanticObjectId, date: datetime.datetime
    ):
        return await BenefitGenerated.find(
            BenefitGenerated.level.id == level,
            BenefitGenerated.start_date <= date,
            BenefitGenerated.end_date >= date,
            fetch_links=True,
        ).to_list()

    async def insert_many_benefit_tickets(self, tickets: List[BenefitTicket]):
        return await BenefitTicket.insert_many(tickets)

    async def update_generated_benefit(
        self, id: str, data: BenefitGenerated
    ) -> BenefitGenerated:
        benefit = await BenefitGenerated.get(id)
        if not benefit:
            raise NotFoundException("Benefit not found")

        new_data = benefit.model_dump() | data.model_dump(exclude_none=True)
        new_model = BenefitGenerated(**new_data)
        await new_model.replace()
        await new_model.fetch_all_links()
        return new_model

    async def get_active_benefits_by_phone(
        self, phone: str, type_: BenefitType, date: datetime.datetime
    ) -> List[BenefitTicket]:
        return await BenefitTicket.find(
            BenefitTicket.customer.phone_number == phone,
            BenefitTicket.redeemed == False,
            BenefitTicket.start_date <= date,
            BenefitTicket.end_date >= date,
            BenefitTicket.benefit_generated.start_date <= date,
            BenefitTicket.benefit_generated.end_date >= date,
            BenefitTicket.benefit_generated.type == type_,
            BenefitTicket.benefit_generated.is_active == True,
            BenefitTicket.benefit_generated.dependency == False,
            fetch_links=True,
        ).to_list()

    async def get_active_dependent_benefits(
        self, phone: str, charged_amount: float, date: datetime.datetime
    ) -> List[BenefitTicket]:
        return await BenefitTicket.find(
            BenefitTicket.customer.phone_number == phone,
            BenefitTicket.redeemed == False,
            BenefitTicket.start_date <= date,
            BenefitTicket.end_date >= date,
            BenefitTicket.benefit_generated.start_date <= date,
            BenefitTicket.benefit_generated.end_date >= date,
            BenefitTicket.benefit_generated.is_active == True,
            BenefitTicket.benefit_generated.min_amount <= charged_amount,
            BenefitTicket.benefit_generated.dependency == True,
            fetch_links=True,
        ).to_list()

    async def get_gas_discount_by_phone(self, phone: str, date: datetime.datetime):
        benefit_ticket = await BenefitTicket.find_one(
            BenefitTicket.customer.phone_number == phone,
            BenefitTicket.redeemed == False,
            BenefitTicket.start_date <= date,
            BenefitTicket.end_date >= date,
            BenefitTicket.benefit_generated.start_date <= date,
            BenefitTicket.benefit_generated.end_date >= date,
            BenefitTicket.benefit_generated.type == BenefitType.gas,
            BenefitTicket.benefit_generated.is_active == True,
            fetch_links=True,
        )

        if not benefit_ticket:
            raise NotFoundException

        return benefit_ticket

    async def get_benefit_ticket_by_customer(self, id: str, customer_id: str):
        benefit_ticket = await BenefitTicket.find_one(
            BenefitTicket.id == ObjectId(id),
            BenefitTicket.benefit_generated.start_date <= date,
            BenefitTicket.benefit_generated.end_date >= date,
            BenefitTicket.benefit_generated.type != BenefitType.gas,
            BenefitTicket.benefit_generated.type != BenefitType.periferics,
            BenefitTicket.customer.id == ObjectId(customer_id),
            BenefitTicket.benefit_generated.is_active == True,
            fetch_links=True,
        )

        if not benefit_ticket:
            raise NotFoundException("No benefit ticket found")

        return benefit_ticket


benefit_repository = BenefitRepository()
