import strawberry
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from strawberry.fastapi import GraphQLRouter

from gql.context import Context
from gql.mutations import Mutation
from gql.query import Query

def init_graphql(app: FastAPI):
    schema = strawberry.Schema(query=Query, mutation=Mutation)
    graphql_app = GraphQLRouter(schema, context_getter=lambda: Context())

    app.add_middleware(
        CORSMiddleware,
        allow_origins=["*"], 
        allow_credentials=True,
        allow_methods=["*"],  
        allow_headers=["*"],  
    )

    app.include_router(graphql_app, prefix="/graphql")
