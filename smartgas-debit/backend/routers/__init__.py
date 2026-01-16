from flask import Flask
from .customers import bp as customer_bp
from .gas_stations import bp as stations_bp
from .deposits import bp as deposits_bp
from .gift_cards import bp as gift_cards_bp
from .payments import bp as payments_bp


def init_routes(app: Flask):
    common_prefix = "/api/v1"
    app.register_blueprint(customer_bp, url_prefix=f"{common_prefix}/customers")
    app.register_blueprint(stations_bp, url_prefix=f"{common_prefix}/gas-stations")
    app.register_blueprint(deposits_bp, url_prefix=f"{common_prefix}/deposits")
    app.register_blueprint(gift_cards_bp, url_prefix=f"{common_prefix}/gift-cards")
    app.register_blueprint(payments_bp, url_prefix=f"{common_prefix}/payments")
