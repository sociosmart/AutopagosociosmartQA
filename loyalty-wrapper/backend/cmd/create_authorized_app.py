import asyncio
import os
import sys

sys.path.append(os.path.join(os.path.dirname(__file__), ".."))

from core.db import init_db
from core.services.authorized_apps import authorized_app_service


async def main():
    await init_db()

    app = await authorized_app_service.create_authorized_apps(input("Authorized App: "))

    print("App Key: ", app.app_key, app.api_key)


if __name__ == "__main__":
    asyncio.run(main())
