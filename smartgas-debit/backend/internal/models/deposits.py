import uuid
import datetime
from typing import TYPE_CHECKING, List, Optional
from sqlalchemy import Boolean, Enum, Float, ForeignKey, DateTime, String
from sqlalchemy.ext.hybrid import hybrid_property
from sqlalchemy.orm import Mapped, mapped_column, relationship

from internal.db import db
from internal.enums.deposits import MovementType
from internal.models.customers import Customer
from internal.models.gas_stations import GasStation, LegalName

if TYPE_CHECKING:
    from internal.models.gift_cards import GiftCard
    from internal.models.payments import PaymentFundsTrack, Payment
else:
    GiftCard = "GiftCard"
    PaymentFundsTrack = "PaymentFundsTrack"
    Payment = "Payment"


class Deposit(db.Model):
    __tablename__ = "deposits"

    id: Mapped[int] = mapped_column(primary_key=True)
    customer_id: Mapped[int] = mapped_column(ForeignKey("customers.id"))
    customer: Mapped[Customer] = relationship(back_populates="deposits")
    legal_name_id: Mapped[int] = mapped_column(ForeignKey("legal_names.id"))
    legal_name: Mapped[LegalName] = relationship(back_populates="deposits")
    amount: Mapped[float] = mapped_column(Float, default=0)
    amount_used: Mapped[float] = mapped_column(Float, default=0)
    is_active: Mapped[Boolean] = mapped_column(Boolean, default=True)
    movements: Mapped[List["Movement"]] = relationship(back_populates="deposit")
    transaction_id: Mapped[str] = mapped_column(String(50), nullable=False, default="")
    payment_funds_track: Mapped[List["PaymentFundsTrack"]] = relationship(
        back_populates="deposit"
    )
    created_at: Mapped[datetime.datetime] = mapped_column(
        DateTime, default=datetime.datetime.now
    )
    updated_at: Mapped[datetime.datetime] = mapped_column(
        DateTime, onupdate=datetime.datetime.now, default=datetime.datetime.now
    )

    @hybrid_property
    def difference(self) -> float:
        return self.amount - self.amount_used


class Movement(db.Model):
    __tablename__ = "movements"

    id: Mapped[int] = mapped_column(primary_key=True)
    description: Mapped[str] = mapped_column(String(500))
    type: Mapped[MovementType] = mapped_column(
        Enum(MovementType), default=MovementType.DEPOSIT
    )
    customer_id: Mapped[int] = mapped_column(ForeignKey("customers.id"))
    customer: Mapped[Customer] = relationship(back_populates="movements")
    gas_station_id: Mapped[Optional[int]] = mapped_column(ForeignKey("gas_stations.id"))
    gas_station: Mapped[Optional[GasStation]] = relationship(back_populates="movements")
    deposit_id: Mapped[Optional[int]] = mapped_column(ForeignKey("deposits.id"))
    deposit: Mapped[Optional[Deposit]] = relationship(back_populates="movements")
    gift_card_id: Mapped[Optional[int]] = mapped_column(ForeignKey("gift_cards.id"))
    gift_card: Mapped[GiftCard] = relationship(back_populates="movements")
    payment_id: Mapped[Optional[uuid.UUID]] = mapped_column(ForeignKey("payments.id"))
    payment: Mapped[Payment] = relationship(back_populates="movements")
    amount: Mapped[float] = mapped_column(Float, default=0)
    is_active: Mapped[Boolean] = mapped_column(Boolean, default=True)
    created_at: Mapped[datetime.datetime] = mapped_column(
        DateTime, default=datetime.datetime.now
    )
    updated_at: Mapped[datetime.datetime] = mapped_column(
        DateTime, onupdate=datetime.datetime.now, default=datetime.datetime.now
    )
