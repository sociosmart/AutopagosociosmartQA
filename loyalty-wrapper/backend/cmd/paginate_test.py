import asyncio

from core.db import init_db
from core.models.users import User
from core.utils.pagination import Paginator
from create_user import create_user


async def create_users():
    for i in range(100_000):
        email = str(i) + "mesterlum@hotmail.com"
        password = email
        await create_user(email, password)


async def main():
    await init_db()

    cursor = await Paginator(User, limit=10).paginate()

    print(cursor.next_cursor)
    print(cursor.prev_cursor)

    print("================== ")


if __name__ == "__main__":
    asyncio.run(create_users())
    #asyncio.run(main())
