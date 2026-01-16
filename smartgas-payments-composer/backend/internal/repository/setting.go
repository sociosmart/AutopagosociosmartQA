package repository

import (
	"smartgas-payment/internal/models"

	"gorm.io/gorm"
)

//go:generate mockery --name SettingRepository --filename=mock_setting.go --inpackage=true
type SettingRepository interface {
	GetAll() ([]*models.Setting, error)
	GetByName(string) (*models.Setting, error)
	GetOrCreate(*models.Setting) (bool, error)
	Update(string, string) (bool, error)
}

type settingRepository struct {
	db *gorm.DB
}

func ProvideSettingRepository(db *gorm.DB) *settingRepository {
	return &settingRepository{
		db: db,
	}
}

func (sr *settingRepository) GetAll() ([]*models.Setting, error) {
	var settings []*models.Setting

	err := sr.db.Find(&settings).Error

	return settings, err

}

func (sr *settingRepository) GetByName(name string) (*models.Setting, error) {
	var setting *models.Setting

	if err := sr.db.Where(models.Setting{Name: name}).First(&setting).Error; err != nil {
		return nil, err
	}

	return setting, nil
}

func (sr *settingRepository) GetOrCreate(setting *models.Setting) (bool, error) {

	result := sr.db.Where(models.Setting{Name: setting.Name}).FirstOrCreate(&setting)

	return result.RowsAffected > 0, result.Error
}

func (sr *settingRepository) Update(name, value string) (bool, error) {
	result := sr.db.Model(models.Setting{}).Where(models.Setting{Name: name}).Update("value", value)

	return result.RowsAffected > 0, result.Error
}
