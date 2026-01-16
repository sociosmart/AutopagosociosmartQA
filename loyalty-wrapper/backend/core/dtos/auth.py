from pydantic import BaseModel, EmailStr, Field


class TokenAccess(BaseModel):
    access_token: str
    refresh_token: str


class CreateUserData(BaseModel):
    first_name: str 
    last_name: str
    email: EmailStr
    password: str
