from enum import Enum


class BenefitType(str, Enum):
    physical = "physical"
    digital = "digital"
    gas = "gas"
    periferics = "periferics"


class BenefitFrequency(str, Enum):
    n_times = "n_times"
    daily = "daily"
    monthly = "montly"
    weekly = "weekly"
    hourly = "hourly"
    always = "always"


class BenefitGeneratorKeys(str, Enum):
    levels_duration = "levels_duration"
    levels_last_n_days = "levels_last_n_days"
