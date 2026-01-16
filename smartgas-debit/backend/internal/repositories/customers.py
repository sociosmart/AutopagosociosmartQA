from sqlalchemy import select, update

from internal.db import db
from internal.models.customers import Customer
from internal.models.deposits import Movement


def get_customer_by_external_id(external_id):
    stmt = select(Customer).where(Customer.external_id == external_id)

    return db.session.execute(stmt).scalar_one_or_none()


def create_customer(cus: Customer) -> Customer:
    db.session.add(cus)
    db.session.commit()
    db.session.refresh(cus)
    return cus


def update_customer_by_id(id, data):
    stmt = update(Customer).where(Customer.id == id).values(**data)

    db.session.execute(stmt)

    return db.session.commit()


def get_paginated_movements_by_customer_id(id):
    stmt = (
        select(Movement)
        .where(Movement.customer_id == id)
        .order_by(Movement.created_at.desc())
    )

    return db.paginate(stmt, error_out=False)
