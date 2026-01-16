import json

import requests
from flask import current_app


def get_campaigns():
    config = current_app.config

    return requests.post(
        f"{config.get('SMARTGAS_API_URL')}/rest/Envionotificacion?Campana"
    )


def get_devices_to_notify(external_campaign_id):
    config = current_app.config

    return requests.post(
        f"{config.get('SMARTGAS_API_URL')}/rest/Envionotificacion?Listado",
        data=json.dumps({"Cve_Campana": external_campaign_id}),
    )


def change_status_campaign(external_campaign_id, status="12"):
    config = current_app.config
    return requests.post(
        f"{config.get('SMARTGAS_API_URL')}/rest/Envionotificacion?CampanaEstatus",
        data=json.dumps([{"Cve_Id": external_campaign_id, "Estatus": status}]),
    )
