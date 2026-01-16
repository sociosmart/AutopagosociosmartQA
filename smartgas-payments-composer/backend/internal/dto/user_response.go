package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type UserMeResponse struct {
	ID          uuid.UUID    `json:"id"`
	Email       string       `json:"email"`
	FirstName   string       `json:"first_name"`
	LastName    string       `json:"last_name"`
	IsAdmin     bool         `json:"is_admin"`
	CreatedAt   time.Time    `json:"created_at"`
	Permissions []Permission `json:"permissions"`
	Groups      []struct {
		Name        string       `json:"name"`
		Permissions []Permission `json:"permissions"`
	} `json:"groups"`
}

type Permission struct {
	Name string `json:"name"`
}

type UserListResponse struct {
	ID          uuid.UUID    `json:"id"`
	Email       string       `json:"email"`
	FirstName   string       `json:"first_name"`
	LastName    string       `json:"last_name"`
	CreatedAt   time.Time    `json:"created_at"`
	Active      bool         `json:"active"`
	IsAdmin     bool         `json:"is_admin"`
	Permissions []Permission `json:"permissions"`
	Groups      []struct {
		Name        string       `json:"name"`
		Permissions []Permission `json:"permissions"`
	} `json:"groups"`
	GasStations []struct {
		Name string `json:"name"`
	} `json:"gas_stations"`
}

type UserDetailResponse struct {
	ID          uuid.UUID `json:"id"`
	Email       string    `json:"email"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Active      bool      `json:"active"`
	IsAdmin     bool      `json:"is_admin"`
	Permissions []struct {
		ID   uuid.UUID `json:"id"`
		Name string    `json:"name"`
	} `json:"permissions"`
	Groups []struct {
		ID   uuid.UUID `json:"id"`
		Name string    `json:"name"`
	} `json:"groups"`
	GasStations []struct {
		ID   uuid.UUID `json:"id"`
		Name string    `json:"name"`
	} `json:"gas_stations"`
}

func (u *UserDetailResponse) MarshalJSON() ([]byte, error) {
	if u.Permissions == nil {
		u.Permissions = make([]struct {
			ID   uuid.UUID "json:\"id\""
			Name string    "json:\"name\""
		}, 0)
	}

	if u.Groups == nil {
		u.Groups = make([]struct {
			ID   uuid.UUID "json:\"id\""
			Name string    "json:\"name\""
		}, 0)

		if u.GasStations == nil {
			u.GasStations = make([]struct {
				ID   uuid.UUID "json:\"id\""
				Name string    "json:\"name\""
			}, 0)
		}
	}

	return json.Marshal(*u)
}

func (u *UserMeResponse) MarshalJSON() ([]byte, error) {
	if u.Permissions == nil {
		u.Permissions = make([]Permission, 0)
	}

	if u.Groups == nil {
		u.Groups = make([]struct {
			Name        string       "json:\"name\""
			Permissions []Permission "json:\"permissions\""
		}, 0)
	}

	return json.Marshal(*u)
}

func (u *UserListResponse) MarshalJSON() ([]byte, error) {
	if u.Permissions == nil {
		u.Permissions = make([]Permission, 0)
	}

	if u.Groups == nil {
		u.Groups = make([]struct {
			Name        string       "json:\"name\""
			Permissions []Permission "json:\"permissions\""
		}, 0)
	}
	if u.GasStations == nil {
		u.GasStations = make([]struct {
			Name string "json:\"name\""
		}, 0)
	}

	return json.Marshal(*u)
}
