package middlewares

import "github.com/google/wire"

var MiddlewaresSet = wire.NewSet(
	ProvideAuthMiddleware,
	ProvideCustomerAUthMiddleware,
	ProvideSecurityMiddleware,
)
