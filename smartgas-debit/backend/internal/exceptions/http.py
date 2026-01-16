class CustomException(Exception):
    pass


class InternalServerError(CustomException):
    code = 500

    def __init__(self, message="Internal server error"):
        self.message = message


class UnauthorizedError(CustomException):
    code = 401

    def __init__(self, message="Unauthorized"):
        self.message = message


class NotFoundError(CustomException):
    code = 404

    def __init__(self, message="Not found"):
        self.message = message


class GenericError(CustomException):
    def __init__(self, message="Error", code=500):
        self.message = message
        self.code = code
