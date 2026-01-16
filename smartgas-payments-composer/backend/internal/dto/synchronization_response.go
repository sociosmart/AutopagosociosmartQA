package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type SynchronizationGetLastSyncResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Status    string    `json:"status"`
}

type SynchronizationListResponse struct {
	ID        uuid.UUID `json:"id"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	Status    string    `json:"status"`
	Errors    []struct {
		Text string `json:"text"`
	} `json:"errors"`
}

func (sl *SynchronizationListResponse) MarshalJSON() ([]byte, error) {
	if sl.Errors == nil {
		sl.Errors = make([]struct {
			Text string "json:\"text\""
		}, 0)
	}

	return json.Marshal(*sl)
}

type SynchronizationListDetailResponse struct {
	ID         uuid.UUID `json:"id"`
	ExternalID string    `json:"external_id"`
	Data       string    `json:"data"`
	Action     string    `json:"action"`
	ErrorText  string    `json:"error_text"`
	CreatedAt  time.Time `json:"created_at"`
}
