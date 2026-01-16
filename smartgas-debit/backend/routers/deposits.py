from typing import List
from flask import Blueprint, current_app, g, jsonify
from flask_apispec import doc, marshal_with, use_kwargs

from internal.db import db
from internal.dtos.common import (
    AuthorizationAppHeaders,
    AuthorizationHeaders,
    GeneralMessage,
    PaginatedResponse,
)
from internal.dtos.deposits import (
    DepositListRequest,
    DepositListResponse,
    GetBalanceDetailedResponse,
    GetBalanceQuery,
    GetBalanceResponse,
    GetDepositDetailResponse,
    MakeADepositRequest,
)
from internal.enums.deposits import MovementType
from internal.exceptions.http import InternalServerError, NotFoundError
from internal.middlewares.authorized_app import authorized_app_auth
from internal.middlewares.customer_auth import customer_auth
from internal.models.deposits import Deposit, Movement
from internal.repositories.customers import get_customer_by_external_id
from internal.repositories.deposits import (
    create_deposit,
    create_movement,
    get_balance_detailed_query,
    get_deposit_by_criterias,
    get_deposits_by_customer_paginated,
    get_balance_by_customer,
)
from internal.repositories.gas_stations import (
    get_legal_name_by_id,
    get_station_by_external_id,
)
from internal.services.swit import reserve_funds
from internal.utils.pagination import make_pagination

bp = Blueprint("deposits", __name__)


@bp.get("/")
@customer_auth()
@doc(tags=["Deposits"], description="Listing all deposits per customer")
@use_kwargs(DepositListRequest, location="query")
@use_kwargs(AuthorizationHeaders, location="headers")
@marshal_with(PaginatedResponse, code="200")
@marshal_with(GeneralMessage, code="500", description="Internal Server Error")
@marshal_with(GeneralMessage, code="401", description="Unauthorization")
def list_deposits(**kwargs):
    customer = g.customer
    try:
        deposits = get_deposits_by_customer_paginated(customer.id)
    except Exception as e:
        current_app.logger.error(
            f"Something wrong has happend while trying to get data - {e}"
        )
        raise InternalServerError

    return jsonify(make_pagination(deposits, DepositListResponse))


@bp.post("/")
@customer_auth()
@doc(tags=["Deposits"], description="Make a deposit for pre-paid cash in account")
@use_kwargs(MakeADepositRequest, location="json")
@use_kwargs(AuthorizationHeaders, location="headers")
@marshal_with(GeneralMessage, code="200", description="Deposit added to account")
@marshal_with(GeneralMessage, code="500", description="Internal Server Error")
@marshal_with(GeneralMessage, code="401", description="Unauthorization")
def make_a_deposit(legal_name_id, amount, source_id, cvv, last_4, **kwargs):
    customer = g.customer
    try:
        legal_name = get_legal_name_by_id(legal_name_id)
    except Exception as e:
        current_app.logger.error(f"Something went wrong - {e}")
        raise InternalServerError

    if not legal_name:
        raise NotFoundError(message="Legal Name not found")

    try:
        transaction_id = reserve_funds(
            customer.swit_customer_id,
            source_id,
            cvv,
            last_4,
            amount,
            customer.email,
            f"Deposito a smartgas debito por {amount}",
        )
    except Exception as e:
        return (
            GeneralMessage().dump({"message": "Incorrect cvv or insufficient funds"}),
            402,
        )

    try:
        deposit = create_deposit(
            Deposit(
                amount=amount,
                legal_name=legal_name,
                customer=customer,
                transaction_id=transaction_id,
            ),
            auto_commit=False,
        )
        db.session.flush()
        create_movement(
            Movement(
                customer=customer,
                description=f"Deposito en la razon social {deposit.legal_name.name} por {amount}",
                deposit=deposit,
                amount=amount,
                type=MovementType.DEPOSIT,
            ),
            auto_commit=False,
        )
        db.session.commit()
    except Exception as e:
        current_app.logger.error(f"Somethimg went wrong while creating deposit - {e}")
        db.session.rollback()
        raise InternalServerError

    return jsonify({"message": "ok"})


@bp.get("/balance")
@customer_auth()
@doc(tags=["Deposits"], description="Get balance for customer")
@use_kwargs(AuthorizationHeaders, location="headers")
@use_kwargs(GetBalanceQuery, location="query")
@marshal_with(GetBalanceResponse, code="200", description="Balance for customer")
@marshal_with(GeneralMessage, code="500", description="Internal Server Error")
@marshal_with(GeneralMessage, code="401", description="Unauthorization")
@marshal_with(
    GeneralMessage, code="404", description="Gas Station not found if provided"
)
def get_balance(external_gas_station_id=None, **kwargs):
    customer = g.customer

    legal_name_id = None
    if external_gas_station_id:
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
        balance = get_balance_by_customer(customer.id, legal_name_id)
    except Exception as e:
        current_app.logger.error(f"Something went wrong while getting subtotal - {e}")
        raise InternalServerError

    return jsonify(GetBalanceResponse().dump(balance))


@bp.get("/balance/<external_id>")
@use_kwargs(AuthorizationAppHeaders, location="headers")
@authorized_app_auth()
@doc(tags=["Deposits"], description="Get balance for provided customer")
@marshal_with(GetBalanceResponse, code="200", description="Balance for customer")
@marshal_with(GeneralMessage, code="500", description="Internal Server Error")
@marshal_with(GeneralMessage, code="404", description="Customer not found")
@marshal_with(GeneralMessage, code="401", description="Unauthorization")
def get_balance_by_external_customer(external_id, **kwargs):
    try:
        customer = get_customer_by_external_id(external_id)
    except Exception as e:
        current_app.logger.error(f"Something went wrong while getting customer - {e}")
        raise InternalServerError

    if not customer:
        raise NotFoundError

    try:
        balance = get_balance_by_customer(customer.id)
    except Exception as e:
        current_app.logger.error(f"Something went wrong while getting subtotal - {e}")
        raise InternalServerError

    return jsonify(GetBalanceResponse().dump(balance))


@bp.get("/balance-detailed/<external_id>")
@doc(tags=["Deposits"], description="Get balance for external customer detailed")
@use_kwargs(AuthorizationAppHeaders, location="headers")
@authorized_app_auth()
@marshal_with(
    List[GetBalanceDetailedResponse], code="200", description="Balance for customer"
)
@marshal_with(GeneralMessage, code="500", description="Internal Server Error")
@marshal_with(GeneralMessage, code="401", description="Unauthorization")
@marshal_with(GeneralMessage, code="404", description="Customer not found")
def get_balance_detailed_by_external_customer(external_id, **kwargs):
    try:
        customer = get_customer_by_external_id(external_id)
    except Exception as e:
        current_app.logger.error(f"Something went wrong while getting customer - {e}")
        raise InternalServerError

    if not customer:
        raise NotFoundError

    try:
        balance_detailed = get_balance_detailed_query(customer.id)
    except Exception as e:
        current_app.logger.error(f"Something went wrong while getting subtotal - {e}")
        raise InternalServerError

    return jsonify(GetBalanceDetailedResponse(many=True).load(balance_detailed))


@bp.get("/balance-detailed")
@customer_auth()
@doc(tags=["Deposits"], description="Get balance for customer detailed")
@use_kwargs(AuthorizationHeaders, location="headers")
@marshal_with(
    List[GetBalanceDetailedResponse], code="200", description="Balance for customer"
)
@marshal_with(GeneralMessage, code="500", description="Internal Server Error")
@marshal_with(GeneralMessage, code="401", description="Unauthorization")
def get_balance_detailed(**kwargs):
    customer = g.customer
    try:
        balance_detailed = get_balance_detailed_query(customer.id)
    except Exception as e:
        current_app.logger.error(f"Something went wrong while getting subtotal - {e}")
        raise InternalServerError

    return jsonify(GetBalanceDetailedResponse(many=True).load(balance_detailed))


@bp.get("/<int:id>")
@customer_auth()
@doc(
    tags=["Deposits"],
    description="Get Deposit detail",
    params={"id": {"description": "the id of the deposit", "example": 1}},
)
@use_kwargs(AuthorizationHeaders, location="headers")
@marshal_with(GetDepositDetailResponse, code="200", description="Deposit detail")
@marshal_with(GeneralMessage, code="500", description="Internal Server Error")
@marshal_with(GeneralMessage, code="401", description="Unauthorization")
def get_deposit_detail(id, **kwargs):
    customer = g.customer
    try:
        deposit = get_deposit_by_criterias(id=id, customer_id=customer.id)
    except Exception as e:
        current_app.logger.error(
            f"Something went wrong when bringing deposit detail - {e}"
        )
        raise InternalServerError

    if not deposit:
        raise NotFoundError("this deposit does not exist for this user")

    return GetDepositDetailResponse().dump(deposit)
