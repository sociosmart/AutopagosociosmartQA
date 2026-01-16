import os
import sys

sys.path.append(os.path.join(os.path.dirname(__file__), "..", "..", ".."))

from api import create_app
from api.internal.db import db
from api.internal.models.authorization import Authorization

if __name__ == "__main__":
    app = create_app()
    with app.app_context():
        app_name = input("Introduce la app name:")

        auth = Authorization(application_name=app_name)

        db.session.add(auth)
        db.session.commit()
        db.session.refresh(auth)

        print(f"Your apikey={auth.api_key}, appkey={auth.app_key}")
