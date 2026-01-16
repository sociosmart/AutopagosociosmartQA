package repository

import (
	"smartgas-payment/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

//go:generate mockery --name SecurityRepository --filename=mock_security.go --inpackage=true
type SecurityRepository interface {
	GetByKeys(uuid.UUID, uuid.UUID) (*models.AuthorizedApplication, error)
	Create(*models.AuthorizedApplication) error
}

type securityRepository struct {
	db *gorm.DB
}

func ProvideSecurityRepository(db *gorm.DB) *securityRepository {
	return &securityRepository{
		db: db,
	}
}

func (sr *securityRepository) GetByKeys(appKey uuid.UUID, apiKey uuid.UUID) (*models.AuthorizedApplication, error) {
	var authorizedApp *models.AuthorizedApplication

	if result := sr.db.Where(models.AuthorizedApplication{Active: true, AppKey: appKey, ApiKey: apiKey}).First(&authorizedApp); result.Error != nil {
		return nil, result.Error
	}

	return authorizedApp, nil
}

func (sr *securityRepository) Create(authorizedApp *models.AuthorizedApplication) error {
	if result := sr.db.Create(&authorizedApp); result.Error != nil {
		return result.Error
	}

	return nil
}
