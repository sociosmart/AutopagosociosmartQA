# TODO: If application gets bigger, consider migrate to a bigger layout for large applications
from flask import Flask, jsonify
from flask_apispec import MethodResource, marshal_with, use_kwargs

from api.internal.services.gorush import GorushException, create_push_notification

from .dtos import (
    CreateNotificationRequest,
    CreateNotificationResponse,
    InternalServerError,
)
from .middlewares import authorization


class CreateNotificationView(MethodResource):
    @authorization
    @use_kwargs(CreateNotificationRequest, location="json")
    @marshal_with(CreateNotificationResponse, code="200")
    @marshal_with(InternalServerError, code="500")
    def post(self, **kwargs):

        try:
            result = create_push_notification(data=kwargs)
        except GorushException:
            return

        return jsonify(result)


def init_views(app: Flask):

    app.add_url_rule(
        "/api/notification",
        view_func=CreateNotificationView.as_view("create_notification"),
    )
