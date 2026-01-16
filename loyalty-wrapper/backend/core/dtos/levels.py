from typing import Optional

from pydantic import BaseModel, Field


class AddLevelBody(BaseModel):
    name: str = Field(max_length=100, description="Level name")
    min_points: float = Field(ge=0, description="Min points to reach this level")
    is_active: bool = Field(description="Level status")


class UpdateLevelBody(BaseModel):
    name: Optional[str] = Field(None, max_length=100, description="Level name")
    min_points: Optional[float] = Field(
        None, ge=0, description="Min points to reach this level"
    )
    is_active: Optional[bool] = Field(None, description="Level status")
