from typing import Optional

from pydantic import BaseModel, Field

from core.enums.margins import MarginType


class ExternalGasStationResponse(BaseModel):
    fk_social_name: str = Field(alias="Fk_RazonSocial")
    company: str = Field(alias="Compania")
    id_group: str = Field(alias="Cve_Grupo")
    external_id: str = Field(alias="Cve_PuntoDeVenta")
    name: str = Field(alias="NombreComercial")
    cre_permission: str = Field(alias="Num_PermisoCRE")
    vpn: Optional[str] = Field(None, alias="Vpn")
    status: str = Field(alias="Estatus")
    latitude: str = Field(alias="Latitud")
    longitude: str = Field(alias="Longitud")
    street: str = Field(alias="Calle")
    external_number: str = Field(alias="Num_Exterior")
    locality: str = Field(alias="Colonia")
    zip_code: str = Field(alias="CP")
    city: str = Field(alias="Ciudad")
    state: str = Field(alias="Estado")
    regular_price: str = Field(alias="Regular")
    premium_price: str = Field(alias="Premier")
    diesel_price: str = Field(alias="Diesel")


class AddGasStationMargin(BaseModel):
    margin_type: MarginType = Field(MarginType.by_margin, description="Margin Type")
    points: float = Field(ge=0, le=100, description="Point per mxn peso")
    margin: float = Field(
        0, ge=0, le=100, description="Margin for product in gas station"
    )
    gas_station: Optional[str] = Field(None, description="Gas station id")
    product: str = Field(description="Product id")


class UpdateGasStationMargin(BaseModel):
    margin_type: Optional[MarginType] = Field(None, description="Margin Type")
    points: Optional[float] = Field(
        None, ge=0, le=100, description="Point per mxn peso"
    )
    margin: Optional[float] = Field(
        None, ge=0, le=100, description="Margin for product in gas station"
    )
    gas_station: Optional[str] = Field("", description="Gas station id")
    product: Optional[str] = Field(None, description="Product id")


class UpdateGasStationMarginInDb(UpdateGasStationMargin):
    pass
