from flask import Flask
from apispec import APISpec
from apispec.ext.marshmallow import MarshmallowPlugin
from flask_apispec.extension import FlaskApiSpec

# registry routes
from routers.customers import (
    get_cards_by_customer,
    get_me_information,
    list_customer_movements,
    list_movements_by_date,
)
from routers.gas_stations import (
    list_gas_stations,
    list_gas_stations_all,
    list_gas_stations_by_legal_name,
    list_groups,
    list_legal_names,
    list_legal_names_all,
)
from routers.deposits import (
    get_balance,
    get_balance_by_external_customer,
    get_balance_detailed,
    get_balance_detailed_by_external_customer,
    get_deposit_detail,
    list_deposits,
    make_a_deposit,
)
from routers.gift_cards import (
    create_gift_card,
    get_card_detail_by_card_key,
    get_gift_card_detail,
    list_gift_cards,
)
from routers.payments import make_a_payment, payment_cancelation, payment_confirmation


def init_docs(app: Flask):
    spec = APISpec(
        title="API Documentation for smart gas debit",
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
    docs.register(list_gas_stations, blueprint="gas-stations")
    docs.register(list_groups, blueprint="gas-stations")
    docs.register(list_legal_names, blueprint="gas-stations")
    docs.register(list_legal_names_all, blueprint="gas-stations")
    docs.register(list_gas_stations_by_legal_name, blueprint="gas-stations")
    docs.register(list_gas_stations_all, blueprint="gas-stations")

    docs.register(list_deposits, blueprint="deposits")
    docs.register(make_a_deposit, blueprint="deposits")
    docs.register(get_balance, blueprint="deposits")
    docs.register(get_balance_detailed, blueprint="deposits")
    docs.register(get_deposit_detail, blueprint="deposits")
    docs.register(get_balance_by_external_customer, blueprint="deposits")
    docs.register(get_balance_detailed_by_external_customer, blueprint="deposits")

    docs.register(list_customer_movements, blueprint="customers")
    docs.register(list_movements_by_date, blueprint="customers")
    docs.register(get_cards_by_customer, blueprint="customers")
    docs.register(get_me_information, blueprint="customers")

    docs.register(list_gift_cards, blueprint="gift_cards")
    docs.register(create_gift_card, blueprint="gift_cards")
    docs.register(get_gift_card_detail, blueprint="gift_cards")
    docs.register(get_card_detail_by_card_key, blueprint="gift_cards")

    docs.register(make_a_payment, blueprint="payments")
    docs.register(payment_confirmation, blueprint="payments")
    docs.register(payment_cancelation, blueprint="payments")
