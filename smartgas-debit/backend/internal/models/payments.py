import datetime
from typing import List, Optional
import uuid
from sqlalchemy import DateTime, Float, ForeignKey, Enum
from sqlalchemy.ext.hybrid import hybrid_property
from sqlalchemy.orm import Mapped, mapped_column, relationship
from sqlalchemy.types import UUID
from internal.db import db
from internal.enums.payments import PaymentStatus
from internal.models.customers import Customer
from internal.models.deposits import Deposit, Movement
from internal.models.gas_stations import LegalName
from internal.models.gift_cards import GiftCard
from internal.models.authorization import Authorization


class Payment(db.Model):
    __tablename__ = "payments"

    id: Mapped[uuid.UUID] = mapped_column(
        UUID(as_uuid=True),
        default=uuid.uuid4,
        primary_key=True,
    )
    legal_name_id: Mapped[int] = mapped_column(ForeignKey("legal_names.id"))
    legal_name: Mapped[Optional[LegalName]] = relationship(back_populates="payments")
    customer_id: Mapped[Optional[int]] = mapped_column(ForeignKey("customers.id"))
    customer: Mapped[Optional[Customer]] = relationship(back_populates="payments")
    gift_card_id: Mapped[Optional[int]] = mapped_column(ForeignKey("gift_cards.id"))
    gift_card: Mapped[Optional[GiftCard]] = relationship(back_populates="payments")
    amount: Mapped[float] = mapped_column(Float, default=0)
    amount_confirmed: Mapped[float] = mapped_column(Float, default=0)
    status: Mapped[PaymentStatus] = mapped_column(
        Enum(PaymentStatus), default=PaymentStatus.FUNDS_RESERVED
    )
    payment_funds_track: Mapped[List["PaymentFundsTrack"]] = relationship(
        back_populates="payment"
    )
    created_by_id: Mapped[Optional[int]] = mapped_column(ForeignKey("authorization.id"))
    created_by: Mapped[Authorization] = relationship(
        back_populates="payments_created",
        foreign_keys=[created_by_id],
    )
    updated_by_id: Mapped[Optional[int]] = mapped_column(
        ForeignKey("authorization.id"),
    )
    updated_by: Mapped[Optional[Authorization]] = relationship(
        back_populates="payments_updated",
        foreign_keys=[updated_by_id],
    )
    canceled_by_id: Mapped[Optional[int]] = mapped_column(
        ForeignKey("authorization.id"),
    )
    canceled_by: Mapped[Optional[Authorization]] = relationship(
        back_populates="payments_canceled",
        foreign_keys=[canceled_by_id],
    )
    movements: Mapped[Movement] = relationship(back_populates="payment")
    created_at: Mapped[datetime.datetime] = mapped_column(
        DateTime, default=datetime.datetime.now
    )
    updated_at: Mapped[datetime.datetime] = mapped_column(
        DateTime, onupdate=datetime.datetime.now, default=datetime.datetime.now
    )

    @hybrid_property
    def paid_by_customer(self):
        return not not self.customer_id


class PaymentFundsTrack(db.Model):
    __tablename__ = "payment_funds_track"

    id: Mapped[int] = mapped_column(primary_key=True)
    payment_id: Mapped[uuid.UUID] = mapped_column(ForeignKey("payments.id"))
    payment: Mapped[Payment] = relationship(back_populates="payment_funds_track")
    deposit_id: Mapped[int] = mapped_column(ForeignKey("deposits.id"))
    deposit: Mapped[Deposit] = relationship(back_populates="payment_funds_track")
    amount: Mapped[float] = mapped_column(Float, default=0)
    prev_amount: Mapped[float] = mapped_column(Float, default=0)
    created_at: Mapped[datetime.datetime] = mapped_column(
        DateTime, default=datetime.datetime.now
    )
