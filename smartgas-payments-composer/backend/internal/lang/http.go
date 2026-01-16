package lang

const (
	InternalServerError = "Internal Server Error"
	// TODO: Modify the vars name for EmailOrPasswordIncorrect
	UserOrPasswordIncorrect      = "Email or Password Incorrect"
	NotFoundRecord               = "Not found"
	NoAuthorizationHeader        = "No authorization header setted"
	AuthorizationHeaderMalformed = "Authorization token not sent correctly, it must be i.e. Bearer ..."
	NoAdminPermissions           = "You do not have admin permissions to perform this action"
	DuplicatedEntry              = "Duplicated entries for: "
	RecordUpdated                = "Record Updated"
	Healthy                      = "healthy"
	NotAcceptable                = "Not acceptable value on field: "
	InvalidOrExpiredToken        = "The provided token is invalid or is already expired"
	NotFuelInPump                = "This fuel type is not present in the current pump"
	Synchronized                 = "Synchronized"
	RecordCreated                = "Record created"
	PaymentRequired              = "Payment Required"
	ApplicationUnauthorized      = "Application Unauthorized"
	NotEnoughPermissions         = "Not Enough Permissions"
	PaymentAlreadyServed         = "Payment Already Served"
	AmountChargedGratherThanPaid = "Amount charged grather than paid amount"
	GasStationNotFound           = "Gas Station Not Found"
	UnauthorizedEmployee         = "No permissions to perform this action"
)
