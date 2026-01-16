package repository

import (
	"smartgas-payment/internal/models"
	"smartgas-payment/internal/schemas"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

//go:generate mockery --name CustomerRepository --filename=mock_customer.go --inpackage=true
type CustomerRepository interface {
	Create(*models.Customer) error
	GetCustomerOrCreate(*schemas.Customer) (*models.Customer, bool, error)
	UpdateByID(uuid.UUID, *models.Customer) (bool, error)
	ListAll() ([]*models.Customer, error)
	GetCustomerByExternalID(string) (*models.Customer, error)
}

type customerRepository struct {
	db *gorm.DB
}

func ProvideCustomerRepository(db *gorm.DB) *customerRepository {
	return &customerRepository{
		db: db,
	}
}

func (cr *customerRepository) Create(customer *models.Customer) error {
	if result := cr.db.Create(&customer); result.Error != nil {
		return result.Error
	}

	return nil
}

func (cr *customerRepository) GetCustomerOrCreate(
	customer *schemas.Customer,
) (*models.Customer, bool, error) {
	var cus models.Customer
	result := cr.db.Where(models.Customer{ExternalID: customer.ExternalID}).
		Attrs(customer).
		FirstOrCreate(&cus)

	if result.Error != nil {
		return nil, false, result.Error
	}

	return &cus, result.RowsAffected > 0, nil
}

func (cr *customerRepository) UpdateByID(id uuid.UUID, customer *models.Customer) (bool, error) {
	result := cr.db.Model(customer).Where("id = ?", id).Updates(customer)

	if result.Error != nil {
		return false, result.Error
	}

	if result.RowsAffected < 1 {
		return false, nil
	}

	return true, nil
}

func (cr *customerRepository) ListAll() ([]*models.Customer, error) {
	var customers []*models.Customer

	err := cr.db.Find(&customers).Error

	return customers, err
}

func (cr *customerRepository) GetCustomerByExternalID(externalID string) (*models.Customer, error) {
	var customer models.Customer

	if result := cr.db.Where("external_id = ?", externalID).First(&customer); result.Error != nil {
		return nil, result.Error
	}

	return &customer, nil
}
