import requests
import json
import sentry_sdk
from app.config import SWIT_BUSINESS, SWIT_API_URL, SWIT_TOKEN, SWIT_API_KEY
from internal.exceptions.swit import SwitNotFundsOrIncorrectException


def get_headers():
    return {
        "business": SWIT_BUSINESS,
        "token": SWIT_TOKEN,
        "x-api-key": SWIT_API_KEY,
    }


def reserve_funds(
    swit_customeer_id, source_id, cvv, last_4, amount, email=None, description=""
):
    headers = get_headers()
    data = {
        "customerId": swit_customeer_id,
        "sourceId": source_id,
        "cvv": cvv,
        "amount": amount,
        "capture": True,
        "cardLastDigits": last_4,
        "description": description,
    }
    response = requests.post(
        f"{SWIT_API_URL}/payments", headers=headers, data=json.dumps(data)
    )

    if response.status_code not in [200, 400]:
        scope = sentry_sdk.get_current_scope()
        scope.set_user({"email": email})
        scope.set_level("error")
        scope.set_transaction_name("SwitPaymentError")
        scope.set_context(
            "Swit",
            {
                "url": SWIT_API_URL,
                "status_code": response.status_code,
            },
        )
        scope.capture_message(
            f"Error while reseerving card funds for user user {email} in swit"
        )
        raise Exception(
            f"Error while reseerving card funds for user user {email} in swit"
        )

    response_data = response.json()

    if response.status_code == 400:
        raise SwitNotFundsOrIncorrectException

    result = response_data["result"]

    return result["transactionId"]


def get_cards_by_user(swit_customeer_id, email=None):
    headers = get_headers()
    response = requests.get(
        f"{SWIT_API_URL}/customers/{swit_customeer_id}/cards", headers=headers
    )

    if response.status_code != 200:
        scope = sentry_sdk.get_current_scope()
        scope.set_user({"email": email})
        scope.set_transaction_name("SwitCardList")
        scope.set_context(
            "Swit",
            {
                "url": SWIT_API_URL,
                "status_code": response.status_code,
            },
        )
        scope.capture_message(
            f"Error while retrieving cards for user user {email} in swit"
        )
        raise Exception(
            f"Unable to get cards for user in swit - {response.status_code}"
        )

    response_data = response.json()
    return response_data["result"]


def create_customer(email, first_name, last_name):
    headers = get_headers()
    data = {
        "email": email,
        "first_name": first_name,
        "last_name": last_name,
    }
    response = requests.post(
        f"{SWIT_API_URL}/customers",
        headers=headers,
        data=json.dumps(data),
    )

    if response.status_code != 200:
        scope = sentry_sdk.get_current_scope()
        scope.set_user(data)
        scope.set_transaction_name("SwitUserCreationError")
        scope.set_context(
            "Swit",
            {
                "url": SWIT_API_URL,
                "status_code": response.status_code,
            },
        )
        scope.capture_message(f"Error while registering user {email} in swit")
        raise Exception(f"Unable to register user in swit - {response.status_code}")

    response_data = response.json()

    if response_data["status"] != "Success":
        scope = sentry_sdk.get_current_scope()
        scope.set_user(data)
        scope.set_transaction_name("SwitUserResponseError")
        scope.set_context(
            "Swit",
            {
                "url": SWIT_API_URL,
                **response_data,
            },
        )
        scope.capture_message(f"Error while registering user {email} in swit")
        raise Exception(f"Unable to register user in swit")

    return response_data["result"]
