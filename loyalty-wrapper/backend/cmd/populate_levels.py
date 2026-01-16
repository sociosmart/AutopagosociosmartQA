import asyncio
import os
import sys

sys.path.append(os.path.join(os.path.dirname(__file__), ".."))
from core.db import init_db
from core.models.levels import Level


async def populate_levels():
    await init_db()
    levels = [
        Level(name="Clasico", min_points=0),
        Level(name="Plata", min_points=120),
        Level(name="Oro", min_points=150),
        Level(name="Diamante", min_points=200),
    ]

    await Level.insert_many(levels)


if __name__ == "__main__":
    asyncio.run(populate_levels())
