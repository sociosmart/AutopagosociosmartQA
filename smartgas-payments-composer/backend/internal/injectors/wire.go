//go:build exclude

package injectors

import (
	"smartgas-payment/api/v1/controllers"
	"smartgas-payment/api/v1/routes"
	"smartgas-payment/config"
	"smartgas-payment/internal/database"
	"smartgas-payment/internal/middlewares"
	"smartgas-payment/internal/repository"
	"smartgas-payment/internal/server/app"
	"smartgas-payment/internal/services"
	"smartgas-payment/internal/tasks"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"gorm.io/gorm"
)

// Mocks provider
func ProvideUserRepositoryMock() *repository.MockUserRepository {
	return &repository.MockUserRepository{}
}

func ProvidePermissionRepositoryMock() *repository.MockPermissionRepository {
	return &repository.MockPermissionRepository{}
}

func ProvideGasStationRepositoryMock() *repository.MockGasStationRepository {
	return &repository.MockGasStationRepository{}
}

func ProvideGasPumpRepositoryMock() *repository.MockGasPumpRepository {
	return &repository.MockGasPumpRepository{}
}

func ProvideCustomerRepositoryMock() *repository.MockCustomerRepository {
	return &repository.MockCustomerRepository{}
}

func ProvidePaymentRepositoryMock() *repository.MockPaymentRepository {
	return &repository.MockPaymentRepository{}
}

func ProvideCustomerServiceMock() *services.MockCustomerService {
	return &services.MockCustomerService{}
}

func ProvideStripeServiceMock() *services.MockStripeService {
	return &services.MockStripeService{}
}

func ProvideSocioSmartServiceMock() *services.MockSocioSmartService {
	return &services.MockSocioSmartService{}
}

func ProvideSynchronizationTaskMock() *tasks.MockSynchronizationTask {
	return &tasks.MockSynchronizationTask{}
}

func ProvideSynchronizationRepository() *repository.MockSynchronizationRepository {
	return &repository.MockSynchronizationRepository{}
}

func ProvideSecurityRepositoryMock() *repository.MockSecurityRepository {
	return &repository.MockSecurityRepository{}
}

func ProvideSwitServiceMock() *services.MockSwitService {
	return &services.MockSwitService{}
}

func ProvideInvoicingServiceMock() *services.MockInvoicingService {
	return &services.MockInvoicingService{}
}

func ProvideMailServiceMock() *services.MockMailService {
	return &services.MockMailService{}
}

func ProvideSettingRepositoryMock() *repository.MockSettingRepository {
	return &repository.MockSettingRepository{}
}

func ProvideCampaignRepositoryMock() *repository.MockCampaignRepository {
	return &repository.MockCampaignRepository{}
}

func ProvideElegibilityRepositoryMock() *repository.MockElegibilityRepository {
	return &repository.MockElegibilityRepository{}
}

func ProvideDebitServiceMock() *services.MockDebitService {
	return &services.MockDebitService{}
}

var MockSet = wire.NewSet(
	ProvideUserRepositoryMock,
	ProvideGasStationRepositoryMock,
	ProvideGasPumpRepositoryMock,
	ProvideCustomerRepositoryMock,
	ProvidePaymentRepositoryMock,
	ProvideCustomerServiceMock,
	ProvideStripeServiceMock,
	ProvideSocioSmartServiceMock,
	ProvideSynchronizationTaskMock,
	ProvideSynchronizationRepository,
	ProvideSecurityRepositoryMock,
	ProvidePermissionRepositoryMock,
	ProvideSwitServiceMock,
	ProvideInvoicingServiceMock,
	ProvideMailServiceMock,
	ProvideSettingRepositoryMock,
	ProvideCampaignRepositoryMock,
	ProvideElegibilityRepositoryMock,
	ProvideDebitServiceMock,

	wire.Bind(new(repository.UserRepository), new(*repository.MockUserRepository)),
	wire.Bind(new(repository.GasStationRepository), new(*repository.MockGasStationRepository)),
	wire.Bind(new(repository.GasPumpRepository), new(*repository.MockGasPumpRepository)),
	wire.Bind(new(repository.CustomerRepository), new(*repository.MockCustomerRepository)),
	wire.Bind(new(repository.PaymentRepository), new(*repository.MockPaymentRepository)),
	wire.Bind(new(services.CustomerService), new(*services.MockCustomerService)),
	wire.Bind(new(services.StripeService), new(*services.MockStripeService)),
	wire.Bind(new(services.SocioSmartService), new(*services.MockSocioSmartService)),
	wire.Bind(new(tasks.SynchronizationTask), new(*tasks.MockSynchronizationTask)),
	wire.Bind(
		new(repository.SynchronizationRepository),
		new(*repository.MockSynchronizationRepository),
	),
	wire.Bind(new(repository.SecurityRepository), new(*repository.MockSecurityRepository)),
	wire.Bind(new(repository.PermissionRepository), new(*repository.MockPermissionRepository)),
	wire.Bind(new(services.SwitService), new(*services.MockSwitService)),
	wire.Bind(new(services.InvoicingService), new(*services.MockInvoicingService)),
	wire.Bind(new(services.MailService), new(*services.MockMailService)),
	wire.Bind(new(repository.SettingRepository), new(*repository.MockSettingRepository)),
	wire.Bind(new(repository.CampaignRepository), new(*repository.MockCampaignRepository)),
	wire.Bind(new(repository.ElegibilityRepository), new(*repository.MockElegibilityRepository)),
	wire.Bind(new(services.DebitService), new(*services.MockDebitService)),
)

type App struct {
	Router *gin.Engine
	DB     *gorm.DB
}

func ProvideApp(router *gin.Engine, db *gorm.DB) *App {
	return &App{
		Router: router,
		DB:     db,
	}
}

type AppWithMock struct {
	Router                        *gin.Engine
	UserRepositoryMock            *repository.MockUserRepository
	GasStationRepositoryMock      *repository.MockGasStationRepository
	GasPumpRepositoryMock         *repository.MockGasPumpRepository
	PaymentRepositoryMock         *repository.MockPaymentRepository
	CustomerRepositoryMock        *repository.MockCustomerRepository
	extCustomerService            *services.MockCustomerService
	stripeServiceMock             *services.MockStripeService
	socioSmartServiceMock         *services.MockSocioSmartService
	synchronizationTaskMock       *tasks.MockSynchronizationTask
	synchronizationRepositoryMock *repository.MockSynchronizationRepository
	securityRepositoryMock        *repository.MockSecurityRepository
	permissionRepositoryMock      *repository.MockPermissionRepository
	switServiceMock               *services.MockSwitService
	invoicingServiceMock          *services.MockInvoicingService
	mailServiceMock               *services.MockMailService
	settingRepositoryMock         *repository.MockSettingRepository
	campaignRepositoryMock        *repository.MockCampaignRepository
	elebilityRepositoryMock       *repository.MockElegibilityRepository
	debitServiceMock              *services.MockDebitService
}

func ProvideAppWithMock(router *gin.Engine,
	userRepository *repository.MockUserRepository,
	gasStationRepository *repository.MockGasStationRepository,
	gasPumpRepository *repository.MockGasPumpRepository,
	customerRepository *repository.MockCustomerRepository,
	paymentRepository *repository.MockPaymentRepository,
	extCustomerService *services.MockCustomerService,
	stripeServiceMock *services.MockStripeService,
	socioSmartServiceMock *services.MockSocioSmartService,
	synchronizationTaskMock *tasks.MockSynchronizationTask,
	synchronizationRepositoryMock *repository.MockSynchronizationRepository,
	securityRepositoryMock *repository.MockSecurityRepository,
	permissionRepositoryMock *repository.MockPermissionRepository,
	switServiceMock *services.MockSwitService,
	invoicingServiceMock *services.MockInvoicingService,
	mailServiceMock *services.MockMailService,
	settingRepositoryMock *repository.MockSettingRepository,
	campaignRepositoryMock *repository.MockCampaignRepository,
	elebilityRepositoryMock *repository.MockElegibilityRepository,
	debitServiceMock *services.MockDebitService,
) *AppWithMock {
	return &AppWithMock{
		Router:                        router,
		UserRepositoryMock:            userRepository,
		GasStationRepositoryMock:      gasStationRepository,
		GasPumpRepositoryMock:         gasPumpRepository,
		PaymentRepositoryMock:         paymentRepository,
		CustomerRepositoryMock:        customerRepository,
		extCustomerService:            extCustomerService,
		stripeServiceMock:             stripeServiceMock,
		socioSmartServiceMock:         socioSmartServiceMock,
		synchronizationTaskMock:       synchronizationTaskMock,
		synchronizationRepositoryMock: synchronizationRepositoryMock,
		securityRepositoryMock:        securityRepositoryMock,
		permissionRepositoryMock:      permissionRepositoryMock,
		switServiceMock:               switServiceMock,
		invoicingServiceMock:          invoicingServiceMock,
		mailServiceMock:               mailServiceMock,
		settingRepositoryMock:         settingRepositoryMock,
		campaignRepositoryMock:        campaignRepositoryMock,
		elebilityRepositoryMock:       elebilityRepositoryMock,
		debitServiceMock:              debitServiceMock,
	}
}

func InitializeServer() (*App, error) {
	wire.Build(
		config.NewConfig,
		database.ConnectDB,
		app.ProvideGinApp,
		services.ServicesSet,
		repository.RepositorySet,
		tasks.TasksSet,
		controllers.ControllersSet,
		routes.RoutesSet,
		middlewares.MiddlewaresSet,
		wire.NewSet(ProvideApp),
	)

	return &App{}, nil
}

func InitializeSynchronizationTask() (tasks.SynchronizationTask, error) {
	wire.Build(
		config.NewConfig,
		database.ConnectDB,
		services.ServicesSet,
		repository.RepositorySet,
		tasks.TasksSet,
	)

	return nil, nil
}

func InitializeDB() (*gorm.DB, error) {
	wire.Build(
		config.NewConfig,
		database.ConnectDB,
	)

	return &gorm.DB{}, nil
}

func InitializeUserRepository() (repository.UserRepository, error) {
	wire.Build(
		config.NewConfig,
		database.ConnectDB,
		repository.RepositorySet,
	)

	return nil, nil
}

func InitializeSecurityRepository() (repository.SecurityRepository, error) {
	wire.Build(
		config.NewConfig,
		database.ConnectDB,
		repository.RepositorySet,
	)

	return nil, nil
}

func InitializeServerWithMocks() (*AppWithMock, error) {
	wire.Build(
		config.NewConfig,
		app.ProvideGinApp,
		MockSet,
		controllers.ControllersSet,
		routes.RoutesSet,
		middlewares.MiddlewaresSet,
		ProvideAppWithMock,
	)

	return &AppWithMock{}, nil
}
