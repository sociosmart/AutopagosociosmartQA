import datetime

from sqlalchemy import func, select, update

from internal.models.gift_cards import GiftCard
from internal.db import db


def get_gift_cards_by_customer_paginated(customer_id, filters={}):
    stmt = (
        select(GiftCard)
        .where(GiftCard.customer_id == customer_id)
        .order_by(-GiftCard.created_at)
    )

    return db.paginate(stmt, error_out=False)


def get_commited_funds_stmt(customer_id, legal_name_id=None, date=None):
    if not date:
        date = datetime.datetime.now()

    conditions = []
    if legal_name_id:
        conditions.append(GiftCard.legal_name_id == legal_name_id)

    return select(func.sum(GiftCard.amount)).where(
        GiftCard.customer_id == customer_id,
        GiftCard.redeemed == False,
        GiftCard.expiration_date >= date,
        *conditions
    )


def create_gift_card(gift_card: GiftCard, auto_commit=True) -> GiftCard:
    db.session.add(gift_card)

    if auto_commit:
        db.session.commit()
        db.session.refresh(gift_card)

    return gift_card


def get_gift_card_by_criterias(**kwargs):
    stmt = select(GiftCard).filter_by(**kwargs)

    return db.session.execute(stmt).scalar_one_or_none()


def get_gift_card_by_card_key(card_key):
    stmt = select(GiftCard).where(
        GiftCard.redeemed == False,
        GiftCard.is_expired == False,
        GiftCard.card_key == card_key,
    )

    return db.session.execute(stmt).scalar_one_or_none()


def update_gift_card_by_id(id, values={}, auto_commit=True):
    stmt = update(GiftCard).where(GiftCard.id == id).values(**values)

    db.session.execute(stmt)

    if auto_commit:
        db.session.commit()
