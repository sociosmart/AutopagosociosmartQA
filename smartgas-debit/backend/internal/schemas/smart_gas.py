from marshmallow import Schema, fields


class GasStationSchema(Schema):
    name = fields.Str()
    external_id = fields.Str()
    external_group_id = fields.Str()
    external_legal_name_id = fields.Str()
    cre_permission = fields.Str()
    latitude = fields.Str()
    longitude = fields.Str()


class GasStationGroupSchema(Schema):
    group_name = fields.Str()
    external_id = fields.Str()


class LegalNameSchema(Schema):
    name = fields.Str()
    external_id = fields.Str()
