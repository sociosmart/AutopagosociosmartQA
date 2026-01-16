from marshmallow import Schema, fields

from internal.dtos.common import PaginationRequest


class GetMeResponse(Schema):
    id = fields.String()
    email = fields.String()
    external_id = fields.String()
    first_name = fields.String()
    last_namee = fields.String()
    phone_number = fields.String()
    swit_customer_id = fields.String()


class CustomerMovementBase(Schema):
    id = fields.Int()
    description = fields.Str()
    type = fields.Str()
    created_at = fields.DateTime()
    amount = fields.Float()


class CardBase(Schema):
    brand = fields.String(default="default")
    last_4 = fields.String()


class PaymentMethod(Schema):
    card = fields.Nested(CardBase)
    id = fields.String()
    is_last_used = fields.Boolean()


class CustomerMovementListRequest(PaginationRequest):
    pass


class CustomerMovementByDateRequest(Schema):
    date_start = fields.Date(
        required=False,
        description="Date used to filter, it must come in yyyy-mm-dd format, used timezone: America/Mazatlan",
        example="2024-08-09",
    )
    date_end = fields.Date(
        required=False,
        description="Date used to filter, it must come in yyyy-mm-dd format, timezone is: America/Mazatlan",
        example="2024-08-09",
    )

    legal_name_id = fields.Int(required=False, description="Legal Name to filter")
