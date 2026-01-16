from sqlalchemy import func, insert, or_, select, update
from internal.db import db
from internal.enums.payments import PaymentStatus
from internal.models.payments import Payment, PaymentFundsTrack


def create_payment(payment: Payment, auto_commit=True) -> Payment:
    db.session.add(payment)

    if auto_commit:
        db.session.commit()
        db.session.refresh(payment)

    return payment


def check_if_gift_card_is_already_used(gift_card_id):
    stmt = select(Payment).where(
        Payment.gift_card_id == gift_card_id,
        or_(
            Payment.status == PaymentStatus.FUNDS_RESERVED,
            Payment.status == PaymentStatus.CONFIRMED,
        ),
    )

    return db.session.execute(stmt).scalar_one_or_none()


def get_reserved_funds_stmt(customer_id, legal_name_id=None):
    extra_conditions = []
    if legal_name_id:
        extra_conditions.append(Payment.legal_name_id == legal_name_id)
    return select(func.sum(Payment.amount)).where(
        Payment.customer_id == customer_id,
        Payment.status == PaymentStatus.FUNDS_RESERVED,
        *extra_conditions
    )


def get_payment_by_id(id):
    stmt = select(Payment).where(Payment.id == id)

    return db.session.execute(stmt).scalar_one_or_none()


def update_payment_by_id(id, values={}, auto_commit=True):
    stmt = update(Payment).where(Payment.id == id).values(**values)

    db.session.execute(stmt)
    if auto_commit:
        db.session.commit()


def bulk_create_funds_track(data, auto_commit=True):
    stmt = insert(PaymentFundsTrack)

    db.session.execute(stmt, data)

    if auto_commit:
        db.session.commit()
