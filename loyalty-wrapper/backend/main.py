from contextlib import asynccontextmanager

from fastapi import FastAPI

from core.config import settings
from core.db import init_db
from core.logging.logger import logger
from gql.app import init_graphql


@asynccontextmanager
async def lifespan(_: FastAPI):
    await init_db()
    # TODO: Implement database connection ping/pong here
    yield


def create_app():
    app = FastAPI(
        docs_url=settings.docs.swagger_url,
        redoc_url=settings.docs.redoc_url,
        debug=settings.debug,
        title=settings.docs.title,
        summary=settings.docs.summary,
        version=settings.docs.version,
        lifespan=lifespan,
    )

    @app.get("/healthcheck")
    async def healthcheck():
        return {"status": "healthy"}

    init_graphql(app)

    return app


app = create_app()

if __name__ == "__main__":
    import uvicorn

    uvicorn.run(
        "main:app",
        host=str(settings.server.host),
        port=settings.server.port,
        reload=settings.server.reload,
    )
