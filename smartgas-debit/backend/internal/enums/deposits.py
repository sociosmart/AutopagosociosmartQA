from enum import StrEnum, auto


class MovementType(StrEnum):
    DEPOSIT = auto()
    FUNDS_RESERVED = auto()
    FUNDS_CONFIRMATION = auto()
    PAYMENT_CANCELED = auto()
    GIFT_CARD_CREATION = auto()
    GIFT_CARD_REDEMPTION = auto()
