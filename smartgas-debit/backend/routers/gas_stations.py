from typing import List
from flask import Blueprint, jsonify, current_app

from flask_apispec import marshal_with, doc, use_kwargs

from internal.repositories.gas_stations import (
    get_gas_stations,
    get_gas_stations_paginated,
    get_groups_paginated,
    get_legal_name_by_id,
    get_legal_names,
    get_legal_names_paginated,
)
from internal.dtos.common import (
    PaginatedResponse,
    GeneralMessage,
    AuthorizationHeaders,
)
from internal.dtos.gas_stations import (
    GasStationBase,
    GasStationListRequest,
    GasStationGroupBase,
    GasStationsListResponse,
    GroupListRequest,
    LegalNameBase,
    LegalNameListAllRequest,
    LegalNameListRequest,
)
from internal.utils.pagination import make_pagination
from internal.exceptions.http import InternalServerError, NotFoundError
from internal.middlewares.customer_auth import customer_auth

bp = Blueprint("gas-stations", __name__)


@bp.get("/")
@customer_auth()
@doc(tags=["Gas Stations"], description="Listing all gas stations paginated")
@use_kwargs(GasStationListRequest, location="query")
@use_kwargs(AuthorizationHeaders, location="headers")
@marshal_with(PaginatedResponse, code="200")
@marshal_with(GeneralMessage, code="500", description="Internal Server Error")
@marshal_with(GeneralMessage, code="401", description="Unauthorization")
def list_gas_stations(**kwargs):
    try:
        stations_paginated = get_gas_stations_paginated()
    except Exception as e:
        current_app.logger.error(
            f"Something wrong has happend while trying to get data - {e}"
        )
        raise InternalServerError

    return jsonify(make_pagination(stations_paginated, GasStationBase))


@bp.get("/all")
@customer_auth()
@doc(tags=["Gas Stations"], description="Listing all gas stations")
@use_kwargs(AuthorizationHeaders, location="headers")
@marshal_with(GasStationsListResponse, code="200")
@marshal_with(GeneralMessage, code="500", description="Internal Server Error")
@marshal_with(GeneralMessage, code="401", description="Unauthorization")
def list_gas_stations_all(**kwargs):
    try:
        stations = get_gas_stations()
    except Exception as e:
        current_app.logger.error(
            f"Something wrong has happend while trying to get data - {e}"
        )
        raise InternalServerError

    stations_response = GasStationsListResponse(many=True)

    return jsonify(stations_response.dump(stations))


@bp.get("/groups")
@customer_auth()
@doc(tags=["Gas Stations"], description="Listing all groups")
@use_kwargs(AuthorizationHeaders, location="headers")
@use_kwargs(GroupListRequest, location="query")
@marshal_with(PaginatedResponse, code="200")
@marshal_with(GeneralMessage, code="500", description="Internal Server Error")
@marshal_with(GeneralMessage, code="401", description="Unauthorization")
def list_groups(**kwargs):
    try:
        groups = get_groups_paginated()
    except Exception as e:
        current_app.logger.error(
            f"Something wrong has happend while trying to get data - {e}"
        )
        raise InternalServerError

    return jsonify(make_pagination(groups, GasStationGroupBase))


@bp.get("/legal-names")
@customer_auth()
@doc(tags=["Gas Stations"], description="Listing all legal names")
@use_kwargs(AuthorizationHeaders, location="headers")
@use_kwargs(LegalNameListRequest, location="query")
@marshal_with(PaginatedResponse, code="200")
@marshal_with(GeneralMessage, code="500", description="Internal Server Error")
@marshal_with(GeneralMessage, code="401", description="Unauthorization")
def list_legal_names(**kwargs):
    try:
        legal_names = get_legal_names_paginated()
    except Exception as e:
        current_app.logger.error(
            f"Something wrong has happend while trying to get data - {e}"
        )
        raise InternalServerError

    return jsonify(make_pagination(legal_names, LegalNameBase))


@bp.get("/legal-names/all")
@customer_auth()
@doc(tags=["Gas Stations"], description="Listing all in plain format legal names")
@use_kwargs(AuthorizationHeaders, location="headers")
@use_kwargs(LegalNameListAllRequest, location="query")
@marshal_with(List[LegalNameBase], code="200")
@marshal_with(GeneralMessage, code="500", description="Internal Server Error")
@marshal_with(GeneralMessage, code="401", description="Unauthorization")
def list_legal_names_all(**kwargs):
    try:
        legal_names = get_legal_names()
    except Exception as e:
        current_app.logger.error(
            f"something wrong has happend while trying to get data - {e}"
        )
        raise InternalServerError

    legal_names_response = LegalNameBase(many=True)

    return jsonify(legal_names_response.dump(legal_names))


@bp.get("/gas-stations-by-legal-name/<int:id>")
@customer_auth()
@doc(
    tags=["Gas Stations"],
    description="Listing all gas stations per legal name",
    params={"id": {"description": "the id of the legal name", "example": 1}},
)
@use_kwargs(AuthorizationHeaders, location="headers")
@marshal_with(List[GasStationBase], code="200")
@marshal_with(GeneralMessage, code="500", description="Internal Server Error")
@marshal_with(GeneralMessage, code="404", description="Legal Name does not exist")
@marshal_with(GeneralMessage, code="401", description="Unauthorization")
def list_gas_stations_by_legal_name(id, **kwargs):

    try:
        legal_name = get_legal_name_by_id(id)
    except Exception as e:
        current_app.logger.error(f"Something wrong has happened - {e}")
        raise InternalServerError

    if not legal_name:
        raise NotFoundError

    gas_stations = GasStationBase(many=True)
    return jsonify(gas_stations.dump(legal_name.gas_stations))
