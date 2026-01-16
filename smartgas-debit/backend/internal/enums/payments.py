from enum import StrEnum, auto


class PaymentStatus(StrEnum):
    FUNDS_RESERVED = auto()
    CONFIRMED = auto()
    CANCELLED = auto()
