from flask import Flask, jsonify
import sentry_sdk

from internal.exceptions.http import (
    CustomException,
    GenericError,
    InternalServerError,
    NotFoundError,
    UnauthorizedError,
)


def init_handlers(app: Flask):
    @app.errorhandler(422)
    @app.errorhandler(400)
    def handle_error(err):
        headers = err.data.get("headers", None)
        messages = err.data.get("messages", ["Invalid request."])
        if headers:
            return jsonify({"errors": messages}), err.code, headers
        else:
            return jsonify({"errors": messages}), err.code

    @app.errorhandler(CustomException)
    def generic_handle_error(err):
        if isinstance(err, InternalServerError):
            sentry_sdk.capture_exception(err)
            return jsonify({"message": err.message}), err.code
        elif isinstance(err, UnauthorizedError):
            return jsonify({"message": err.message}), err.code
        elif isinstance(err, NotFoundError):
            return jsonify({"message": err.message}), err.code
        elif isinstance(err, GenericError):
            return jsonify({"message": err.message}), err.code

        raise err
