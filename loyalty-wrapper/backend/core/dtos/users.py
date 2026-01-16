from pydantic import BaseModel
from pydantic import Field
from typing import Optional

class CreateUserInDB(BaseModel):
    pass


class AddUserBody(BaseModel):
    first_name: str = Field(max_length=100, description="First name")
    last_name: str = Field(max_length=100, description="Last name")
    email: str = Field(max_length=100, description="Email")
    password: str = Field(max_length=100, description="Password")
    is_active: bool = Field(description="Status")




class UpdateUserBody(BaseModel):
    first_name: Optional[str] = Field(None, description="First name", max_length=100)
    last_name: Optional[str] = Field(None, description="Last name", max_length=100)
    email: Optional[str] = Field(None, description="Email", max_length=100)
    password: Optional[str] = Field(None, description="Password", max_length=100)
    is_active: Optional[bool] = Field(None, description="Status")


