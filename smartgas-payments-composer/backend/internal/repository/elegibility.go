package repository

import (
	"smartgas-payment/internal/models"
	"smartgas-payment/internal/schemas"
	"smartgas-payment/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

//go:generate mockery --name ElegibilityRepository --filename=mock_elegibility.go --inpackage=true
type ElegibilityRepository interface {
	LevelList(*schemas.Pagination, any) ([]*models.Level, error)
	CreateLevel(*models.Level) error
	UpdateLevelByID(uuid.UUID, *models.Level) (bool, error)
	CustomerLevelList(*schemas.Pagination, any) ([]*models.CustomerLevel, error)
	LevelListAll() ([]*models.Level, error)
	UpdateCustomerLevelByID(uuid.UUID, *models.CustomerLevel) (bool, error)
	CreateCustomerLevel(*models.CustomerLevel) error
	LevelListAllActive() ([]*models.Level, error)
	GetCustomerLevelByCriterias(any) (*models.CustomerLevel, error)
	GetLevelByCriterias(any) (*models.Level, error)
}

type elegibilityRepository struct {
	db *gorm.DB
}

func ProvideElegibilityRepository(
	db *gorm.DB,
) *elegibilityRepository {
	return &elegibilityRepository{
		db: db,
	}
}

func (er *elegibilityRepository) LevelList(
	pagination *schemas.Pagination,
	filters any,
) ([]*models.Level, error) {
	var levels []*models.Level

	sql := `
    name LIKE @search
  `

	result := er.db.
		Scopes(utils.Paginate(pagination, levels, er.db, sql, filters, "")).
		Where(sql, filters).
		Order("created_at desc").
		Find(&levels)

	if result.Error != nil {
		return nil, result.Error
	}

	return levels, nil
}

func (er *elegibilityRepository) CreateLevel(level *models.Level) error {
	return er.db.Create(level).Error
}

func (er *elegibilityRepository) UpdateLevelByID(id uuid.UUID, values *models.Level) (bool, error) {
	result := er.db.Model(models.Level{}).Where("id = ?", id).Updates(&values)

	if result.Error != nil {
		return false, result.Error
	}

	return result.RowsAffected > 0, nil
}

func (er *elegibilityRepository) CustomerLevelList(
	pagination *schemas.Pagination,
	filters any,
) ([]*models.CustomerLevel, error) {
	var customerLevel []*models.CustomerLevel

	sql := `
    CONCAT(Customer.first_name, ' ', Customer.first_last_name, ' ', Customer.second_last_name) LIKE @search
  `

	result := er.db.
		InnerJoins("Customer").
		InnerJoins("Level").
		Scopes(utils.Paginate(pagination, customerLevel, er.db, sql, filters, "Customer", "Level")).
		Where(sql, filters).
		Order("created_at desc").
		Find(&customerLevel)

	if result.Error != nil {
		return nil, result.Error
	}

	return customerLevel, nil
}

func (er *elegibilityRepository) LevelListAll() ([]*models.Level, error) {
	var levels []*models.Level

	err := er.db.Find(&levels).Error

	return levels, err
}

func (er *elegibilityRepository) LevelListAllActive() ([]*models.Level, error) {
	var levels []*models.Level

	err := er.db.Model(models.Level{}).
		Where("active = TRUE").
		Order("min_charges ASC, min_amount ASC").
		Find(&levels).Error

	return levels, err
}

func (er *elegibilityRepository) UpdateCustomerLevelByID(
	id uuid.UUID,
	customerLevel *models.CustomerLevel,
) (bool, error) {
	result := er.db.Model(models.CustomerLevel{}).
		Omit("Level", "Customer").
		Where("id = ?", id).
		Updates(&customerLevel)

	if result.Error != nil {
		return false, result.Error
	}

	return result.RowsAffected > 0, nil
}

func (er *elegibilityRepository) CreateCustomerLevel(cusLevel *models.CustomerLevel) error {
	return er.db.Omit("Level.*", "Customer.*").Create(&cusLevel).Error
}

func (er *elegibilityRepository) GetCustomerLevelByCriterias(
	filters any,
) (*models.CustomerLevel, error) {
	var cusLevel models.CustomerLevel

	result := er.db.
		InnerJoins("Level").
		Where(filters).First(&cusLevel)

	if result.Error != nil {
		return nil, result.Error
	}

	return &cusLevel, nil
}

func (er *elegibilityRepository) GetLevelByCriterias(filters any) (*models.Level, error) {
	var level models.Level

	if err := er.db.Model(models.Level{}).Where(filters).First(&level).Error; err != nil {
		return nil, err
	}

	return &level, nil
}
