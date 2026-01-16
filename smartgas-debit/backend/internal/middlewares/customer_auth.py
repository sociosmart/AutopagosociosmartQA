from flask import current_app, g
from functools import wraps
from internal.exceptions.http import InternalServerError, UnauthorizedError
from internal.services.smart_gas import verify_customer, update_swit
from internal.repositories.customers import (
    get_customer_by_external_id,
    create_customer,
    update_customer_by_id,
)
from internal.models.customers import Customer
from internal.services.swit import create_customer as swit_create_customer

from flask import request


def customer_auth():
    def _customer_auth(f):
        @wraps(f)
        def __customer_auth(*args, **kwargs):
            authorization = request.headers.get("Authorization")

            splitted = authorization.split(" ")
            smart_gas_api = current_app.config["SMARTGAS_API_URL"]

            token = splitted[1]
            try:
                customer_data = verify_customer(smart_gas_api, token)
            except UnauthorizedError:
                raise UnauthorizedError
            except Exception as e:
                current_app.logger.error(
                    f"Unable to bring data from smart gas servers - {e}"
                )
                raise InternalServerError

            try:
                customer = get_customer_by_external_id(customer_data["external_id"])
            except Exception as e:
                current_app.logger.error(f"Unable to get customer data - {e}")
                raise InternalServerError

            if not customer:
                try:
                    c = Customer(**customer_data)
                    customer = create_customer(c)
                except Exception as e:
                    current_app.logger.error(f"Unable to create customer - {e}")
                    raise InternalServerError
            else:
                try:
                    update_customer_by_id(customer.id, customer_data)
                except Exception as e:
                    current_app.logger.error(f"Unable to update customer - {e}")
                    raise InternalServerError

            if not customer.swit_customer_id:
                try:
                    id = swit_create_customer(
                        email=customer.email,
                        first_name=customer.first_name,
                        last_name=customer.last_name,
                    )
                    update_customer_by_id(customer.id, {"swit_customer_id": id})
                except Exception as e:
                    current_app.logger.error(f"Unable to create customer in swit - {e}")
                    raise InternalServerError
                else:
                    try:
                        update_swit(smart_gas_api, token, id)
                    except Exception as e:
                        current_app.logger.error(f"Unable to update swit token - {e}")

            g.customer = customer
            result = f(*args, **kwargs)
            return result

        return __customer_auth

    return _customer_auth
