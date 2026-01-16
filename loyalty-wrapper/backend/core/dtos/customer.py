from typing import Optional

from pydantic import BaseModel, Field


class ExternalCostumerResponse(BaseModel):
    external_id: str = Field(alias="Id")
    name: str = Field(alias="Nombre")
    last_name: str = Field(alias="Ap_Paterno")
    second_last_name: str = Field(alias="Ap_Materno")
    status: str = Field(alias="Estatus")
    token: Optional[str] = Field(alias="Token")
    phone_number: str = Field(alias="Num_celular")
    email: str = Field(alias="correo")
    push_token: Optional[str] = Field(alias="tokenM")
    # token_swit: str = Field(alias="TokenSwit")
