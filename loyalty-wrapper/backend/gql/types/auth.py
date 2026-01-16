from typing import Annotated
import strawberry


@strawberry.type
class LoginSuccess:
    access_token: str
    refresh_token: str


@strawberry.type
class LoginError:
    message: str
    type: str


LoginResult = Annotated[
    LoginSuccess | LoginError,
    strawberry.union("LoginResult"),
]
