package repository

import (
	"smartgas-payment/internal/models"
	"smartgas-payment/internal/schemas"
	"smartgas-payment/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

//go:generate mockery --name GasPumpRepository --filename=mock_gas_pump.go --inpackage=true
type GasPumpRepository interface {
	List(*schemas.Pagination, any) ([]*models.GasPump, error)
	Create(*models.GasPump) error
	GetByID(uuid.UUID) (*models.GasPump, error)
	UpdateByID(uuid.UUID, *models.GasPump) (bool, error)
	GetByExternalIDOrCreate(string, *models.GasPump) (bool, error)
	GetActiveByID(uuid.UUID) (*models.GasPump, error)
	GetByGasStationAndNumber(uuid.UUID, string) (*models.GasPump, error)
}

type gasPumpRepository struct {
	db *gorm.DB
}

func ProvideGasPumpRepository(db *gorm.DB) *gasPumpRepository {
	return &gasPumpRepository{
		db: db,
	}
}

func (gp *gasPumpRepository) List(
	pagination *schemas.Pagination,
	filters any,
) ([]*models.GasPump, error) {
	var pumps []*models.GasPump

	filterQuery := "number LIKE @search OR gas_pumps.id LIKE @search OR GasStation.name LIKE @search OR GasStation.ip LIKE @search"
	if utils.CheckIfStationsExist(filters) {
		filterQuery = " GasStation.id IN @stations AND (number LIKE @search OR gas_pumps.id LIKE @search OR GasStation.name LIKE @search OR GasStation.ip LIKE @search)"
	}

	result := gp.db.
		InnerJoins("GasStation").
		Scopes(utils.Paginate(pagination, pumps, gp.db, filterQuery, filters, "GasStation")).
		Order("created_at desc").
		Where(filterQuery, filters).
		Find(&pumps)

	if result.Error != nil {
		return nil, result.Error
	}

	return pumps, nil
}

func (gp *gasPumpRepository) Create(pump *models.GasPump) error {
	if result := gp.db.Create(&pump); result.Error != nil {
		return result.Error
	}

	return nil
}

func (gp *gasPumpRepository) GetByID(id uuid.UUID) (*models.GasPump, error) {
	var pump models.GasPump
	if result := gp.db.InnerJoins("GasStation").First(&pump, id); result.Error != nil {
		return nil, result.Error
	}
	return &pump, nil
}

func (gp *gasPumpRepository) GetByGasStationAndNumber(
	gasStationID uuid.UUID,
	number string,
) (*models.GasPump, error) {
	var pump models.GasPump
	if result := gp.db.InnerJoins("GasStation").Where("gas_pumps.gas_station_id = ? AND gas_pumps.number = ? AND gas_pumps.active = true", gasStationID, number).First(&pump); result.Error != nil {
		return nil, result.Error
	}

	return &pump, nil
}

func (gp *gasPumpRepository) UpdateByID(id uuid.UUID, values *models.GasPump) (bool, error) {
	result := gp.db.Model(models.GasPump{}).Where("id = ?", id).Updates(&values)
	if result.Error != nil {
		return false, result.Error
	}

	if result.RowsAffected > 0 {
		return true, nil
	}

	return false, nil
}

func (gp *gasPumpRepository) GetByExternalIDOrCreate(
	externalID string,
	gasPump *models.GasPump,
) (bool, error) {
	result := gp.db.Where(models.GasPump{ExternalID: externalID}).
		Attrs(&gasPump).
		FirstOrCreate(&gasPump)
	if result.Error != nil {
		return false, result.Error
	}

	return result.RowsAffected > 0, nil
}

func (gp *gasPumpRepository) GetActiveByID(id uuid.UUID) (*models.GasPump, error) {
	var pump models.GasPump
	if result := gp.db.InnerJoins("GasStation").Where(&models.GasPump{Active: utils.BoolAddr(true)}).First(&pump, id); result.Error != nil {
		return nil, result.Error
	}
	return &pump, nil
}
