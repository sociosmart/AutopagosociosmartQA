import datetime
from typing import List, Optional

from beanie import PydanticObjectId, WriteRules

from core.dtos.benefits import (
    AddBenefit,
    UpdateBenefit,
    UpdateBenefitInDb,
    UpdateGeneratedBenefit,
)
from core.dtos.pagination import PaginationParams
from core.enums.benefits import BenefitFrequency, BenefitType
from core.exceptions.common import (
    BenefitAlreadyRedeemedException,
    BenefitNotInDateRangeException,
    BenefitNotStockLeftException,
    NotFoundException,
)
from core.models.benefits import Benefit, BenefitGenerated
from core.repositories.benefits import benefit_repository
from core.services.levels import level_service


class BenefitService:
    def __init__(self):
        self.__repo = benefit_repository
        self.__level_service = level_service

    async def get_benefits_paginated(
        self, pagination: Optional[PaginationParams], search: Optional[dict] = None
    ):
        return await self.__repo.get_benefits_paginated(pagination, search)

    async def get_benefits(self, **filters):
        return await self.__repo.get_benefits(**filters)

    async def get_benefits_generated_paginated(
        self, pagination: Optional[PaginationParams], search: Optional[dict] = None
    ):
        return await self.__repo.get_benefits_generated_paginated(pagination, search)

    async def get_benefits_ticket_paginated(
        self, pagination: Optional[PaginationParams], search: Optional[dict] = None
    ):
        return await self.__repo.get_benefits_tickets_paginated(pagination, search)

    async def get_benefit_by_id(self, id):
        return await self.__repo.get_benefit_by_id(id)

    async def get_benefit_generated_by_id(self, id):
        return await self.__repo.get_benefit_generated_by_id(id)

    async def create_benefit(self, data: AddBenefit):
        try:
            await self.__level_service.get_level_by_id(data.level)
        except NotFoundException:
            raise NotFoundException(message="Level not found")
        except Exception as e:
            raise e

        return await self.__repo.create_benefit(Benefit(**data.model_dump()))

    async def get_generated_benefits_in_period(
        self, ids: List[PydanticObjectId], date: datetime.datetime
    ) -> List[BenefitGenerated]:
        return await self.__repo.get_generated_benefits_in_period(ids, date)

    async def get_generated_benefits_in_period_for_level(
        self, level: PydanticObjectId, date: datetime.datetime
    ):
        return await self.__repo.get_generated_benefits_in_period_for_level(level, date)

    async def insert_generated_benefits_many(
        self, generated_benefit: List[BenefitGenerated]
    ):
        await self.__repo.insert_generated_benefits_many(generated_benefit)

    async def update_benefit(self, id: str, data: UpdateBenefit):
        # Validate level
        if data.level:
            try:
                await self.__level_service.get_level_by_id(data.level)
            except NotFoundException:
                raise NotFoundException(message="Level not found")
            except Exception as e:
                raise e
        return await self.__repo.update_benefit(
            id, UpdateBenefitInDb(**data.model_dump())
        )

    async def update_generated_benefit(
        self, id: str, data: UpdateGeneratedBenefit
    ) -> BenefitGenerated:
        try:
            return await self.__repo.update_generated_benefit(id, data)
        except Exception as e:
            raise e

    async def get_active_benefits_by_phone(
        self, phone: str, type_: BenefitType = BenefitType.digital
    ):
        return await self.__repo.get_active_benefits_by_phone(
            phone, type_, datetime.datetime.now()
        )

    async def get_gas_discount_by_phone(self, phone: str):
        return await self.__repo.get_gas_discount_by_phone(
            phone, datetime.datetime.now()
        )

    async def redeem_benefit(self, id: str, customer: str):
        benefit_ticket = await self.__repo.get_benefit_ticket_by_customer(id, customer)

        now = datetime.datetime.now()
        # Benefit already redeemed
        if benefit_ticket.redeemed:
            raise BenefitAlreadyRedeemedException

        # Not in data range
        if not (benefit_ticket.start_date <= now and benefit_ticket.end_date >= now):
            raise BenefitNotInDateRangeException

        # If benefit has -1 in stock that means that there is no limit
        if benefit_ticket.benefit_generated.stock > -1:
            # Benefit already used
            if (
                benefit_ticket.benefit_generated.stock_used
                >= benefit_ticket.benefit_generated.stock
            ):
                raise BenefitNotStockLeftException

        # changing values
        if not benefit_ticket.benefit_generated.frequency == BenefitFrequency.always:
            benefit_ticket.redeemed = True
        benefit_ticket.benefit_generated.stock_used += 1
        await benefit_ticket.save(link_rule=WriteRules.WRITE)

    async def get_active_dependent_benefits(self, phone: str, charged_amount: float):
        return await self.__repo.get_active_dependent_benefits(
            phone, charged_amount, datetime.datetime.now()
        )


benefit_service = BenefitService()
