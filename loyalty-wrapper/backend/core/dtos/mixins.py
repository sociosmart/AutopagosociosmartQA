import datetime

from pydantic import BaseModel


class WithCreatedAt(BaseModel):
    created_at: datetime.datetime = datetime.datetime.now()


class WithUpdatedAt(BaseModel):
    updated_at: datetime.datetime = datetime.datetime.now()
