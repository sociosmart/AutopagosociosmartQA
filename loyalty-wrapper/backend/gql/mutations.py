import strawberry

from gql.resolvers.accumulations import AccumulationResolverMutation
from gql.resolvers.auth import AuthResolverMutation
from gql.resolvers.benefits import BenefitResolverMutation
from gql.resolvers.gas_stations import GasStationResolverMutation
from gql.resolvers.levels import LevelResolverMutation
from gql.resolvers.products import ProductResolverMutation
from gql.resolvers.users import UserResolverMutation


@strawberry.type
class Mutation(
    AuthResolverMutation,
    ProductResolverMutation,
    GasStationResolverMutation,
    LevelResolverMutation,
    AccumulationResolverMutation,
    BenefitResolverMutation,
    UserResolverMutation,
):
    pass
