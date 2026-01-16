package repository

import (
	"github.com/google/wire"
)

var RepositorySet = wire.NewSet(
	ProvideUserRepository,
	ProvideGasStationRepository,
	ProvideGasPumpRepository,
	ProvideCustomerRepository,
	ProvidePaymentRepository,
	ProvideSynchronizationRepository,
	ProvideSecurityRepository,
	ProvidePermissionRepository,
	ProvideSettingRepository,
	ProvidePromotionRepository,
	ProvideElegibilityRepository,

	wire.Bind(new(UserRepository), new(*userRepository)),
	wire.Bind(new(GasStationRepository), new(*gasStationRepository)),
	wire.Bind(new(GasPumpRepository), new(*gasPumpRepository)),
	wire.Bind(new(CustomerRepository), new(*customerRepository)),
	wire.Bind(new(PaymentRepository), new(*paymentRepository)),
	wire.Bind(new(SynchronizationRepository), new(*synchronizationRepository)),
	wire.Bind(new(SecurityRepository), new(*securityRepository)),
	wire.Bind(new(PermissionRepository), new(*permissionRepository)),
	wire.Bind(new(SettingRepository), new(*settingRepository)),
	wire.Bind(new(CampaignRepository), new(*campaignRepository)),
	wire.Bind(new(ElegibilityRepository), new(*elegibilityRepository)),
)
