from marshmallow import Schema, ValidationError, fields, validate, validates_schema


class MakeAPaymentRequest(Schema):
    amount = fields.Float(
        required=True,
        description="The amount for payment",
        example=10,
        validate=validate.Range(min=1),
    )
    card_key = fields.String(
        required=False,
        description="The gift card key",
        examples=["1111-1111-1111-1111"],
    )
    external_customer_id = fields.String(
        required=False,
        description="If wanted to charge a customer, then provide this field.",
        examples=["1"],
    )

    external_legal_name_id = fields.String(
        required=True,
        description="Legal Name where to apply payment",
        examples=["1"],
    )

    @validates_schema
    def valida_card_key_or_customer_id(self, data, **kwargs):
        if not data.get("card_key") and not data.get("external_customer_id"):
            raise ValidationError(
                "you must set at least card_key or customer_id parameter"
            )


class MakeAPaymentResponse(Schema):
    id = fields.UUID()


class PaymentConfirmationPathParams(Schema):
    id = fields.UUID(required=True)


class PaymentConfirmationRequest(Schema):
    amount = fields.Float(
        validate=validate.Range(min=0.1),
        description="The amount to confirm in the payment",
    )
