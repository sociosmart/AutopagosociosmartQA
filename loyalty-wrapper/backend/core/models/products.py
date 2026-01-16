from beanie import Document, Indexed
from pydantic import Field


class Product(Document):
    name: str = Field(max_length=100)
    codename: Indexed(str, unique=True)
    is_active: bool = True

    class Settings:
        name = "products"
