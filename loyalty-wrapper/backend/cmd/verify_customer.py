import asyncio
from core.services.smart_gas import smart_gas_sevice

async def verify_costumer():

    token = input("Paste your token: ")

    try:
        customer_data = await smart_gas_sevice.verify_customer(token)

        print("Customer verified successfully!")
        print(f"Customer Data: {customer_data}")
    except Exception as e:
        print(f"Error verifying customer: {e}")

if __name__ == "__main__":
    asyncio.run(verify_costumer())
