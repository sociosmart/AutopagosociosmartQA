import datetime
from typing import List, TYPE_CHECKING

from sqlalchemy.orm import Mapped, mapped_column, relationship
from sqlalchemy import String, DateTime, UniqueConstraint

from internal.db import db

if TYPE_CHECKING:
    from internal.models.deposits import Deposit, Movement
    from internal.models.gift_cards import GiftCard
    from internal.models.payments import Payment
else:
    Deposit = "Deposit"
    Movement = "Movement"
    GiftCard = "GiftCard"
    Payment = "Payment"


class Customer(db.Model):
    __tablename__ = "customers"
    __table_args__ = (UniqueConstraint("external_id", name="external_id_unique"),)

    id: Mapped[int] = mapped_column(primary_key=True)
    external_id: Mapped[str] = mapped_column(String(20), nullable=False)
    swit_customer_id: Mapped[str] = mapped_column(String(100), default="")
    first_name: Mapped[str] = mapped_column(String(255), nullable=False)
    last_name: Mapped[str] = mapped_column(String(255), nullable=False)
    phone_number: Mapped[str] = mapped_column(String(20), nullable=False)
    email: Mapped[str] = mapped_column(String(100), nullable=False)
    deposits: Mapped[List[Deposit]] = relationship(back_populates="customer")
    movements: Mapped[List[Movement]] = relationship(back_populates="customer")
    gift_cards: Mapped[List[GiftCard]] = relationship(back_populates="customer")
    payments: Mapped[List["Payment"]] = relationship(back_populates="customer")
    created_at: Mapped[datetime.datetime] = mapped_column(
        DateTime, default=datetime.datetime.now
    )
    updated_at: Mapped[datetime.datetime] = mapped_column(
        DateTime, onupdate=datetime.datetime.now, default=datetime.datetime.now
    )
