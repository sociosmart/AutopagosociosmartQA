from core.repositories.settings import setting_repository


class SettingService:
    def __init__(self):
        self.__repo = setting_repository

    async def get_by_key(self, key: str):
        return await self.__repo.get_setting_by_key(key)


settings_service = SettingService()
