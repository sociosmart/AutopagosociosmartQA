# from core.utils.bson import PyObjectId
from beanie import Document, Indexed
from pydantic import EmailStr, Field


# Default User model, this will be the 'standard' model to share across multiple layers
class User(Document):
    first_name: str
    last_name: str
    email: Indexed(EmailStr, unique=True)
    password: str = Field(exclude=True)
    is_active: bool = True

    class Settings:
        name = "users"
