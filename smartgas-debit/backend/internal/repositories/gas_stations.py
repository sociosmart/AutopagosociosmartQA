from sqlalchemy import select, update
from internal.db import db

from internal.models.gas_stations import GasStation, GasStationGroup, LegalName


def get_station_by_external_id(external_id):
    stmt = select(GasStation).where(GasStation.external_id == external_id)
    return db.session.execute(stmt).scalar_one_or_none()


def update_station_by_external_id(external_id, data):
    stmt = (
        update(GasStation).where(GasStation.external_id == external_id).values(**data)
    )

    db.session.execute(stmt)
    db.session.commit()


def create_station(station: GasStation) -> GasStation:
    db.session.add(station)
    db.session.commit()
    db.session.refresh(station)

    return station


def get_group_by_external_id(external_id):
    stmt = select(GasStationGroup).where(GasStationGroup.external_id == external_id)

    return db.session.execute(stmt).scalar_one_or_none()


def create_group(group: GasStationGroup) -> GasStationGroup:
    db.session.add(group)
    db.session.commit()
    db.session.refresh(group)

    return group


def update_group_by_external_id(external_id, data):
    stmt = (
        update(GasStationGroup)
        .where(GasStationGroup.external_id == external_id)
        .values(**data)
    )

    db.session.execute(stmt)
    db.session.commit()


def get_gas_stations_paginated(filters={}):
    stmt = select(GasStation).where(GasStation.legal_name_id != None)

    return db.paginate(stmt, error_out=False)


def get_groups_paginated(filters={}):
    stmt = select(GasStationGroup).filter_by(**filters)

    return db.paginate(stmt, error_out=False)


def get_legal_name_by_external_id(external_id) -> LegalName | None:
    stmt = select(LegalName).where(LegalName.external_id == external_id)

    return db.session.execute(stmt).scalar_one_or_none()


def create_legal_name(legal_name: LegalName) -> LegalName:
    db.session.add(legal_name)
    db.session.commit()
    db.session.refresh(legal_name)

    return legal_name


def update_legal_name_by_external_id(external_id, data):
    stmt = update(LegalName).where(LegalName.external_id == external_id).values(**data)

    db.session.execute(stmt)
    db.session.commit()


def get_legal_names_paginated(filters={}):
    stmt = select(LegalName).filter_by(**filters)

    return db.paginate(stmt, error_out=False)


def get_legal_names(filters={}):
    stmt = select(LegalName).filter_by(**filters)

    return db.session.execute(stmt).scalars()


def get_legal_name_by_id(id) -> LegalName | None:
    stmt = select(LegalName).where(LegalName.id == id)
    return db.session.execute(stmt).scalar_one_or_none()


def get_gas_stations(filters={}):
    stmt = select(GasStation).where(GasStation.legal_name_id != None)

    return db.session.execute(stmt).scalars()
