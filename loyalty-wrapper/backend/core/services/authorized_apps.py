from core.models.authorized_apps import AuthorizedApp
from core.repositories.authorized_apps import authorized_app_repository


class AuthorizedAppService:
    def __init__(self):
        self.__repo = authorized_app_repository

    async def create_authorized_apps(self, name: str) -> AuthorizedApp:
        return await self.__repo.create_authorized_app(AuthorizedApp(name=name))

    async def get_authorized_app_by_criterias(self, **filters) -> AuthorizedApp:
        return await self.__repo.get_authorized_app_by_criterias(**filters)


authorized_app_service = AuthorizedAppService()
