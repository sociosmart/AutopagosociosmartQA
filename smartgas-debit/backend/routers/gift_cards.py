import datetime
from flask import Blueprint, current_app, jsonify, g
from flask_apispec import doc, marshal_with, use_kwargs
from sqlalchemy import except_

from internal.db import db
from internal.dtos.common import AuthorizationHeaders, GeneralMessage, PaginatedResponse
from internal.dtos.gift_cards import (
    GetCardByKeyQuery,
    GiftCardBase,
    GiftCardCreationRequest,
    GiftCardCreationResponse,
    GiftCardDetailByCardKeyResponse,
    GiftCardDetailResponse,
    ListGiftCardsPaginatedRequest,
)
from internal.enums.deposits import MovementType
from internal.exceptions.http import GenericError, InternalServerError, NotFoundError
from internal.middlewares.customer_auth import customer_auth
from internal.models.deposits import Movement
from internal.models.gift_cards import GiftCard
from internal.repositories.gas_stations import (
    get_legal_name_by_id,
    get_station_by_external_id,
)
from internal.repositories.gift_cards import (
    get_gift_card_by_card_key,
    get_gift_card_by_criterias,
    get_gift_cards_by_customer_paginated,
    create_gift_card as create_gift_card_repo,
)
from internal.utils.pagination import make_pagination
from internal.repositories.deposits import create_movement, get_balance_by_customer


bp = Blueprint("gift_cards", __name__)


@bp.get("/")
@customer_auth()
@doc(tags=["Gift Cards"], description="List Gift Cards by customer paginated")
@use_kwargs(ListGiftCardsPaginatedRequest, location="query")
@use_kwargs(AuthorizationHeaders, location="headers")
@marshal_with(PaginatedResponse, code="200")
@marshal_with(GeneralMessage, code="500", description="Internal Server Error")
@marshal_with(GeneralMessage, code="401", description="Unauthorization")
def list_gift_cards(**kwargs):
    customer = g.customer
    try:
        gift_cards = get_gift_cards_by_customer_paginated(customer.id)
    except Exception as e:
        current_app.logger.error(f"Something went wrong while getting gift cards - {e}")
        raise InternalServerError

    return make_pagination(gift_cards, GiftCardBase)


@bp.post("/")
@customer_auth()
@doc(tags=["Gift Cards"], description="Gift Card creation")
@use_kwargs(AuthorizationHeaders, location="headers")
@use_kwargs(
    GiftCardCreationRequest, location="json", description="Gift card Creation body"
)
@marshal_with(GiftCardCreationResponse, code="201")
@marshal_with(GeneralMessage, code="500", description="Internal Server Error")
@marshal_with(GeneralMessage, code="406", description="Unsuficient funds")
@marshal_with(GeneralMessage, code="401", description="Unauthorization")
def create_gift_card(amount, legal_name_id, **kwargs):
    customer = g.customer

    try:
        legal_name = get_legal_name_by_id(legal_name_id)
    except Exception as e:
        current_app.logger.error(f"Something went wrong - {e}")
        raise InternalServerError

    if not legal_name:
        raise NotFoundError(message="Legal Name not found")

    try:
        balance = get_balance_by_customer(customer.id, legal_name_id=legal_name_id)
    except Exception as e:
        current_app.logger.error(f"Something went wrong while bringing balance - {e}")
        raise InternalServerError

    total = balance["total"]

    if amount > total:
        raise GenericError(message="Unsuficient funds", code=406)

    try:
        gift_card = create_gift_card_repo(
            GiftCard(
                customer=customer,
                legal_name=legal_name,
                expiration_date=datetime.datetime.now() + datetime.timedelta(days=15),
                amount=amount,
            ),
            auto_commit=False,
        )

        db.session.flush()
        create_movement(
            Movement(
                customer=customer,
                description=f"Tarjeta de regalo creada en {legal_name.name} por {amount}",
                gift_card=gift_card,
                amount=amount,
                type=MovementType.GIFT_CARD_CREATION,
            ),
            auto_commit=False,
        )
        db.session.commit()
    except Exception as e:
        current_app.logger.error(f"Error while creating gift card - {e}")
        db.session.rollback()
        raise InternalServerError

    return GiftCardCreationResponse().dump(gift_card), 201


@bp.get("/<int:id>")
@customer_auth()
@doc(
    tags=["Gift Cards"],
    description="Get gift card detail",
    params={"id": {"description": "the id of the gift card", "example": 1}},
)
@use_kwargs(AuthorizationHeaders, location="headers")
@marshal_with(GiftCardDetailResponse, code="200", description="Gift Card detail")
@marshal_with(GeneralMessage, code="500", description="Internal Server Error")
@marshal_with(GeneralMessage, code="401", description="Unauthorization")
def get_gift_card_detail(id, **kwargs):
    customer = g.customer
    try:
        gift_card = get_gift_card_by_criterias(id=id, customer_id=customer.id)
    except Exception as e:
        current_app.logger.error(
            f"Something went wrong when bringing gift card detail - {e}"
        )
        raise InternalServerError

    if not gift_card:
        raise NotFoundError("this gift card does not exist for this user")

    return GiftCardDetailResponse().dump(gift_card)


@bp.get("/by-key/<card_key>")
@doc(
    tags=["Gift Cards"],
    description="Get gift card detail by card key",
    params={
        "card_key": {"description": "the card key", "example": "1111-1111-1111-1111"}
    },
)
@marshal_with(
    GiftCardDetailByCardKeyResponse, code="200", description="Gift Card detail"
)
@marshal_with(GeneralMessage, code="500", description="Internal Server Error")
@use_kwargs(GetCardByKeyQuery, location="query")
@marshal_with(
    GeneralMessage,
    code="404",
    description="Card already redeemed or does not exist or expired",
)
@marshal_with(
    GeneralMessage,
    code="409",
    description="Card cannot be redeemed in given gas station",
)
def get_card_detail_by_card_key(card_key, external_gas_station_id, **kwargs):
    try:
        gas_station = get_station_by_external_id(external_gas_station_id)
    except Exception as e:
        current_app.logger.error(
            f"Something went wrong while getting gas station - {e}"
        )
        raise InternalServerError
    else:
        if not gas_station:
            raise NotFoundError
    legal_name_id = gas_station.legal_name_id
    try:
        gift_card = get_gift_card_by_card_key(card_key)
    except Exception as e:
        current_app.logger.error(
            f"Something went wrong when bringing gift card detail - {e}"
        )
        raise InternalServerError

    if not gift_card:
        raise NotFoundError(
            "this gift card does not exist, is already expired or redeemed"
        )

    if gift_card.legal_name_id != legal_name_id:
        raise GenericError(
            message="This giftcard cannot be redemeed in this gas station", code=409
        )

    return GiftCardDetailByCardKeyResponse().dump(gift_card)
