from flask import Flask
from flask_migrate import Migrate
from flask_cors import CORS

from app.sentry import init_sentry

from .handlers import init_handlers

from .doc import init_docs
from routers import init_routes
from internal.db import db
from internal.cli import init_cli


def create_app():
    init_sentry()
    app = Flask(__name__)
    CORS(app)

    # Global config from env
    app.config.from_pyfile("config.py")

    @app.route("/health")
    def health_check():
        return "healthy"

    init_handlers(app)
    db.init_app(app)
    init_cli(app)
    migrate = Migrate(app, db)
    init_routes(app)
    init_docs(app)

    return app
