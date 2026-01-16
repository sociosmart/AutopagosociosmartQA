import pymongo
from beanie import Document, Indexed


class Setting(Document):
    key: Indexed(str, pymongo.DESCENDING, unique=True)
    value: str

    class Settings:
        name = "settings"
