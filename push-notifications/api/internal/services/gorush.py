import requests
from flask import current_app, json


class GorushException(Exception):
    pass


def create_push_notification(data):
    config = current_app.config
    result = requests.post(
        f"{config.get("PUSH_NOTIFICATIONS_API_URL")}/api/push", data=json.dumps(data)
    )

    if result.status_code != 200:
        raise GorushException

    return result.json()
