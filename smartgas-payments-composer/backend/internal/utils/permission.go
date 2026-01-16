package utils

import (
	"smartgas-payment/internal/enums"
	"smartgas-payment/internal/models"
)

func GatherUserPermissions(user *models.User) map[string]bool {
	permissions := make(map[string]bool)
	for _, permission := range user.Permissions {
		permissions[permission.Name] = true
	}

	for _, group := range user.Groups {
		for _, permission := range group.Permissions {
			permissions[permission.Name] = true
		}
	}

	return permissions
}

func RequiredPermissionInUser(requiredPermission enums.Permission, user *models.User) (exist bool) {
	permissions := GatherUserPermissions(user)

	_, exist = permissions[string(requiredPermission)]

	return
}

func GatherGasStationIds(user *models.User) (gasStations []string) {
	for _, station := range user.GasStations {
		gasStations = append(gasStations, station.ID.String())
	}

	return
}

func AddStationsFilter(user *models.User, filters map[string]any) {
	if !*user.IsAdmin {
		filters["stations"] = GatherGasStationIds(user)
	}
}

func CheckIfStationsExist(filters any) bool {
	filtersMap := filters.(map[string]any)

	// Means that the user is not admin
	_, ok := filtersMap["stations"]

	return ok
}
