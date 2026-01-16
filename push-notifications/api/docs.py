from apispec import APISpec
from apispec.ext.marshmallow import MarshmallowPlugin
from flask import Flask
from flask_apispec.extension import FlaskApiSpec

from .views import CreateNotificationView


def init_docs(app: Flask):
    spec = APISpec(
        title="API Documentation for smart gas push notifications module",
        version="v1",
        openapi_version="2.0",
        plugins=[MarshmallowPlugin()],
    )
    app.config.update(
        {
            "APISPEC_SPEC": spec,
            "APISPEC_SWAGGER_URL": "/swagger/",
            "APISPEC_SWAGGER_UI_URL": "/swagger-ui/",
        }
    )
    docs = FlaskApiSpec(app)

    docs.register(CreateNotificationView, endpoint="create_notification")
