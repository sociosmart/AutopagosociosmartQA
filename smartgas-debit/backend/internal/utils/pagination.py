from flask_sqlalchemy.pagination import Pagination
from marshmallow import Schema

from internal.dtos.common import PaginatedResponse


def make_pagination(pagination: Pagination, item_model: Schema) -> PaginatedResponse:
    data = {
        "total": pagination.total,
        "page": pagination.page,
        "previous": pagination.prev_num,
        "next": pagination.next_num,
        "items": item_model().dump(pagination.items, many=True),
    }

    return PaginatedResponse().load(data)
