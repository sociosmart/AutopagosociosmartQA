from typing import TYPE_CHECKING, List
import uuid
import datetime

from sqlalchemy import Boolean, DateTime, String
from sqlalchemy.orm import Mapped, mapped_column, relationship
from sqlalchemy.types import UUID
from internal.db import db

if TYPE_CHECKING:
    from internal.models.payments import Payment
else:
    Payment = "Payment"


class Authorization(db.Model):
    __tablename__ = "authorization"

    id: Mapped[int] = mapped_column(primary_key=True)
    app_name: Mapped[str] = mapped_column(String(255))
    app_key: Mapped[uuid.UUID] = mapped_column(UUID(as_uuid=True), default=uuid.uuid4)
    api_key: Mapped[uuid.UUID] = mapped_column(UUID(as_uuid=True), default=uuid.uuid4)
    active: Mapped[Boolean] = mapped_column(Boolean, default=True)
    payments_created: Mapped[List[Payment]] = relationship(
        back_populates="created_by", foreign_keys="Payment.created_by_id"
    )
    payments_updated: Mapped[List[Payment]] = relationship(
        back_populates="updated_by", foreign_keys="Payment.updated_by_id"
    )
    payments_canceled: Mapped[List[Payment]] = relationship(
        back_populates="canceled_by", foreign_keys="Payment.canceled_by_id"
    )
    created_at: Mapped[datetime.datetime] = mapped_column(
        DateTime, default=datetime.datetime.now
    )
    updated_at: Mapped[datetime.datetime] = mapped_column(
        DateTime, onupdate=datetime.datetime.now, default=datetime.datetime.now
    )
