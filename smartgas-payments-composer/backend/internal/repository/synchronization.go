package repository

import (
	"smartgas-payment/internal/models"
	"smartgas-payment/internal/schemas"
	"smartgas-payment/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

//go:generate mockery --name SynchronizationRepository --filename=mock_synchronization.go --inpackage=true
type SynchronizationRepository interface {
	CreateBatchErrors([]*models.SynchronizationError) error
	CreateBatchDetails([]*models.SynchronizationDetail) error
	Create(*models.Synchronization) error
	GetLastByType(string) (*models.Synchronization, error)
	UpdateStatusByID(uuid.UUID, string) (bool, error)
	List(*schemas.Pagination, any) ([]*models.Synchronization, error)
	ListDetails(*schemas.Pagination, any) ([]*models.SynchronizationDetail, error)
}

type synchronizationRepository struct {
	db *gorm.DB
}

func ProvideSynchronizationRepository(db *gorm.DB) *synchronizationRepository {
	return &synchronizationRepository{
		db: db,
	}
}

func (sr *synchronizationRepository) Create(sync *models.Synchronization) error {
	if result := sr.db.Create(&sync); result.Error != nil {
		return result.Error
	}

	return nil
}

func (sr *synchronizationRepository) CreateBatchErrors(errors []*models.SynchronizationError) error {
	if result := sr.db.CreateInBatches(errors, 100); result.Error != nil {
		return result.Error
	}

	return nil
}

func (sr *synchronizationRepository) CreateBatchDetails(details []*models.SynchronizationDetail) error {
	if result := sr.db.CreateInBatches(details, 100); result.Error != nil {
		return result.Error
	}

	return nil
}

func (sr *synchronizationRepository) GetLastByType(t string) (*models.Synchronization, error) {
	var sync models.Synchronization

	if result := sr.db.Order("created_at desc").Where(models.Synchronization{Type: t}).First(&sync); result.Error != nil {
		return nil, result.Error
	}

	return &sync, nil
}

func (sr *synchronizationRepository) UpdateStatusByID(id uuid.UUID, status string) (bool, error) {
	result := sr.db.Model(models.Synchronization{}).Where("id = ?", id).Update("status", status)

	if result.Error != nil {
		return false, result.Error
	}

	return result.RowsAffected > 0, nil
}

func (sr *synchronizationRepository) List(pagination *schemas.Pagination, filters any) ([]*models.Synchronization, error) {
	var syncs []*models.Synchronization

	result := sr.db.
		Preload("Errors").
		Scopes(utils.Paginate(pagination, syncs, sr.db, "", filters, "")).
		Where(filters).
		Order("created_at desc").
		Find(&syncs)

	if result.Error != nil {
		return nil, result.Error
	}

	return syncs, nil
}

func (sr *synchronizationRepository) ListDetails(pagination *schemas.Pagination, filters any) ([]*models.SynchronizationDetail, error) {
	var syncDetails []*models.SynchronizationDetail

	result := sr.db.
		Scopes(utils.Paginate(pagination, syncDetails, sr.db, "", filters, "")).
		Where(filters).
		Order("created_at desc").
		Find(&syncDetails)

	if result.Error != nil {
		return nil, result.Error
	}

	return syncDetails, nil
}
