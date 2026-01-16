from flask import Flask
from internal.services.smart_gas import (
    get_legal_names,
    get_stations,
    get_stations_groups,
)
from internal.repositories.gas_stations import (
    create_legal_name,
    get_legal_name_by_external_id,
    get_station_by_external_id,
    create_station,
    update_legal_name_by_external_id,
    update_station_by_external_id,
    get_group_by_external_id,
    create_group,
    update_group_by_external_id,
)
from internal.models.gas_stations import GasStation, GasStationGroup, LegalName


def init_smartgas_cli(app: Flask):
    # @app.cli.command("populate-gas-stations")
    def populate_gas_stations():
        smart_gas_api = app.config["SMARTGAS_API_URL"]

        try:
            stations = get_stations(smart_gas_api)
        except Exception as e:
            app.logger.error(f"Unable to bring gas stations from endpoint - {e}")
        else:
            for station in stations:
                legal_name_id = None
                group_id = None
                try:
                    legal_name = get_legal_name_by_external_id(
                        station["external_legal_name_id"]
                    )
                    legal_name_id = legal_name.id
                except:
                    app.logger.error(
                        f"Legal Name {station['external_legal_name_id']} not loaded yet"
                    )

                try:
                    group = get_group_by_external_id(station["external_group_id"])
                    group_id = group.id
                except:
                    app.logger.error(
                        f"Group {station['external_group_id']} is not loaded yet"
                    )

                try:
                    s = get_station_by_external_id(station["external_id"])
                except Exception as e:
                    app.logger.error(f"something wrong happend - {e}")
                    continue

                del station["external_group_id"]
                del station["external_legal_name_id"]

                new_station_data = station | {
                    "group_id": group_id,
                    "legal_name_id": legal_name_id,
                }
                #
                if not s:
                    try:
                        create_station(GasStation(**new_station_data))
                    except Exception as e:
                        app.logger.error(f"unable to create gas station - {e}")
                else:
                    try:
                        update_station_by_external_id(
                            station["external_id"], new_station_data
                        )
                    except Exception as e:
                        app.logger.error(f"unable to update gas station - {e}")

    #
    # @app.cli.command("populate-groups")
    def populate_groups():
        smart_gas_api = app.config["SMARTGAS_API_URL"]

        try:
            groups = get_stations_groups(smart_gas_api)
        except Exception as e:
            app.logger.error(f"Unable to bring data from gas station groups - {e}")
        else:
            for group in groups:
                try:
                    g = get_group_by_external_id(group["external_id"])
                except Exception as e:
                    app.logger.error(f"something wrong happend - {e}")
                    continue

                if not g:
                    try:
                        create_group(GasStationGroup(**group))
                    except Exception as e:
                        app.logger.error(f"Unable to create group - {e}")
                else:
                    try:
                        update_group_by_external_id(group["external_id"], group)
                    except Exception as e:
                        app.logger.error(f"Unable to update group - {e}")

    # @app.cli.command("populate-legal-names")
    def populate_legal_names():
        smart_gas_api = app.config["SMARTGAS_API_URL"]

        try:
            legal_names = get_legal_names(smart_gas_api)
        except Exception as e:
            app.logger.error(f"Unable to bring data from legal names - {e}")
        else:
            for legal_name in legal_names:
                try:
                    l = get_legal_name_by_external_id(legal_name["external_id"])
                except Exception as e:
                    app.logger.error(f"something wrong happend - {e}")
                    continue

                if not l:
                    try:
                        create_legal_name(LegalName(**legal_name))
                    except Exception as e:
                        app.logger.error(f"Unable to create legal name - {e}")
                else:
                    try:
                        update_legal_name_by_external_id(
                            legal_name["external_id"], legal_name
                        )
                    except Exception as e:
                        app.logger.error(f"Unable to update legal name - {e}")

    @app.cli.command("populate-all")
    def populate_all():
        populate_groups()
        populate_legal_names()
        populate_gas_stations()
