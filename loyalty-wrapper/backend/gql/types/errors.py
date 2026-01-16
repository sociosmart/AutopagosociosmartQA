from typing import List

import strawberry


@strawberry.type
class FieldError:
    field: str = strawberry.field(description="Actual field that contains error")
    message: str = strawberry.field(description="Field error")
    type: str = strawberry.field(description="Actual error code")


@strawberry.type
class InputValidationError:
    errors: List[FieldError] = strawberry.field(description="Validation errors")


@strawberry.type
class GeneralError:
    code: int = strawberry.field(description="HTTP code reference")
    message: str = strawberry.field(description="Actual error message")
