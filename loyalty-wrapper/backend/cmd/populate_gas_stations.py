import asyncio
import os
import sys

sys.path.append(os.path.join(os.path.dirname(__file__), ".."))
from core.db import init_db
from core.services.smart_gas import smart_gas_sevice


async def process_gas_stations():
    try:
        await init_db()

        external_data = await smart_gas_sevice.fetch_external_gas_stations()

        new_stations = await smart_gas_sevice.save_gas_stations(external_data)

        if new_stations:
            print("Nuevas estaciones guardadas")
        else:
            print("No se encontraron estaciones nuevas.")
    except Exception as e:
        print(f"Error durante el procesamiento de estaciones de gasolina: {e}")


if __name__ == "__main__":
    asyncio.run(process_gas_stations())
