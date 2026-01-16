from marshmallow import Schema, fields, validate
from .common import PaginationRequest


class GasStationBase(Schema):
    id = fields.Int()
    name = fields.Str()
    latitude = fields.Str()
    longitude = fields.Str()


class GasStationGroupBase(Schema):
    id = fields.Int()
    group_name = fields.Str()
    external_id = fields.Str()


class LegalNameBase(Schema):
    id = fields.Int()
    name = fields.Str()


class GasStationListRequest(PaginationRequest):
    name = fields.Str()


class GroupListRequest(PaginationRequest):
    group_name = fields.Str()


class LegalNameListRequest(PaginationRequest):
    name = fields.Str()


class LegalNameListAllRequest(Schema):
    pass


class GasStationsListResponse(GasStationBase):
    legal_name = fields.Nested(LegalNameBase)


class GasStationsByLegalNamePathRequest(Schema):
    id = fields.Int(required=True, validate=validate.Range(min=1), dump_only=True)


class LegalNameWithGasStations(LegalNameBase):
    gas_stations = fields.List(fields.Nested(GasStationBase))
