package repository

import (
	"smartgas-payment/internal/models"

	"gorm.io/gorm"
)

//go:generate mockery --name PermissionRepository --filename=mock_permission.go --inpackage=true
type PermissionRepository interface {
	ListAll() ([]*models.Permission, error)
	ListAllGroups() ([]*models.Group, error)
}

type permissionRepository struct {
	db *gorm.DB
}

func ProvidePermissionRepository(db *gorm.DB) *permissionRepository {
	return &permissionRepository{
		db: db,
	}
}

func (pr *permissionRepository) ListAll() ([]*models.Permission, error) {
	var permissions []*models.Permission

	result := pr.db.Find(&permissions)

	return permissions, result.Error
}

func (pr *permissionRepository) ListAllGroups() ([]*models.Group, error) {
	var groups []*models.Group

	result := pr.db.Find(&groups)

	return groups, result.Error
}
