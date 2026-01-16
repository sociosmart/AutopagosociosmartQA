from core.exceptions.common import BaseErrorException


class InvalidUserCredentials(BaseErrorException):
    pass


class InactiveUserException(BaseErrorException):
    def __init__(self):
        self.message = "Inactive user"


class UnauthorizationException(BaseErrorException):
    def __init__(self):
        self.message = "Current user not authorized to perform this action"


class InvalidJWTokenException(BaseErrorException):
    def __init__(self):
        self.message = "Token malformed"


class TokenExpiredException(BaseErrorException):
    def __init__(self):
        self.message = "Token Already expired"


class EmailAlreadyTakenException(BaseErrorException):
    def __init__(self):
        self.message = "Email already taken"
