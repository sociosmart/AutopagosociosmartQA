package dto

type SecurityHeadersRequest struct {
	// TODO: implement custom tag in order to provide it when displaying 400 error
	AppKey string `json:"APP-KEY" header:"APP-KEY" binding:"required,uuid4" validate:"required,uuid4"`
	ApiKey string `json:"API-KEY" header:"API-KEY" binding:"required,uuid4" validate:"required,uuid4"`
}

type EmployeeValidationHeadersRequest struct {
	// TODO: implement custom tag in order to provide it when displaying 400 error
	ExternalGasStationID string `json:"X-GAS-STATION-ID" header:"X-GAS-STATION-ID" binding:"required" validate:"required"`
	EmployeeID           string `json:"X-EMPLOYEE-ID"    header:"X-EMPLOYEE-ID"    binding:"required" validate:"required"`
	EmployeeNIP          string `json:"X-EMPLOYEE-NIP"   header:"X-EMPLOYEE-NIP"   binding:"required" validate:"required"`
}
