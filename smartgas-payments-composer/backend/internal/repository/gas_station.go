package repository

import (
	"smartgas-payment/internal/models"
	"smartgas-payment/internal/schemas"
	"smartgas-payment/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

//go:generate mockery --name GasStationRepository --filename=mock_gas_station.go --inpackage=true
type GasStationRepository interface {
	List(*schemas.Pagination, any) ([]*models.GasStation, error)
	Create(*models.GasStation) error
	GetByID(uuid.UUID) (*models.GasStation, error)
	UpdateByID(uuid.UUID, *models.GasStation) (bool, error)
	ListAll(filters any) ([]*models.GasStation, error)
	GetByExternalIDOrCreate(string, *models.GasStation) (bool, error)
	GetByExternalID(string) (*models.GasStation, error)
}

type gasStationRepository struct {
	db *gorm.DB
}

func ProvideGasStationRepository(db *gorm.DB) *gasStationRepository {
	return &gasStationRepository{
		db: db,
	}
}

func (gs *gasStationRepository) List(
	pagination *schemas.Pagination,
	filters any,
) ([]*models.GasStation, error) {
	var stations []*models.GasStation

	filterQuery := "id LIKE @search OR name LIKE @search OR ip LIKE @search"
	// Means that the user is not admin
	if utils.CheckIfStationsExist(filters) {
		filterQuery = "id IN @stations AND (id LIKE @search OR name LIKE @search OR ip LIKE @search)"
	}

	result := gs.db.
		Scopes(utils.Paginate(pagination, stations, gs.db, filterQuery, filters, "")).
		Order("created_at desc").
		Where(filterQuery, filters).
		Find(&stations)

	if result.Error != nil {
		return nil, result.Error
	}

	return stations, nil
}

func (gs *gasStationRepository) ListAll(filters any) ([]*models.GasStation, error) {
	var stations []*models.GasStation

	var result *gorm.DB
	if utils.CheckIfStationsExist(filters) {
		filterQuery := "id IN @stations"
		result = gs.db.
			Order("created_at desc").
			Where(filterQuery, filters).
			Find(&stations)
	} else {
		result = gs.db.
			Order("created_at desc").
			Find(&stations)
	}

	if result.Error != nil {
		return nil, result.Error
	}

	return stations, nil
}

func (gs *gasStationRepository) Create(station *models.GasStation) error {
	if result := gs.db.Create(&station); result.Error != nil {
		return result.Error
	}

	return nil
}

func (gs *gasStationRepository) GetByID(id uuid.UUID) (*models.GasStation, error) {
	var station models.GasStation
	if result := gs.db.First(&station, id); result.Error != nil {
		return nil, result.Error
	}
	return &station, nil
}

func (gs *gasStationRepository) GetByExternalID(externalID string) (*models.GasStation, error) {
	var station models.GasStation

	if result := gs.db.Where("external_id = ? AND active = true", externalID).First(&station); result.Error != nil {
		return nil, result.Error
	}

	return &station, nil
}

func (gs *gasStationRepository) UpdateByID(id uuid.UUID, values *models.GasStation) (bool, error) {
	result := gs.db.Model(models.GasStation{}).Where("id = ?", id).Updates(&values)
	if result.Error != nil {
		return false, result.Error
	}

	if result.RowsAffected > 0 {
		return true, nil
	}

	return false, nil
}

func (gs *gasStationRepository) GetByExternalIDOrCreate(
	externalID string,
	gasStation *models.GasStation,
) (bool, error) {
	result := gs.db.Where(models.GasStation{ExternalID: externalID}).
		Attrs(&gasStation).
		FirstOrCreate(&gasStation)
	if result.Error != nil {
		return false, result.Error
	}

	return result.RowsAffected > 0, nil
}
