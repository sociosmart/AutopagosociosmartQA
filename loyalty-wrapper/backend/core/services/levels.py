import datetime
import logging
from typing import List, Optional

from bson import ObjectId
from dateutil.relativedelta import relativedelta

from core.dtos.levels import AddLevelBody, UpdateLevelBody
from core.dtos.pagination import CursorPage, PaginationParams
from core.enums.benefits import BenefitFrequency, BenefitGeneratorKeys, BenefitType
from core.exceptions.common import NotFoundException
from core.logging.logger import logger
from core.models.benefits import BenefitGenerated, BenefitTicket
from core.models.customers import Customer
from core.models.levels import CustomerLevel, Level
from core.repositories.benefits import benefit_repository
from core.repositories.customers import customer_repository
from core.repositories.levels import level_repository
from core.services.accumulations import accumulation_service
from core.services.settings import settings_service

logger = logging.getLogger(__name__)


class LevelService:
    def __init__(self):
        self.__repo = level_repository
        self.__customer_repo = customer_repository
        self.__benefit_repo = benefit_repository

    async def get_levels_paginated(
        self, pagination: Optional[PaginationParams], search: Optional[dict] = None
    ) -> CursorPage[Level]:
        return await self.__repo.get_levels_paginated(pagination, search)

    async def get_level_by_id(self, id: str) -> Level:
        return await self.__repo.get_level_by_id(id)

    async def update_level(self, id: str, data: UpdateLevelBody) -> Level:
        return await self.__repo.update_level(id, data)

    async def add_level(self, data: AddLevelBody) -> Level:
        return await self.__repo.add_level(Level(**data.model_dump()))

    async def get_customer_level(self, customer_id: str):
        now = datetime.datetime.now()
        return await self.__repo.get_customer_level(
            customer_id,
            now,
        )

    async def get_customer_level_by_phone(self, phone: str):
        now = datetime.datetime.now()
        return await self.__repo.get_customer_level_by_phone(
            phone,
            now,
        )

    async def get_customer_levels_paginated(
        self, pagination: Optional[PaginationParams]
    ) -> CursorPage[CustomerLevel]:
        return await self.__repo.get_customer_levels_paginated(pagination)

    async def get_suitable_level(self, points: float):
        return await self.__repo.get_suitable_level(points)

    async def create_customer_level(
        self,
        customer_id: str,
        level_id: str,
        start_date: datetime.datetime,
        end_date: datetime.datetime,
    ):
        return await self.__repo.create_customer_level(
            CustomerLevel(
                customer=ObjectId(customer_id),
                level=ObjectId(level_id),
                start_date=start_date,
                end_date=end_date,
            )
        )

    async def __create_benefit_tickets_for_customer(
        self,
        customer: Customer,
        benefit: BenefitGenerated,
    ) -> List[BenefitTicket]:
        start_date = benefit.start_date
        end_date = benefit.end_date
        # Checking types
        # always
        # benefit is gas | periferics | benefit frequency no stock
        benefits_to_add = []
        if (
            benefit.type == BenefitType.gas
            or benefit.type == BenefitType.periferics
            or benefit.frequency == BenefitFrequency.always
        ):
            benefits_to_add.append(
                BenefitTicket(
                    customer=customer,
                    benefit_generated=benefit,
                    start_date=start_date,
                    end_date=end_date,
                ),
            )

        # Applicable benefits X times
        elif benefit.frequency == BenefitFrequency.n_times:
            for _ in range(benefit.num_times):
                benefits_to_add.append(
                    BenefitTicket(
                        customer=customer,
                        benefit_generated=benefit,
                        start_date=start_date,
                        end_date=end_date,
                    ),
                )
        # Applicable benefits every hour
        elif benefit.frequency == BenefitFrequency.hourly:
            # avoding start date renaming
            s_date = start_date
            while s_date < end_date:
                benefits_to_add.append(
                    BenefitTicket(
                        customer=customer,
                        benefit_generated=benefit,
                        start_date=s_date,
                        end_date=s_date + datetime.timedelta(hours=1),
                    )
                )
                s_date = s_date + datetime.timedelta(hours=1)
        # Applicable benefits every day
        elif benefit.frequency == BenefitFrequency.daily:
            # avoding start date renaming
            s_date = start_date
            while s_date < end_date:
                benefits_to_add.append(
                    BenefitTicket(
                        customer=customer,
                        benefit_generated=benefit,
                        start_date=s_date,
                        end_date=s_date + datetime.timedelta(days=1),
                    )
                )
                s_date = s_date + datetime.timedelta(days=1)
        # Applicable benefits every week
        elif benefit.frequency == BenefitFrequency.weekly:
            # avoding start date renaming
            s_date = start_date
            while s_date < end_date:
                benefits_to_add.append(
                    BenefitTicket(
                        customer=customer,
                        benefit_generated=benefit,
                        start_date=s_date,
                        end_date=s_date + datetime.timedelta(days=7),
                    )
                )
                s_date = s_date + datetime.timedelta(days=7)

        # Applicable benefit every month
        elif benefit.frequency == BenefitFrequency.monthly:
            # avoding start date renaming
            s_date = start_date
            while s_date < end_date:
                benefits_to_add.append(
                    BenefitTicket(
                        customer=customer,
                        benefit_generated=benefit,
                        start_date=s_date,
                        end_date=s_date + relativedelta(months=1),
                    )
                )
                s_date = s_date + relativedelta(months=1)

        return benefits_to_add

    async def generate_customers_level(
        self,
    ) -> List[CustomerLevel]:
        # Getting customers to apply the level in period
        customers = await self.__customer_repo.get_customers()

        # Making sure it is the same period
        now = datetime.datetime.now()
        customer_levels = []
        for customer in customers:
            customer_level = await self.generate_customer_level(customer, date=now)
            if customer_level:
                customer_levels.append(customer_level)

        return customer_levels

    async def generate_customer_level(
        self, customer: Customer, date: datetime.datetime, is_new_user=False
    ) -> CustomerLevel | None:
        # Getting active customer level
        try:
            await self.get_customer_level(customer.id)
            return
        except NotFoundException as e:
            pass

        # Getting reported points in period
        points = 0
        if not is_new_user:
            setting = await settings_service.get_by_key(
                BenefitGeneratorKeys.levels_last_n_days
            )

            days = 30

            if setting:
                days = int(setting.value)

            start_date = date - datetime.timedelta(days=days)
            end_date = date
            report = await accumulation_service.get_accumulations_report_in_period(
                customer.id, start_date, end_date
            )
            if report:
                points = report.total

        # Getting suitable level for customer
        suitable_level = await self.get_suitable_level(points)

        # If not suitable level that means we can skip this
        if not suitable_level:
            return

        try:
            applicable_benefits = (
                await self.__benefit_repo.get_generated_benefits_in_period_for_level(
                    suitable_level.id,
                    date,
                )
            )
        except Exception as e:
            logger.error(
                f"Something went wrong while trying to generate benefits for customer - {customer.id} error: {e}"
            )
            return

        if not applicable_benefits:
            return

        try:
            customer_level = await self.create_customer_level(
                customer.id,
                suitable_level.id,
                applicable_benefits[0].start_date,
                applicable_benefits[0].end_date,
            )
        except Exception as e:
            logger.error(f"Something went wrong while creating customer level - {e}")
            # TODO: Something went wrong
            return
        else:
            benefits = []
            for ab in applicable_benefits:
                applicable_tickets = await self.__create_benefit_tickets_for_customer(
                    customer_level.customer,
                    ab,
                )
                benefits.extend(applicable_tickets)

            try:
                await self.__benefit_repo.insert_many_benefit_tickets(benefits)
            except Exception as e:
                logger.error(
                    f"Something went wrong while inserting benefit ticket for customer - {customer_level.customer.id}"
                )

        return customer_level


level_service = LevelService()
