import asyncio
import os
import sys

sys.path.append(os.path.join(os.path.dirname(__file__), ".."))
from core.db import init_db
from core.dtos.auth import CreateUserData
from core.exceptions.auth import EmailAlreadyTakenException
from core.services.auth import auth_service


async def create_user(first_name: str, last_name: str, email: str, password: str):
    await init_db()
    try:
        await auth_service.create_user(
            CreateUserData(
                first_name=first_name,
                last_name=last_name,
                email=email,
                password=password,
            ),
        )
    except EmailAlreadyTakenException:
        print("Email already taken")
    except Exception as e:
        print("Something went wrong - ", e)


if __name__ == "__main__":
    first_name = input("Firstname: ")
    last_name = input("Lastname: ")
    email = input("Email: ")
    password = input("Password: ")

    asyncio.run(create_user(first_name, last_name, email, password))
