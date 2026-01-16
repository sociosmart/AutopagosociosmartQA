package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

type Routes []Route

var RoutesSet = wire.NewSet(
	ProvideUserRoutes,
	ProvideAuthRoutes,
	ProvideGasStationRoutes,
	ProvideGasPumpRoutes,
	ProvideV1Routes,
	ProvidePaymenRoutes,
	ProvideCustomerRoutes,
	ProvideSynchronizationRoutes,
	ProvidePermissionRoutes,
	ProvideSettingRoutes,
	ProvideCampaingRoutes,
	ProvideElebilityRoutes,
)

type Route interface {
	Setup(*gin.RouterGroup)
}

func ProvideV1Routes(
	userRoutes *UserRoutes,
	authRoutes *AuthRoutes,
	gasStationRoutes *GasStationRoutes,
	gasPumps *GasPumpRoutes,
	paymentRoutes *PaymentRoutes,
	customerRoutes *CustomerRoutes,
	syncRoutes *SynchronizationRoute,
	permissionRoutes *PermissionRoutes,
	settingRoutes *SettingRoutes,
	promotionRoutes *CampaignRoutes,
	elegibilityRoutes *ElebilityRoutes,
) Routes {
	return Routes{
		userRoutes,
		authRoutes,
		gasStationRoutes,
		gasPumps,
		paymentRoutes,
		customerRoutes,
		syncRoutes,
		permissionRoutes,
		settingRoutes,
		promotionRoutes,
		elegibilityRoutes,
	}
}
