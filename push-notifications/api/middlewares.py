from functools import wraps

from flask import request
from sqlalchemy import select
from sqlalchemy.exc import NoResultFound

from .internal.db import db
from .internal.models.authorization import Authorization


def authorization(f):
    @wraps(f)
    def __wrapped(*args, **kwargs):
        # just do here everything what you need
        app_key = request.headers.get("X-APP-KEY")
        api_key = request.headers.get("X-API-KEY")

        if not app_key or not api_key:
            return {"message": "Unauthorize to perform this action"}

        stmt = select(Authorization).filter_by(
            app_key=app_key, api_key=api_key, is_active=True
        )

        try:
            db.session.execute(stmt).scalar_one()
        except NoResultFound:
            return {"message": "Unauthorize to perform this action"}

        result = f(*args, **kwargs)
        return result

    return __wrapped
