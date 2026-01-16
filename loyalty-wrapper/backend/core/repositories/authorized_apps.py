from core.exceptions.common import NotFoundException
from core.models.authorized_apps import AuthorizedApp


class AuthorizedAppRepository:
    async def create_authorized_app(self, app: AuthorizedApp) -> AuthorizedApp:
        return await app.create()

    async def get_authorized_app_by_criterias(self, **filters) -> AuthorizedApp:
        app = await AuthorizedApp.find_one(filters)
        if not app:
            raise NotFoundException
        return app


authorized_app_repository = AuthorizedAppRepository()
