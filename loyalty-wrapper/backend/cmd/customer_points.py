import asyncio

from beanie import Link, PydanticObjectId
from pydantic import BaseModel, ConfigDict, Field, ValidationError

from core.db import init_db
from core.dtos.pagination import PaginationParams
from core.models.accumulations import Accumulation, AccumulationReportView
from core.models.customers import Customer
from core.utils.pagination import Paginator


async def main():
    await init_db()
    

    # try:
    #     results = await AccumulationReportView.all(limit=100).to_list()
    #     for r in results:
    #         # await r.fetch_all_links()
    #         print(r)
    # except ValidationError as e:
    #     for v in e.errors():
    #         print(v)

    paginator = await Paginator(AccumulationReportView, PaginationParams(limit=5, prev_cursor=None, next_cursor=None)) \
    .paginate(
        fetch_links=True,
        on_demand=False
        )
    
    for i in paginator.items:
        print(i.id, i.customer.name, i.total_generated_points, i.total_transactions)
    


if __name__ == "__main__":
    asyncio.run(main())
