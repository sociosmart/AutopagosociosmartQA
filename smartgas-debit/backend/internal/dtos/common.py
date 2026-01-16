from marshmallow import Schema, fields, validate


class GeneralMessage(Schema):
    message = fields.Str()


class PaginatedResponse(Schema):
    total = fields.Int()
    page = fields.Int()
    items = fields.List(fields.Dict())
    next = fields.Int(allow_none=True)
    previous = fields.Int(allow_none=True)


class PaginationRequest(Schema):
    per_page = fields.Int(validate=validate.Range(min=10, max=100), load_default=10)
    page = fields.Int(load_default=1)


class AuthorizationHeaders(Schema):
    authorization = fields.Str(required=True)


class AuthorizationAppHeaders(Schema):
    x_app_key = fields.UUID(
        required=True,
        description="App Key",
        data_key="X-APP-KEY",
    )
    x_api_key = fields.UUID(
        required=True,
        description="Api Key",
        data_key="X-API-KEY",
    )
