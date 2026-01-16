import os

DEBUG = int(os.environ.get("DEBUG", 0))
ENVIRONMENT = os.environ.get("ENVIRONMENT", "development")
SECRET_KEY = os.environ.get(
    "SECRET_KEY",
    default="c33b119b6008faa9acc0e0acee41c536b2a087b1e4973f7dc3ca4a65d0d1ca0c",
)
SQLALCHEMY_DATABASE_URI = os.environ.get("SQLALCHEMY_DATABASE_URI")
SMARTGAS_API_URL = os.environ.get("SMARTGAS_API_URL")
SENTRY_ENABLED = int(os.environ.get("SENTRY_ENABLED", 0))
SENTRY_DSN = os.environ.get("SENTRY_DSN")

SWIT_API_URL = os.environ.get("SWIT_API_URL")
SWIT_BUSINESS = os.environ.get("SWIT_BUSINESS")
SWIT_TOKEN = os.environ.get("SWIT_TOKEN")
SWIT_API_KEY = os.environ.get("SWIT_API_KEY")
