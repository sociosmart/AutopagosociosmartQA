import datetime
from typing import Optional

from beanie import Replace, before_event
from pydantic import BaseModel, Field


class CreatedAtMixin(BaseModel):
    created_at: datetime.datetime = Field(default_factory=datetime.datetime.now)


class UpdatedAtMixin(BaseModel):
    updated_at: Optional[datetime.datetime] = None

    @before_event(Replace)
    def set_updated_at(self):
        self.updated_at = datetime.datetime.now()
