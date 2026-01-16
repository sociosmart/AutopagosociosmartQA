import os

SECRET_KEY = os.environ.get(
    "SECRET_KEY",
    default="c33b119b6008faa9acc0e0acee41c536b2a087b1e4973f7dc3ca4a65d0d1ca0c",
)

PUSH_NOTIFICATIONS_API_URL = os.environ.get("PUSH_NOTIFICATIONS_API_URL")
SQLALCHEMY_DATABASE_URI = os.environ.get("SQLALCHEMY_DATABASE_URI")
SMARTGAS_API_URL = os.environ.get("SMARTGAS_API_URL")
