package repository

import (
	"smartgas-payment/internal/models"
	"smartgas-payment/internal/schemas"
	"smartgas-payment/internal/utils"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

//go:generate mockery --name CampaignRepository --filename=mock_campaign.go --inpackage=true
type CampaignRepository interface {
	List(*schemas.Pagination, any) ([]*models.Campaign, error)
	Create(*models.Campaign) error
	UpdateByID(uuid.UUID, *models.Campaign) error
	GetCampaignByID(uuid.UUID, map[string]any) (*models.Campaign, error)
	GetApplicableCampaign(time.Time, uuid.UUID) (*models.Campaign, error)
}

type campaignRepository struct {
	db *gorm.DB
}

func ProvidePromotionRepository(db *gorm.DB) *campaignRepository {
	return &campaignRepository{
		db: db,
	}
}

func (pc *campaignRepository) List(
	pagination *schemas.Pagination,
	filters any,
) ([]*models.Campaign, error) {
	var promotions []*models.Campaign

	filterQuery := "name LIKE @search"
	if utils.CheckIfStationsExist(filters) {
		filterQuery = `
    EXISTS (
      SELECT gas_station_id, campaign_id FROM gas_stations_campaigns as gs
        WHERE campaigns.id = gs.campaign_id AND gs.gas_station_id IN @stations
    )
    AND name LIKE @search
`
	}

	result := pc.db.
		Preload("GasStations").
		Scopes(utils.Paginate(pagination, promotions, pc.db, filterQuery, filters, "")).
		Order("created_at desc").
		Where(filterQuery, filters).
		Find(&promotions)

	if result.Error != nil {
		return nil, result.Error
	}

	return promotions, nil
}

func (cr *campaignRepository) Create(campaign *models.Campaign) error {
	return cr.db.Omit("GasStations.*").Create(campaign).Error
}

func (cr *campaignRepository) UpdateByID(id uuid.UUID, campaing *models.Campaign) error {
	campaing.ID = id
	result := cr.db.
		Omit("GasStations").
		Updates(campaing)

	if campaing.GasStations != nil {
		if err := cr.db.Model(campaing).Omit("GasStations.*").Association("GasStations").Replace(campaing.GasStations); err != nil {
			return err
		}
	}

	return result.Error
}

func (cr *campaignRepository) GetCampaignByID(
	id uuid.UUID,
	filters map[string]any,
) (*models.Campaign, error) {
	filters["id"] = id

	var campaign models.Campaign

	if err := cr.db.Preload("GasStations").Where(filters).First(&campaign).Error; err != nil {
		return nil, err
	}

	return &campaign, nil
}

func (cr *campaignRepository) GetApplicableCampaign(
	date time.Time,
	stationID uuid.UUID,
) (*models.Campaign, error) {
	var campaign models.Campaign

	if err := cr.
		db.
		Model(models.Campaign{}).
		Where("valid_from <= ? AND valid_to >= ? AND active = TRUE AND EXISTS (SELECT * FROM gas_stations_campaigns as gs WHERE gs.campaign_id = campaigns.id AND gs.gas_station_id = ?)", date, date, stationID).
		Order("discount desc").
		First(&campaign).Error; err != nil {
		return nil, err
	}

	return &campaign, nil
}
