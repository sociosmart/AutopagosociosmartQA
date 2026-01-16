from sqlalchemy import and_, func, join, or_, select, update
from sqlalchemy.orm import joinedload
from internal.db import db

from internal.enums.payments import PaymentStatus
from internal.models.deposits import Deposit, Movement
from internal.models.gas_stations import GasStation, LegalName
from internal.models.gift_cards import GiftCard
from internal.models.payments import Payment
from internal.repositories.gift_cards import get_commited_funds_stmt
from internal.repositories.payments import get_reserved_funds_stmt
from internal.utils.db_funcs import subtotal, reserved_total


def get_deposits_by_customer_paginated(customer_id):
    stmt = (
        select(Deposit)
        .where(Deposit.customer_id == customer_id)
        .order_by(-Deposit.created_at)
    )

    return db.paginate(stmt, error_out=False)


def create_deposit(deposit: Deposit, auto_commit=True) -> Deposit:
    db.session.add(deposit)

    if auto_commit:
        db.session.commit()
        db.session.refresh(deposit)

    return deposit


def create_movement(movement: Movement, auto_commit=True) -> Movement:
    db.session.add(movement)

    if auto_commit:
        db.session.commit()
        db.session.refresh(movement)

    return movement


def get_balance_by_customer(
    customer_id,
    legal_name_id=None,
    date=None,
):
    conditions = []
    if legal_name_id:
        conditions.append(Deposit.legal_name_id == legal_name_id)

    commited_funds = func.coalesce(
        (
            get_commited_funds_stmt(customer_id, legal_name_id, date)
            .scalar_subquery()
            .label("commited_funds")
        ),
        0.0,
    )
    reserved_funds = func.coalesce(
        (
            get_reserved_funds_stmt(customer_id, legal_name_id)
            .scalar_subquery()
            .label("reserved_funds")
        ),
        0.0,
    )
    stmt = select(
        subtotal,
        commited_funds,
        reserved_funds,
        (subtotal - commited_funds - reserved_funds).label("total"),
    ).where(
        and_(
            Deposit.customer_id == customer_id,
            *conditions,
        )
    )

    try:
        balance = db.session.execute(stmt).first()
    except Exception as e:
        raise e

    data = dict(zip(["subtotal", "gift_cards", "funds_reserved", "total"], balance))
    return data


def get_balance_detailed_query(customer_id, date=None):
    commited_gift_cards = func.coalesce(
        (
            get_commited_funds_stmt(customer_id, date=date)
            .where(GiftCard.legal_name_id == Deposit.legal_name_id)
            .scalar_subquery()
            .label("commited_funds")
        ),
        0.0,
    )
    reserved_funds = func.coalesce(
        (
            get_reserved_funds_stmt(customer_id)
            .where(Payment.legal_name_id == Deposit.legal_name_id)
            .scalar_subquery()
            .label("reserved_funds")
        ),
        0.0,
    )
    stmt = (
        select(
            LegalName.id,
            LegalName.name,
            subtotal.label("subtotal"),
            commited_gift_cards,
            reserved_funds.label("reserved_funds"),
            (subtotal - commited_gift_cards - reserved_funds).label("total"),
        )
        .join(Deposit)
        .where(Deposit.customer_id == customer_id)
        .group_by(LegalName.id)
    )

    try:
        balance = db.session.execute(stmt).all()
    except Exception as e:
        raise e

    data = list(
        map(
            lambda x: dict(
                zip(
                    [
                        "legal_id",
                        "legal_name",
                        "subtotal",
                        "gift_cards",
                        "funds_reserved",
                        "total",
                    ],
                    x,
                ),
            ),
            balance,
        )
    )

    return data


def get_deposit_by_criterias(**filters):
    stmt = select(Deposit).filter_by(**filters)

    return db.session.execute(stmt).scalar_one_or_none()


def get_movement_by_date_range(customer_id, start_date, end_date, legal_name_id=None):
    extra_filters = []
    if legal_name_id:
        extra_filters.append(
            or_(
                GasStation.legal_name_id == legal_name_id,
                Deposit.legal_name_id == legal_name_id,
                GiftCard.legal_name_id == legal_name_id,
                Payment.legal_name_id == legal_name_id,
            )
        )
    stmt = (
        select(Movement)
        .join(GasStation, GasStation.id == Movement.gas_station_id, isouter=True)
        .join(Deposit, Deposit.id == Movement.deposit_id, isouter=True)
        .join(GiftCard, GiftCard.id == Movement.gift_card_id, isouter=True)
        .join(Payment, Payment.id == Movement.payment_id, isouter=True)
        .filter(
            and_(
                Movement.customer_id == customer_id,
                Movement.created_at.between(start_date, end_date),
            ),
            *extra_filters,
        )
        .order_by(-Movement.created_at)
    )

    return db.session.execute(stmt).scalars()


def get_deposits_with_funds(customer_id, legal_name_id):
    stmt = (
        select(Deposit)
        .where(
            Deposit.customer_id == customer_id,
            Deposit.legal_name_id == legal_name_id,
            Deposit.difference > 0,
        )
        .order_by(Deposit.created_at)
    )

    return db.session.execute(stmt).scalars()


def bulk_deposit_update_by_id(deposits, auto_commit=True):
    stmt = update(Deposit)

    db.session.execute(stmt, deposits)
    if auto_commit:
        db.session.commit()
