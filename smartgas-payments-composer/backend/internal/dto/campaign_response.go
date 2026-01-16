package dto

import (
	"encoding/json"
	"time"
)

type CampaignListResponse struct {
	ID          string    `json:"id"`
	Discount    float64   `json:"discount"`
	Name        string    `json:"name"`
	ValidFrom   time.Time `json:"valid_from"`
	ValidTo     time.Time `json:"valid_to"`
	GasStations []struct {
		Name string `json:"name"`
	} `json:"gas_stations"`
	Active bool `json:"active"`
}

func (plr *CampaignListResponse) MarshalJSON() ([]byte, error) {
	if plr.GasStations == nil {
		plr.GasStations = make([]struct {
			Name string `json:"name"`
		}, 0)
	}

	return json.Marshal(*plr)
}

type CampaignDetailResponse struct {
	ID          string    `json:"id"`
	Discount    float64   `json:"discount"`
	Name        string    `json:"name"`
	ValidFrom   time.Time `json:"valid_from"`
	ValidTo     time.Time `json:"valid_to"`
	GasStations []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"gas_stations"`
	Active bool `json:"active"`
}

func (plr *CampaignDetailResponse) MarshalJSON() ([]byte, error) {
	if plr.GasStations == nil {
		plr.GasStations = make([]struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		}, 0)
	}

	return json.Marshal(*plr)
}
