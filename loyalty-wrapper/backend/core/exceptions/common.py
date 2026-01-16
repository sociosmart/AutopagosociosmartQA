class BaseErrorException(Exception):
    message: str

    def __init__(self, message):
        self.message = message


class InternalServerError(BaseErrorException):
    def __init__(self):
        self.message = "Internal server error"


class NotFoundException(BaseErrorException):
    def __init__(self, message="Not found"):
        self.message = message


class DuplicatedKeyException(BaseErrorException):
    def __init__(self, message="Duplicated key"):
        self.message = message


class BenefitAlreadyRedeemedException(BaseErrorException):
    def __init__(self, message="Benefit already redeemed"):
        self.message = message


class BenefitNotInDateRangeException(BaseErrorException):
    def __init__(self, message="Benefit not in date range"):
        self.message = message


class BenefitNotStockLeftException(BaseErrorException):
    def __init__(self, message="No more stock for this benefit"):
        self.message = message
