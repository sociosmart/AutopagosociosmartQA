import asyncio

from core.db import init_db
from core.models.customers import Customer
from core.utils.pagination import Paginator
from create_customer import create_customer


async def create_users():
    await init_db()  

    for i in range(100_000):
      #  external_id: str, name: str, last_name: str, status: str, phone_number: str, email: str
        external_id = f"external_id_{i}"
        name = "juan pablo"
        last_name = "najera"
        status = "1"
        phone_number = "1234567890"
        email = "jpnajeracastro@hotmail.com"


        await create_customer(external_id, name, last_name, status, phone_number, email)


async def main():
    await init_db()

    cursor = await Paginator(Customer, limit=10).paginate()

    print(cursor.next_cursor)
    print(cursor.prev_cursor)

    print("================== ")



if __name__ == "__main__":
    #asyncio.run(create_users())
    asyncio.run(main())
