import os

from flask import Flask
from flask_migrate import Migrate

from .celery import celery_init_app
from .docs import init_docs
from .handlers import init_handlers
from .internal.db import db

# Adding models, so script can be able to migrate them
from .internal.models import *
from .views import init_views


def create_app():
    app = Flask(__name__)

    # config
    app.config.from_pyfile("config.py")
    app.config.from_mapping(
        CELERY={
            "broker_url": os.environ.get("REDIS_CONN"),
            "result_backend": os.environ.get("REDIS_CONN"),
            "task_ignore_result": True,
            "broker_connection_retry_on_startup": True,
            "beat_schedule": {
                "check-notifications": {
                    "task": "api.tasks.check_notifications",
                    "schedule": 10,
                }
            },
        }
    )

    @app.route("/healthz")
    def healthz():
        return "ok"

    db.init_app(app)
    celery_init_app(app)

    # Migrations
    migrate = Migrate()
    migrate.init_app(app, db)

    init_handlers(app)
    init_views(app)
    init_docs(app)

    return app
