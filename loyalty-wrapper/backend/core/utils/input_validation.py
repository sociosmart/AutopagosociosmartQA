from pydantic import ValidationError

from gql.types.errors import FieldError, InputValidationError


def map_errors(e: ValidationError) -> InputValidationError:
    errors = []
    for error in e.errors():
        errors.append(
            FieldError(field=error["loc"][0], message=error["msg"], type=error["type"])
        )
    return InputValidationError(errors=errors)


def validate_body(input) -> InputValidationError | None:
    if not input:
        return None
    try:
        input.to_pydantic()
    except ValidationError as e:
        return map_errors(e)

    return None
