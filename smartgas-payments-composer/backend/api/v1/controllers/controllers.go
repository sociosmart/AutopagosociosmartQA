package controllers

import (
	"github.com/google/wire"
)

var ControllersSet = wire.NewSet(
	ProvideUserController,
	ProvideAuthController,
	ProvideGasStationProvider,
	ProvideGasPumpProvider,
	ProvidePaymentController,
	ProvideCustomerController,
	ProvideSynchronizationController,
	ProvidePermissionController,
	ProvideSettingController,
	ProvideCampaignController,
	ProvideElegibityController,

	wire.Bind(new(UserController), new(*userController)),
	wire.Bind(new(IAUthController), new(*AuthController)),
	wire.Bind(new(GasStationController), new(*gasStationController)),
	wire.Bind(new(GasPumpController), new(*gasPumpController)),
	wire.Bind(new(PaymentController), new(*paymentController)),
	wire.Bind(new(CustomerController), new(*customerController)),
	wire.Bind(new(SynchronizationController), new(*synchronizationController)),
	wire.Bind(new(PermissionController), new(*permissionController)),
	wire.Bind(new(SettingController), new(*settingController)),
	wire.Bind(new(CampaignController), new(*campaignController)),
	wire.Bind(new(ElegibilityController), new(*elegibilityController)),
)
