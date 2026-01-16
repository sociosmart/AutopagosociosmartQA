import strawberry

from core.exceptions.auth import InactiveUserException, InvalidUserCredentials
from core.logging.logger import logger
from core.services.auth import auth_service
from gql.types.auth import LoginError, LoginResult, LoginSuccess


@strawberry.type
class AuthResolverMutation:
    @strawberry.field
    async def login(self, email: str, password: str) -> LoginResult:
        try:
            data = await auth_service.login(email, password)
            return LoginSuccess(
                access_token=data.access_token, refresh_token=data.refresh_token
            )
        except InvalidUserCredentials as e:
            return LoginError(message=e.message, type="invalid_credentials")
        except InactiveUserException as e:
            return LoginError(message=e.message, type="inactive_user")
        except Exception as e:
            logger.error(f"Something went wrong while logging in user - {e}")
            return LoginError(
                message="Internal Server Error", type="internal_server_error"
            )
