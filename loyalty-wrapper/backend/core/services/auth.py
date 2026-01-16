import datetime

import bcrypt
import jwt
from beanie import PydanticObjectId
from bson import ObjectId
from pymongo.errors import DuplicateKeyError

from core.config import settings
from core.dtos.auth import CreateUserData, TokenAccess
from core.exceptions.auth import (
    EmailAlreadyTakenException,
    InactiveUserException,
    InvalidJWTokenException,
    InvalidUserCredentials,
    TokenExpiredException,
)
from core.exceptions.common import InternalServerError
from core.models.users import User
from core.repositories.users import user_repository
from core.dtos.users import AddUserBody, UpdateUserBody

class AuthService:
    def __init__(self):
        # TODO: Take repo by param
        self.__repo = user_repository
        self.__salt = bcrypt.gensalt()

    async def authorize_user(self, token: str) -> User | None:
        try:
            payload = jwt.decode(
                token, settings.jwt.jwt_secret_key, algorithms=[settings.jwt.algorithm]
            )
            id: PydanticObjectId = payload.get("sub")
        except jwt.ExpiredSignatureError:
            raise TokenExpiredException
        except jwt.InvalidTokenError:
            raise InvalidJWTokenException
        except:
            raise InternalServerError

        # getting user
        try:
            _id = ObjectId(id)
            user = await self.__repo.get_user_by_criterias(_id=_id)
        except:
            raise InternalServerError

        return user

    async def get_active_user(self, token) -> User | None:
        try:
            user = await self.authorize_user(token)
        except Exception as e:
            raise e

        if not user:
            return None

        if not user.is_active:
            raise InactiveUserException

        return user

    def verify_password(self, plain_password: str, hashed_password: str) -> bool:
        return bcrypt.checkpw(
            plain_password.encode("utf-8"), hashed_password.encode("utf-8")
        )

    def hash_password(self, password) -> str:
        bytes_ = password.encode("utf-8")
        return bcrypt.hashpw(bytes_, self.__salt).decode("utf-8")

    def create_jwt_token(self, data: str, expiration: int, key: str, algorithm: str):
        to_encode = {
            "exp": datetime.datetime.now() + datetime.timedelta(minutes=expiration),
            "sub": data,
        }
        return jwt.encode(to_encode, key, algorithm)

    async def login(self, email, password) -> TokenAccess:
        try:
            user = await self.__repo.get_user_by_email(email)
        except:
            raise InternalServerError
        else:

            if not user:
                raise InvalidUserCredentials(message="Invalid user credentials")

            if not self.verify_password(password, user.password):
                raise InvalidUserCredentials(message="Invalid user credentials")

            if not user.is_active:
                raise InactiveUserException

            return TokenAccess(
                access_token=self.create_jwt_token(
                    str(user.id),
                    settings.jwt.access_token_expire_minutes,
                    settings.jwt.jwt_secret_key,
                    settings.jwt.algorithm,
                ),
                refresh_token=self.create_jwt_token(
                    str(user.id),
                    settings.jwt.refresh_token_expire_minutes,
                    settings.jwt.jwt_refresh_secret_key,
                    settings.jwt.algorithm,
                ),
            )

    async def create_user(self, data: AddUserBody) -> User:
        data.password = self.hash_password(data.password)

        try:
            return await self.__repo.create_user(User(**data.model_dump()))
        except DuplicateKeyError:
            raise EmailAlreadyTakenException
        except:
            raise InternalServerError
        

        
    async def update_user(self, id: str, data: UpdateUserBody) -> User:
        if data.password: 
            data.password = self.hash_password(data.password)
        else:
            data.password = None
            
        try:
            return await self.__repo.update_user(id, data)
        except DuplicateKeyError:
            raise EmailAlreadyTakenException
        except:
            raise InternalServerError



auth_service = AuthService()
