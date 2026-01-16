from typing import Optional

import strawberry

from core.dtos.pagination import PaginationParams as Params


# in order to use this common class, we have to add the items attribute.
# there is a known issue regarding with strawberry and generics.
@strawberry.type
class CursorPage:
    next_cursor: Optional[str]
    prev_cursor: Optional[str]
    # items: List[T]
    total: int = 0


@strawberry.experimental.pydantic.input(model=Params, all_fields=True)
class PaginationParams:
    pass
