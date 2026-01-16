import datetime
import random

from typing import List, TYPE_CHECKING
from sqlalchemy import Boolean, DateTime, Float, ForeignKey, String, UniqueConstraint
from sqlalchemy.orm import Mapped, mapped_column, relationship
from sqlalchemy.ext.hybrid import hybrid_property
from internal.db import db
from internal.models.customers import Customer

if TYPE_CHECKING:
    from internal.models.deposits import Movement
    from internal.models.gas_stations import LegalName
    from internal.models.payments import Payment
else:
    Movement = "Movement"
    LegalName = "LegalName"
    Payment = "Payment"


def rand_16_digits():
    return "{}-{}-{}-{}".format(
        random.randint(1000, 9999),
        random.randint(1000, 9999),
        random.randint(1000, 9999),
        random.randint(1000, 9999),
    )


def default_card_key():
    card_key = rand_16_digits()
    while (
        db.session.query(GiftCard).where(GiftCard.card_key == card_key).first()
        is not None
    ):
        card_key = rand_16_digits()

    return card_key


class GiftCard(db.Model):
    __tablename__ = "gift_cards"
    __table_args__ = (UniqueConstraint("card_key", name="card_key_unique"),)

    id: Mapped[int] = mapped_column(primary_key=True)
    card_key: Mapped[str] = mapped_column(String(22), default=default_card_key)
    customer_id: Mapped[int] = mapped_column(ForeignKey("customers.id"))
    customer: Mapped[Customer] = relationship(back_populates="gift_cards")
    legal_name_id: Mapped[int] = mapped_column(ForeignKey("legal_names.id"))
    legal_name: Mapped[LegalName] = relationship(back_populates="gift_cards")
    amount: Mapped[float] = mapped_column(Float, default=0)
    amount_used: Mapped[float] = mapped_column(Float, default=0)
    movements: Mapped[List[Movement]] = relationship(back_populates="gift_card")
    redeemed: Mapped[bool] = mapped_column(Boolean, default=False)
    payments: Mapped[List["Payment"]] = relationship(back_populates="gift_card")
    expiration_date: Mapped[datetime.datetime] = mapped_column(
        DateTime, default=datetime.datetime.now
    )
    created_at: Mapped[datetime.datetime] = mapped_column(
        DateTime, default=datetime.datetime.now
    )
    updated_at: Mapped[datetime.datetime] = mapped_column(
        DateTime, onupdate=datetime.datetime.now, default=datetime.datetime.now
    )

    @hybrid_property
    def is_expired(self) -> bool:
        now = datetime.datetime.now()
        return now > self.expiration_date
