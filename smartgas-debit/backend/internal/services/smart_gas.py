import json
import requests
from typing import List

from internal.schemas.smart_gas import (
    GasStationSchema,
    GasStationGroupSchema,
    LegalNameSchema,
)
from internal.exceptions.http import UnauthorizedError


def get_stations(api_url) -> List[GasStationSchema]:
    r = requests.get(f"{api_url}/rest/operacion")

    if r.status_code != 200:
        raise Exception(f"error, code {r.status_code}")

    stations = []
    for station in r.json():
        s = dict(
            name=station.get("NombreComercial"),
            external_id=station.get("Cve_PuntoDeVenta"),
            external_group_id=station.get("Cve_Grupo"),
            external_legal_name_id=station.get("Fk_RazonSocial"),
            cre_permission=station.get("Num_PermisoCRE"),
            latitude=station.get("Latitud"),
            longitude=station.get("Longitud"),
        )
        schema = GasStationSchema()
        stations.append(schema.load(s))

    return stations


def get_stations_groups(api_url):
    r = requests.get(f"{api_url}/rest/operacion?Grupos")

    if r.status_code != 200:
        raise Exception(f"error, code {r.status_code}")

    groups = []

    for group in r.json():
        g = dict(
            group_name=group.get("NombreComercial"),
            external_id=group.get("Cve_Grupo"),
        )
        groups.append(GasStationGroupSchema().load(g))

    return groups


def get_legal_names(api_url) -> List[LegalNameSchema]:
    r = requests.get(f"{api_url}/rest/operacion?RazonSocial")

    if r.status_code != 200:
        raise Exception(f"error, code {r.status_code}")

    legal_names = []

    for legal_name in r.json():
        l = dict(name=legal_name.get("Nombre"), external_id=legal_name.get("Cve_Razon"))
        legal_names.append(LegalNameSchema().load(l))

    return legal_names


def verify_customer(api_url, token):
    r = requests.post(
        f"{api_url}/rest/clientes?Verifica", data=json.dumps({"Token": token})
    )
    if r.status_code != 200:
        raise Exception(f"Error while verifing customer in smartgas {r.status_code}")

    data = r.json()

    if data[0].get("status") == "error":
        raise UnauthorizedError

    swit_token = data[0].get("TokenSwit")
    d = {
        "external_id": data[0].get("Id"),
        "first_name": data[0].get("Nombre"),
        "last_name": data[0].get("Ap_Paterno") + " " + data[0].get("Ap_Materno"),
        "phone_number": data[0].get("Num_celular"),
        "email": data[0].get("correo"),
        "swit_customer_id": swit_token if swit_token else "",
    }

    return d


def update_swit(api_url, token, token_swit):
    r = requests.post(
        f"{api_url}/rest/auth?Swit",
        data=json.dumps({"Token": token, "TokenSwit": token_swit}),
    )

    if r.status_code != 200:
        raise Exception(f"Error while updating swit token")

    data = r.json()

    if data["Estatus"] == "2":
        raise Exception("Unable to update swit status")
