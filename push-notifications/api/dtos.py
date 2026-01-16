from marshmallow import Schema, ValidationError, fields, validate, validates


class InternalServerError(Schema):
    message = fields.String()


class NotificationLog(Schema):
    error = fields.String()
    message = fields.String()
    platform = fields.String()
    token = fields.String()
    type = fields.String()


class CreateNotificationResponse(Schema):
    count = fields.Integer()
    success = fields.String()
    logs = fields.List(fields.Nested(NotificationLog))


class Notification(Schema):
    notif_id = fields.String(
        required=False,
        metadata={
            "description": "A unique string that identifies the notification for async feedback"
        },
    )
    tokens = fields.List(
        fields.String(), required=True, metadata={"description": "device tokens"}
    )
    platform = fields.Integer(
        required=True,
        metadata={
            "example": "2",
            "description": "1=iOS(Apn), 2=Android/iOS(FCM), 3=Huawei",
        },
        validate=validate.OneOf([1, 2, 3]),
    )
    message = fields.String(
        metadata={"description": "Message for notification", "example": "Hello world!"}
    )
    title = fields.String(
        metadata={"description": "Title for notification", "example": "Title"}
    )
    priority = fields.String(
        metadata={"description": "Sets the priority of the message."},
        validate=validate.OneOf(["normal", "high"]),
    )
    content_available = fields.Bool(
        metadata={"description": "data messages wake the app by default."}
    )
    sound = fields.Dict(
        metadata={
            "description": "sound type.",
        }
    )
    data = fields.Dict(metadata={"description": "Custom data to send in notification."})
    huawei_data = fields.String(metadata={"description": "Huawei extra data"})
    topic = fields.String(
        metadata={"description": "send message to topics, for example, pizza."}
    )

    @validates("tokens")
    def validate_tokens(self, value):
        if not value:
            raise ValidationError("Tokens must have at least one element")


class CreateNotificationRequest(Schema):
    notifications = fields.List(fields.Nested(Notification), required=True)

    @validates("notifications")
    def validate_tokens(self, value):
        if not value:
            raise ValidationError("notifications must have at least one element")
