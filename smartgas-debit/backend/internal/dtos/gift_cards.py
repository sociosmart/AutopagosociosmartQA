from marshmallow import Schema, fields, validate
from .common import PaginationRequest
from .gas_stations import LegalNameBase, LegalNameWithGasStations


class GiftCardBase(Schema):
    id = fields.Int()
    card_key = fields.String()
    amount = fields.Float()
    amount_used = fields.Float()
    expiration_date = fields.DateTime()
    redeemed = fields.Boolean()
    legal_name = fields.Nested(LegalNameBase)
    is_expired = fields.Boolean()


class ListGiftCardsPaginatedRequest(PaginationRequest):
    pass


class GiftCardCreationRequest(Schema):
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


class GiftCardCreationResponse(Schema):
    id = fields.Int()
    card_key = fields.String()


class GiftCardDetailResponse(GiftCardBase):
    legal_name = fields.Nested(LegalNameWithGasStations)


class GiftCardDetailByCardKeyResponse(Schema):
    amount = fields.Float()


class GetCardByKeyQuery(Schema):
    external_gas_station_id = fields.String(required=True)
