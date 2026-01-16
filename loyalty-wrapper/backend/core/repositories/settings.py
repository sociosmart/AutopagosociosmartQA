from core.models.settings import Setting


class SettingRepository:
    async def get_setting_by_key(self, key: str) -> Setting | None:
        return await Setting.find_one(Setting.key == key)


setting_repository = SettingRepository()
