import strawberry

from gql.resolvers.accumulations import AccumulationResolverQuery
from gql.resolvers.benefits import BenefitResolverQuery
from gql.resolvers.customers import CustomersResolverQuery
from gql.resolvers.gas_stations import GasStationResolverQuery
from gql.resolvers.levels import LevelResolverQuery
from gql.resolvers.products import ProductResolverQuery
from gql.resolvers.users import UsersResolverQuery


@strawberry.type
class Query(
    UsersResolverQuery,
    CustomersResolverQuery,
    ProductResolverQuery,
    GasStationResolverQuery,
    AccumulationResolverQuery,
    LevelResolverQuery,
    BenefitResolverQuery,
):
    @strawberry.field
    def health_check(self) -> str:
        return "healthy"
