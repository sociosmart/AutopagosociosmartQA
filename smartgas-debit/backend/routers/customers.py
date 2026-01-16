import datetime
from typing import List
from flask import Blueprint, current_app, g, jsonify
from flask_apispec import doc, marshal_with, use_kwargs

from internal.dtos.common import AuthorizationHeaders, GeneralMessage, PaginatedResponse
from internal.dtos.customers import (
    CustomerMovementBase,
    CustomerMovementByDateRequest,
    CustomerMovementListRequest,
    GetMeResponse,
    PaymentMethod,
)
from internal.dtos.deposits import ListMovementsByDateResponse
from internal.exceptions.http import InternalServerError
from internal.middlewares.customer_auth import customer_auth
from internal.repositories.customers import get_paginated_movements_by_customer_id
from internal.repositories.deposits import get_movement_by_date_range
from internal.services.swit import get_cards_by_user
from internal.utils.dates import localize_datetime, localized_date, localized_now
from internal.utils.pagination import make_pagination

bp = Blueprint("customers", __name__)


@bp.get("/me")
@customer_auth()
@doc(tags=["Customers"], description="Get Customer's profile information")
@use_kwargs(AuthorizationHeaders, location="headers")
@marshal_with(GetMeResponse, code="200")
@marshal_with(GeneralMessage, code="500", description="Internal Server Error")
@marshal_with(GeneralMessage, code="401", description="Unauthorization")
def get_me_information(**kwargs):
    customer = g.customer

    return GetMeResponse().dump(customer)


@bp.get("/cards")
@customer_auth()
@doc(tags=["Customers"], description="List all credit cards saved")
@use_kwargs(AuthorizationHeaders, location="headers")
@marshal_with(List[PaymentMethod], code="200")
@marshal_with(GeneralMessage, code="500", description="Internal Server Error")
@marshal_with(GeneralMessage, code="401", description="Unauthorization")
def get_cards_by_customer(**kwargs):
    customer = g.customer
    try:
        cards = get_cards_by_user(customer.swit_customer_id, customer.email)
    except Exception as e:
        current_app.logger.error(f"Unable to bring cardss - {e}")
        raise InternalServerError

    data = []
    for card in cards:
        c = {
            "id": card["card_id"],
            "is_last_used": card["isLastUsed"],
            "card": {
                "last_4": card["last4"],
            },
        }
        data.append(c)

    return PaymentMethod(many=True).dump(data)


@bp.get("/movements")
@customer_auth()
@doc(tags=["Customers"], description="List all movements per customer")
@use_kwargs(CustomerMovementListRequest, location="query")
@use_kwargs(AuthorizationHeaders, location="headers")
@marshal_with(PaginatedResponse, code="200")
@marshal_with(GeneralMessage, code="500", description="Internal Server Error")
@marshal_with(GeneralMessage, code="401", description="Unauthorization")
def list_customer_movements(**kwargs):
    customer = g.customer
    try:
        movements = get_paginated_movements_by_customer_id(customer.id)
    except Exception as e:
        current_app.logger.error(f"Something went wrong - {e}")
        raise InternalServerError

    return jsonify(make_pagination(movements, CustomerMovementBase))


@bp.get("/movements/by-date")
@customer_auth()
@doc(tags=["Customers"], description="Listing movements by date")
@use_kwargs(CustomerMovementByDateRequest, location="query")
@use_kwargs(AuthorizationHeaders, location="headers")
@marshal_with(ListMovementsByDateResponse, code="200")
@marshal_with(GeneralMessage, code="500", description="Internal Server Error")
@marshal_with(GeneralMessage, code="401", description="Unauthorization")
def list_movements_by_date(
    date_start: datetime.date | None = None,
    date_end: datetime.date | None = None,
    legal_name_id=None,
    **kwargs,
):
    customer = g.customer
    max_date_prevention = False
    min_date_prevention = False
    if not date_start and not date_end:
        d = datetime.datetime.now()
        date = localized_date(datetime.datetime(d.year, d.month, d.day))
        date_start = date
        date_end = date
    else:
        if not date_start:
            date_start = datetime.datetime.min
            min_date_prevention = True
        else:
            date_start = localized_date(
                datetime.datetime(date_start.year, date_start.month, date_start.day)
            )
        if not date_end:
            max_date_prevention = True
            date_end = datetime.datetime.max
        else:
            date_end = localized_date(
                datetime.datetime(date_end.year, date_end.month, date_end.day)
            )

    # This is needed in order to get values from range
    if not min_date_prevention:
        min_date = localize_datetime(date_start, tz=datetime.UTC)
    else:
        min_date = date_start

    if not max_date_prevention:
        max_date = (
            localize_datetime(date_end, tz=datetime.UTC)
            + datetime.timedelta(days=1)
            - datetime.timedelta(seconds=1)
        )
    else:
        max_date = date_end

    try:
        movements = get_movement_by_date_range(
            customer.id, min_date, max_date, legal_name_id
        )
    except Exception as e:
        current_app.logger.error(f"Something went wrong - {e}")
        raise InternalServerError

    return ListMovementsByDateResponse(many=True).dump(movements)
