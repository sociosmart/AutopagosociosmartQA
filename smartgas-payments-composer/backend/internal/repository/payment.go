package repository

import (
	"smartgas-payment/internal/models"
	"smartgas-payment/internal/schemas"
	"smartgas-payment/internal/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StatsForCustomerOpts struct {
	LowDate  time.Time
	HighDate time.Time
}

type CustomerStats struct {
	TotalReported     float64
	TotalTransactions int
}

//go:generate mockery --name PaymentRepository --filename=mock_payment.go --inpackage=true
type PaymentRepository interface {
	CreatePaymentIntent(*models.Payment) error
	GetPaymentByStripePaymentIntentID(string) (*models.Payment, error)
	UpdateByID(uuid.UUID, *models.Payment) (bool, error)
	List(*schemas.Pagination, any) ([]*models.Payment, error)
	GetByID(uuid.UUID) (*models.Payment, error)
	CreateEvent(*models.PaymentEvent) error
	GetLastEventByPaymentID(uuid.UUID) (*models.PaymentEvent, error)
	GetByIDForCustomer(uuid.UUID, uuid.UUID) (*models.Payment, error)
	GetByIDPreloaded(uuid.UUID) (*models.Payment, error)
	GetStatsForCustomer(uuid.UUID, StatsForCustomerOpts) (*CustomerStats, error)
}

type paymentRepository struct {
	db *gorm.DB
}

func ProvidePaymentRepository(db *gorm.DB) *paymentRepository {
	return &paymentRepository{
		db: db,
	}
}

func (pr *paymentRepository) List(
	pagination *schemas.Pagination,
	filters any,
) ([]*models.Payment, error) {
	var payments []*models.Payment

	relatedTables := []string{
		"LEFT JOIN customers as Customer ON Customer.id = payments.customer_id",
		"GasPump",
		"GasPump.GasStation",
		"INNER JOIN gas_stations as GasStation ON GasStation.id = GasPump.gas_station_id",
	}
	filterQuery := `payment_provider LIKE @search OR
  payments.id LIKE @search OR
  CONCAT(Customer.first_name, ' ', Customer.first_last_name, ' ', Customer.second_last_name) LIKE @search OR
  CONCAT(GasStation.name, ' ', GasPump.number) LIKE @search 
  `
	if utils.CheckIfStationsExist(filters) {
		filterQuery = `GasStation.id IN @stations AND  (payment_provider LIKE @search OR
  payments.id LIKE @search OR
  CONCAT(Customer.first_name, ' ', Customer.first_last_name, ' ', Customer.second_last_name) LIKE @search OR
  CONCAT(GasStation.name, ' ', GasPump.number) LIKE @search)
  `
	}

	result := pr.db.
		Joins(relatedTables[0]).
		InnerJoins(relatedTables[1]).
		Preload("Customer").
		Preload(relatedTables[2]).
		Preload("Events", func(db *gorm.DB) *gorm.DB {
			return db.Order("payment_events.created_at DESC")
		}).
		Joins(relatedTables[3]).
		Scopes(utils.Paginate(pagination, payments, pr.db, filterQuery, filters, relatedTables...)).
		Order("created_at desc").
		Where(filterQuery, filters).
		Find(&payments)

	if result.Error != nil {
		return nil, result.Error
	}

	return payments, nil
}

func (pr *paymentRepository) CreatePaymentIntent(payment *models.Payment) error {
	if result := pr.db.Create(&payment); result.Error != nil {
		return result.Error
	}

	return nil
}

func (pr *paymentRepository) GetPaymentByStripePaymentIntentID(
	pid string,
) (*models.Payment, error) {
	var payment models.Payment
	if result := pr.db.
		Preload("GasPump.GasStation").
		Preload("Customer").
		First(&payment, "external_transaction_id = ?", pid); result.Error != nil {
		return nil, result.Error
	}

	return &payment, nil
}

func (pr *paymentRepository) UpdateByID(id uuid.UUID, payment *models.Payment) (bool, error) {
	result := pr.db.Model(payment).Omit("Campaign").Where("id = ?", id).Updates(payment)

	if result.Error != nil {
		return false, result.Error
	}

	if result.RowsAffected < 1 {
		return false, nil
	}

	return true, nil
}

func (pr *paymentRepository) GetByIDForCustomer(
	id uuid.UUID,
	customerID uuid.UUID,
) (*models.Payment, error) {
	var payment models.Payment

	result := pr.db.
		Preload("GasPump.GasStation").
		Preload("Events", func(tx *gorm.DB) *gorm.DB {
			return tx.Order("payment_events.created_at desc").Limit(1)
		}).
		Where("id  = ? AND customer_id = ?", id, customerID).
		First(&payment)

	if result.Error != nil {
		return nil, result.Error
	}

	return &payment, nil
}

func (pr *paymentRepository) GetByID(id uuid.UUID) (*models.Payment, error) {
	var payment *models.Payment

	if result := pr.db.First(&payment, id); result.Error != nil {
		return nil, result.Error
	}

	return payment, nil
}

func (pr *paymentRepository) GetByIDPreloaded(id uuid.UUID) (*models.Payment, error) {
	var payment *models.Payment

	if result := pr.db.
		Preload("GasPump.GasStation").
		Preload("Customer").
		Preload("Campaign").
		Preload("Level").
		First(&payment, id); result.Error != nil {
		return nil, result.Error
	}

	return payment, nil
}

func (pr *paymentRepository) CreateEvent(event *models.PaymentEvent) error {
	if result := pr.db.Create(&event); result.Error != nil {
		return result.Error
	}

	return nil
}

func (pr *paymentRepository) GetLastEventByPaymentID(
	paymentID uuid.UUID,
) (*models.PaymentEvent, error) {
	var event models.PaymentEvent
	err := pr.db.Order("created_at desc").Where("payment_id = ?", paymentID).First(&event).Error
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (pr *paymentRepository) GetStatsForCustomer(
	cusId uuid.UUID,
	opts StatsForCustomerOpts,
) (*CustomerStats, error) {
	var data CustomerStats
	result := pr.db.Model(models.Payment{}).
		Select("SUM(real_amount_reported) as total_reported, COUNT(*) as total_transactions").
		Group("customer_id").
		Where("real_amount_reported > 0 AND created_at BETWEEN ? AND ?", opts.LowDate, opts.HighDate).
		Having("customer_id = ?", cusId).
		Limit(1).
		Scan(&data)

	if result.Error != nil {
		return nil, result.Error
	}

	return &data, nil
}
