from flask import Flask

from internal.models.authorization import Authorization
from internal.repositories.authorization import create_authorized_app


def init_authorization_cli(app: Flask):

    @app.cli.command("create-authorized-app")
    def create_authorization():
        app_name = input("app name: ")
        app = create_authorized_app(Authorization(app_name=app_name))

        print(f"app_key={app.app_key}")
        print(f"api_key={app.api_key}")
