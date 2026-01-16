from typing import Generic, List, Optional, TypeVar

from pydantic import BaseModel, Field

T = TypeVar("T")


class CursorPage(BaseModel, Generic[T]):
    items: List[T]
    next_cursor: Optional[str]
    prev_cursor: Optional[str]
    total: int = 0


class PaginationParams(BaseModel):
    prev_cursor: Optional[str] = Field(None, description="Previous cursor page")
    next_cursor: Optional[str] = Field(None, description="Next Cursor page")
    limit: int = Field(10, gt=0, lt=250)
