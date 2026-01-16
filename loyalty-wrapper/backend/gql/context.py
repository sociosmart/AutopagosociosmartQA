from fastapi import HTTPException
from graphql.pyutils import cached_property
from strawberry.fastapi import BaseContext

from core.exceptions.common import BaseErrorException, InternalServerError
from core.logging.logger import logger
from core.models.authorized_apps import AuthorizedApp
from core.models.customers import Customer
from core.models.users import User
from core.services.auth import auth_service
from core.services.authorized_apps import authorized_app_service
from core.services.smart_gas import smart_gas_sevice


class Context(BaseContext):
    @cached_property
    async def user(self) -> User | None:
        if not self.request:
            return None

        request = self.request

        if not "Authorization" in request.headers:
            # raise GraphQLError("Missing Authorization header")
            return None

        authorization_header = request.headers.get("Authorization", "").split(" ")

        if len(authorization_header) != 2:
            # raise GraphQLError("Authorization header malformed, example: Bearer xyz")
            return None

        if authorization_header[0] != "Bearer":
            # raise GraphQLError("First Part of Authorization header must be Bearer")
            return None

        try:
            user = await auth_service.get_active_user(authorization_header[1])
        except InternalServerError as e:
            logger.error("Error while getting active users -", e.message)
            # raise GraphQLError(e.message)
            return None
        except BaseErrorException as e:
            # raise GraphQLError(e.message)
            return None

        self.user_data = user
        return user

    # @cached_property
    @property
    async def customer(self) -> Customer | None:
        if not self.request:
            return None

        request = self.request

        if not "Authorization" in request.headers:
            # raise GraphQLError("Missing Authorization header")
            return None

        authorization_header = request.headers.get("Authorization", "").split(" ")

        if len(authorization_header) != 2:
            # raise GraphQLError("Authorization header malformed, example: Bearer xyz")
            return None

        if authorization_header[0] != "Token":
            # raise GraphQLError("First Part of Authorization header must be Bearer")
            return None

        try:
            customer = await smart_gas_sevice.update_or_create_customer(
                key="Token", value=authorization_header[1]
            )
        except InternalServerError as e:
            logger.error("Error while getting active users -", e.message)
            # raise GraphQLError(e.message)
            return None
        except BaseErrorException as e:
            # raise GraphQLError(e.message)
            return None

        self.customer_data = customer

        return customer

    # @cached_property
    @property
    async def authorized_app(self) -> AuthorizedApp | None:
        if not self.request:
            return None

        request = self.request

        if not "X-APP-KEY" in request.headers:
            # raise GraphQLError("Missing Authorization header")
            return None

        if not "X-API-KEY" in request.headers:
            # raise GraphQLError("Missing Authorization header")
            return None

        app_key = request.headers.get("X-APP-KEY")
        api_key = request.headers.get("X-API-KEY")

        try:
            authorized_app = (
                await authorized_app_service.get_authorized_app_by_criterias(
                    app_key=app_key,
                    api_key=api_key,
                    is_active=True,
                )
            )
        except InternalServerError as e:
            logger.error("Error while getting authorized app -", e.message)
            # raise GraphQLError(e.message)
            return None
        except BaseErrorException as e:
            logger.error("Exception Error while getting authorized app -", e.message)
            # raise GraphQLError(e.message)
            return None

        self.authorized_app_data = authorized_app

        return authorized_app
