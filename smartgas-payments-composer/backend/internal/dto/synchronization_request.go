package dto

type SynchronizationGetLastSyncQueryRequest struct {
	Type string `form:"type" binding:"required,oneof=gas_pumps gas_stations customer_levels" validate:"required,one_of=gas_pumps gas_stations customer_levels"`
}

type SynchronizationNowQueryRequest struct {
	Type string `json:"type" binding:"required,oneof=gas_pumps gas_stations customer_levels" validate:"required,one_of=gas_pumps gas_stations customer_levels"`
}

type SynchronizationListQueryRequest struct {
	Type string `form:"type" binding:"omitempty,oneof=gas_pumps gas_stations customer_levels" validate:"omitempty,oneof=gas_pumps gas_stations customer_levels"`
}

type SynchronizationListDetailPathRequest struct {
	ID string `json:"id" uri:"id" binding:"required,uuid4" validate:"required,uuid4" example:"23ae8c18-4d7a-41a3-a148-8ae2d0a75690"`
}
