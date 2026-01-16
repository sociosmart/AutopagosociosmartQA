import uuid

from sqlalchemy import select
from internal.models.authorization import Authorization
from internal.db import db


def get_authorized_app(app_key: uuid.UUID, api_key: uuid.UUID):
    stmt = select(Authorization).where(
        Authorization.app_key == app_key,
        Authorization.api_key == api_key,
        Authorization.active == True,
    )

    return db.session.execute(stmt).scalar_one_or_none()


def create_authorized_app(app: Authorization, auto_commit=True) -> Authorization:
    db.session.add(app)
    if auto_commit:
        db.session.commit()
        db.session.refresh(app)
    return app
