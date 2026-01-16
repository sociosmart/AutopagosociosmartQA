from functools import wraps
import uuid

from flask import current_app, request, g

from internal.exceptions.http import InternalServerError, UnauthorizedError
from internal.repositories.authorization import get_authorized_app


def authorized_app_auth():
    def _authorized_app(f):
        @wraps(f)
        def __authorized_app(*args, **kwargs):
            app_key = uuid.UUID(request.headers.get("X-APP-KEY"))
            api_key = uuid.UUID(request.headers.get("X-API-KEY"))

            try:
                authorized_app = get_authorized_app(app_key, api_key)
            except Exception as e:
                current_app.logger.error(
                    f"Something went wrong while getting authorized app credentials - {e}"
                )
                raise InternalServerError

            if not authorized_app:
                raise UnauthorizedError("Unauthorized App")

            g.authorized_app = authorized_app
            result = f(*args, **kwargs)
            return result

        return __authorized_app

    return _authorized_app
