from core.repositories.customers import customer_repository
from core.models.customers import Customer

async def create_customer(external_id: str, name: str, last_name: str, status: str, phone_number: str, email: str):
    try:
        customer = Customer(
            external_id=external_id,
            name=name,
            last_name=last_name,
            status=status,
            phone_number=phone_number,
            email=email,
        )

        await customer_repository.create_customer(customer)

    except Exception as e:
        print("Something went wrong - ", e)