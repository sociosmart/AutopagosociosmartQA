from marshmallow import Schema, fields
from webargs import validate
from .common import PaginationRequest
from .gas_stations import GasStationBase, LegalNameBase, LegalNameWithGasStations
from .gift_cards import GiftCardBase


class DepositBase(Schema):
    id = fields.Int()


class DepositListRequest(PaginationRequest):
    name = fields.Str()


class GetBalanceQuery(Schema):
    external_gas_station_id = fields.String()


class DepositListResponse(Schema):
    id = fields.Int()
    amount = fields.Float()
    amount_used = fields.Float()
    created_at = fields.DateTime()
    legal_name = fields.Nested(LegalNameBase)
    created_at = fields.DateTime()


class MakeADepositRequest(Schema):
    legal_name_id = fields.Int(
        validate=validate.Range(min=1),
        description="Legal Name to link the deposit",
        required=True,
    )
    amount = fields.Float(
        validate=validate.Range(min=10),
        description="Amount to deposit",
        required=True,
    )
    source_id = fields.String(
        required=True,
        description="The actual card that is going to be used: bd660382-a226-493a-8541-0bc0e0b12345",
    )
    cvv = fields.String(
        required=True,
        description="The actual cvv of the card",
    )
    last_4 = fields.String(
        required=True,
        description="Last 4 digits for credit card",
        validate=validate.Length(min=4, max=4),
    )


class GetBalanceResponse(Schema):
    subtotal = fields.Float()
    total = fields.Float()
    gift_cards = fields.Float()
    funds_reserved = fields.Float()


class GetBalanceDetailedResponse(Schema):
    legal_id = fields.Int()
    legal_name = fields.Str()
    subtotal = fields.Float()
    total = fields.Float()
    gift_cards = fields.Float()
    funds_reserved = fields.Float()


class GetDepositDetailResponse(Schema):
    amount = fields.Float()
    amount_used = fields.Float()
    created_at = fields.DateTime()
    legal_name = fields.Nested(LegalNameWithGasStations)
    created_at = fields.DateTime()


class ListMovementsByDateResponse(Schema):
    id = fields.Int()
    amount = fields.Float()
    description = fields.String()
    type = fields.String()
    created_at = fields.DateTime()
    deposit_id = fields.Int()
    gift_card_id = fields.Int()
    gas_station_id = fields.Int()
