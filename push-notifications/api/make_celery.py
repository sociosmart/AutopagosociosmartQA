import os
import sys

sys.path.append(os.path.join(os.path.dirname(__file__), ".."))

from api import create_app

flask_app = create_app()
celery_app = flask_app.extensions["celery"]
