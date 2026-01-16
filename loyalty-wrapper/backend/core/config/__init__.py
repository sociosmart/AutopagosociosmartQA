from ipaddress import IPv4Address

from pydantic import IPvAnyAddress, MongoDsn, RedisDsn
from pydantic_settings import BaseSettings


class Server(BaseSettings):
    host: IPvAnyAddress = IPv4Address("127.0.0.1")
    port: int = 8000
    reload: bool = False


class PushNotifications(BaseSettings):
    push_notifications_url: str = "http://localhost:8009"
    push_app_key: str = ""
    push_api_key: str = ""


class Doc(BaseSettings):
    swagger_url: str = "/swagger"
    redoc_url: str = "/redoc"
    title: str = "Backend Service"
    version: str = "0.0.1"
    summary: str = "A backend service written in Python 3.13 and FastApi."


class DB(BaseSettings):
    connection_uri: MongoDsn = MongoDsn("mongodb://127.0.0.1/")
    collection: str = "loyalty"


class Redis(BaseSettings):
    redis_connection_uri: RedisDsn = RedisDsn("redis://localhost:6379")


class JWT(BaseSettings):
    access_token_expire_minutes: int = 30
    refresh_token_expire_minutes: int = 60 * 24 * 7
    jwt_secret_key: str = (
        "01f663927dba9dc318a4f7569cc78289c2a619ecefc653d1218891d350d2c24c"
    )
    jwt_refresh_secret_key: str = (
        "5fe1514dd87887a48356eb75f6386cd5e71683820a0a0e6c09ae4d0c9eedfe5d"
    )
    algorithm: str = "HS256"


class ThirdPartiesAPI(BaseSettings):
    smartgas_api_url: str = "https://sociosmart.ddns.net"


class Settings(BaseSettings):
    debug: bool = False
    secret: str = "f36e422a6d06ade1406ec409feca98616cd160df11e2bce4409bef7bf4b9bfad"
    server: Server = Server()
    docs: Doc = Doc()
    db: DB = DB()
    jwt: JWT = JWT()
    push: PushNotifications = PushNotifications()
    external_services: ThirdPartiesAPI = ThirdPartiesAPI()
    redis: Redis = Redis()
    tz: str = "America/Mazatlan"


settings = Settings()
