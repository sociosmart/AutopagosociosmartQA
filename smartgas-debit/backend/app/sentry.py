import sentry_sdk
from .config import SENTRY_DSN, SENTRY_ENABLED, ENVIRONMENT, DEBUG


def init_sentry():
    if SENTRY_ENABLED:
        sentry_sdk.init(
            dsn=SENTRY_DSN,
            debug=not not DEBUG,
            environment=ENVIRONMENT,
            # Set traces_sample_rate to 1.0 to capture 100%
            # of transactions for tracing.
            traces_sample_rate=1.0,
            _experiments={
                # Set continuous_profiling_auto_start to True
                # to automatically start the profiler on when
                # possible.
                "continuous_profiling_auto_start": True,
            },
        )
    else:
        sentry_sdk.init(dsn="")
