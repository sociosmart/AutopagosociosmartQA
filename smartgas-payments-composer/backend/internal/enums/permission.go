package enums

type Permission string

const (
	ViewUsers Permission = "view_users"
	AddUser              = "add_user"
	EditUser             = "edit_user"

	ViewSynchronizations = "view_synchronizations"
	AddSynchronization   = "add_synchronization"

	ViewGasStations = "view_gas_stations"
	EditGasStation  = "edit_gas_station"
	AddGasStation   = "add_gas_station"

	ViewGasPumps = "view_gas_pumps"
	AddGasPump   = "add_gas_pump"
	EditGasPump  = "edit_gas_pump"

	ViewPayments = "view_payments"

	CanDoActionsPayments = "can_do_payment_actions"

	ViewCampaigns = "view_campaigns"
	AddCampaign   = "add_campaign"
	EditCampaign  = "edit_campaign"

	ViewElegibilityLevels = "view_elegibility_levels"
	AddElegibilityLevel   = "add_elegibility_level"
	EditElegibilityLevel  = "edit_elegibility_level"

	ViewCustomerLevels = "view_customer_levels"
	AddCustomerLevel   = "add_customer_level"
	EditCustomerLevel  = "edit_customer_level"

	ViewAllCustomers         = "view_all_customers"
	ViewAllElegebilityLevels = "view_all_elegibility_levels"
)
