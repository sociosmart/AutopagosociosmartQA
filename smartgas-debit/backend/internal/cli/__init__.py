from flask import Flask

from internal.cli.authorization import init_authorization_cli

from .smart_gas import init_smartgas_cli


def init_cli(app: Flask):
    init_smartgas_cli(app)
    init_authorization_cli(app)
