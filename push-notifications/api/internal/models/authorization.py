import datetime
import uuid

from sqlalchemy import Boolean, DateTime, String
from sqlalchemy.orm import Mapped, mapped_column

from api.internal.db import db


class Authorization(db.Model):
    __tablename__ = "authorization"

    id: Mapped[int] = mapped_column(primary_key=True)
    application_name: Mapped[str] = mapped_column(String(255))
    app_key: Mapped[str] = mapped_column(String(40), unique=True, default=uuid.uuid4)
    api_key: Mapped[str] = mapped_column(String(40), unique=True, default=uuid.uuid4)
    is_active: Mapped[bool] = mapped_column(Boolean(), default=True)
    created_at: Mapped[datetime.datetime] = mapped_column(
        DateTime, default=datetime.datetime.now
    )
