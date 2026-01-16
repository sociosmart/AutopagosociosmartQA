import asyncio
import os
import random
import sys
from string import ascii_lowercase, digits

sys.path.append(os.path.join(os.path.dirname(__file__), ".."))

from core.db import init_db
from core.models.customers import Customer

first_names = [
    "Eduardo",
    "Juan",
    "Pablo",
    "Cuauhtemoc",
    "Mariana",
    "Meredith",
    "Jane",
    "Soe",
    "Marina",
    "Ana",
    "Bronce",
    "Oscar",
    "Cristiano",
]

last_names = [
    "Paez",
    "Najera",
    "Verdugo",
    "Palafox",
    "Gonzalez",
    "Jimenez",
    "Ronaldo",
    "Martinez",
]


async def create_customer(n):
    await init_db()

    customers = []
    for i in range(n):
        email = "{}@{}.com".format(
            "".join(random.choice(ascii_lowercase) for _ in range(20)),
            "".join(random.choice(ascii_lowercase) for _ in range(8)),
        )
        external_id = f"external_id_{i}"
        name = random.choice(first_names)
        last_name = random.choice(last_names)
        status = "1"
        phone_number = "".join(random.choice(digits) for _ in range(10))

        customers.append(
            Customer(
                external_id=external_id,
                name=name,
                last_name=last_name,
                status=status,
                phone_number=phone_number,
                email=email,
            )
        )

    await Customer.insert_many(customers)


if __name__ == "__main__":
    asyncio.run(create_customer(1))
