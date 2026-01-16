from typing import Optional

from beanie import Document, Indexed


class Customer(Document):
    external_id: Indexed(str, unique=True)
    name: str
    last_name: str
    status: str = ""
    phone_number: Indexed(str, unique=True)
    email: Indexed(str, unique=True)
    push_token: str = ""
    token: Optional[str] = ""

    class Settings:
        name = "customers"
