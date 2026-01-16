import datetime

from typing import List, TYPE_CHECKING, Optional

from sqlalchemy import ForeignKey, String, DateTime, UniqueConstraint
from sqlalchemy.orm import Mapped, mapped_column, relationship

from internal.db import db
from internal.models.gift_cards import GiftCard

if TYPE_CHECKING:
    from internal.models.deposits import Deposit, Movement
    from internal.models.payments import Payment
else:
    Deposit = "Deposit"
    Movement = "Movement"
    Payment = "Payment"


class GasStationGroup(db.Model):
    __tablename__ = "gas_station_groups"
    __table_args__ = (UniqueConstraint("external_id", name="external_id_unique"),)

    id: Mapped[int] = mapped_column(primary_key=True)
    external_id: Mapped[str] = mapped_column(String(10), nullable=False)
    group_name: Mapped[str] = mapped_column(String(255), nullable=False, default="")
    gas_stations: Mapped[List["GasStation"]] = relationship(back_populates="group")
    created_at: Mapped[datetime.datetime] = mapped_column(
        DateTime, default=datetime.datetime.now
    )
    updated_at: Mapped[datetime.datetime] = mapped_column(
        DateTime, onupdate=datetime.datetime.now, default=datetime.datetime.now
    )


class LegalName(db.Model):
    __tablename__ = "legal_names"
    __table_args__ = (UniqueConstraint("external_id", name="external_id_unique"),)

    id: Mapped[int] = mapped_column(primary_key=True)
    external_id: Mapped[str] = mapped_column(String(10))
    name: Mapped[str] = mapped_column(String(255), default="")
    gas_stations: Mapped[List["GasStation"]] = relationship(back_populates="legal_name")
    deposits: Mapped[List[Deposit]] = relationship(back_populates="legal_name")
    gift_cards: Mapped[List[GiftCard]] = relationship(back_populates="legal_name")
    payments: Mapped[Optional[Payment]] = relationship(back_populates="legal_name")
    created_at: Mapped[datetime.datetime] = mapped_column(
        DateTime, default=datetime.datetime.now
    )
    updated_at: Mapped[datetime.datetime] = mapped_column(
        DateTime, onupdate=datetime.datetime.now, default=datetime.datetime.now
    )


class GasStation(db.Model):
    __tablename__ = "gas_stations"
    __table_args__ = (UniqueConstraint("external_id", name="external_id_unique"),)

    id: Mapped[int] = mapped_column(primary_key=True)
    name: Mapped[str] = mapped_column(String(255), nullable=False)
    external_id: Mapped[str] = mapped_column(String(10), nullable=False)
    group_id: Mapped[Optional[int]] = mapped_column(ForeignKey("gas_station_groups.id"))
    group: Mapped[Optional[GasStationGroup]] = relationship(
        back_populates="gas_stations"
    )
    legal_name_id: Mapped[Optional[int]] = mapped_column(ForeignKey("legal_names.id"))
    legal_name: Mapped[Optional[LegalName]] = relationship(
        back_populates="gas_stations"
    )
    cre_permission: Mapped[str] = mapped_column(String(255), nullable=False)
    latitude: Mapped[str] = mapped_column(String(20), default="")
    longitude: Mapped[str] = mapped_column(String(20), default="")
    movements: Mapped[List[Movement]] = relationship(back_populates="gas_station")
    created_at: Mapped[datetime.datetime] = mapped_column(
        DateTime, default=datetime.datetime.now
    )
    updated_at: Mapped[datetime.datetime] = mapped_column(
        DateTime, onupdate=datetime.datetime.now, default=datetime.datetime.now
    )
