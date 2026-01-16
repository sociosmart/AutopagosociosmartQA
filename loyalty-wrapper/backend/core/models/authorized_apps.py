import uuid

from beanie import Document
from pydantic import Field

from core.models.mixins import CreatedAtMixin


class AuthorizedApp(Document, CreatedAtMixin):
    name: str = Field(max_length=50)
    app_key: str = Field(default_factory=lambda: str(uuid.uuid4()))
    api_key: str = Field(default_factory=lambda: str(uuid.uuid4()))
    is_active: bool = True

    class Settings:
        name = "authorized_apps"
