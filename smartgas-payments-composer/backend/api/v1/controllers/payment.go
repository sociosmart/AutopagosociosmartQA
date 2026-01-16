package controllers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"smartgas-payment/config"
	"smartgas-payment/internal/dto"
	"smartgas-payment/internal/lang"
	"smartgas-payment/internal/models"
	"smartgas-payment/internal/repository"
	"smartgas-payment/internal/schemas"
	"smartgas-payment/internal/services"
	"smartgas-payment/internal/utils"
	"strings"
	"time"

	internalWebsocket "smartgas-payment/internal/websocket"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/copier"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/charge"
	"github.com/stripe/stripe-go/v72/webhook"
	"gorm.io/gorm"
)

const (
	minChargeAmount float32 = 10.0
)

var paymentWebsocket = internalWebsocket.InitChannels()

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// TODO: Add allowed origins
		return true
	},
}

type PaymentController interface {
	List(*gin.Context)
	CreateIntent(*gin.Context)
	StripeWebhook(*gin.Context)
	AddEvent(*gin.Context)
	GetByIDForCustomer(*gin.Context)
	PaymentNotifierWS(*gin.Context)
	GetPaymentProvider(*gin.Context)
	SignInvoice(*gin.Context)
	ResendInvoice(c *gin.Context)
	GetInvoicePDF(c *gin.Context)
	DoPaymentAction(c *gin.Context)
	CreateIntentOperation(c *gin.Context)
}

type paymentController struct {
	repository        repository.PaymentRepository
	gasPumpRepository repository.GasPumpRepository
	stripeService     services.StripeService
	switService       services.SwitService
	config            config.Config
	socioSmartService services.SocioSmartService
	invoicingService  services.InvoicingService
	mailService       services.MailService
	settingsRepo      repository.SettingRepository
	campaignRepo      repository.CampaignRepository
	elegibilityRepo   repository.ElegibilityRepository
	debitService      services.DebitService
	customerRepo      repository.CustomerRepository
}

func ProvidePaymentController(repository repository.PaymentRepository,
	gasPumpRepository repository.GasPumpRepository,
	stripeService services.StripeService,
	config config.Config,
	socioSmartService services.SocioSmartService,
	switService services.SwitService,
	invoicingService services.InvoicingService,
	mailService services.MailService,
	settingsRepo repository.SettingRepository,
	campaignRepo repository.CampaignRepository,
	elegibilityRepo repository.ElegibilityRepository,
	debitService services.DebitService,
	customerRepo repository.CustomerRepository,
) *paymentController {
	return &paymentController{
		repository:        repository,
		gasPumpRepository: gasPumpRepository,
		stripeService:     stripeService,
		config:            config,
		socioSmartService: socioSmartService,
		switService:       switService,
		invoicingService:  invoicingService,
		mailService:       mailService,
		settingsRepo:      settingsRepo,
		campaignRepo:      campaignRepo,
		elegibilityRepo:   elegibilityRepo,
		debitService:      debitService,
		customerRepo:      customerRepo,
	}
}

// @Summary Payment List
// @Description Payment List
// @Tags Payments
// @Produce json
// @Router /api/v1/payments [GET]
// @Security Bearer
// @Param param query dto.PaginateRequest true "Pagination"
// @Param search query string false "Lookup in payments"
// @Success 200 {object} dto.PaginationResponse{data=[]dto.PaymentListResponse} "Payments paginated"
// @Failure 400 {array} dto.BadRequestMessage "Bad Request, failed on body, form, query..."
// @Failure 401 {object} dto.GeneralMessage "Unauthorized"
// @Failure 500 {object} dto.GeneralMessage "Internal server error"
func (pc *paymentController) List(c *gin.Context) {
	var pagination dto.PaginateRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, utils.MapValidatorError[dto.PaginateRequest](err))
		return
	}

	user := c.MustGet("user").(*models.User)

	var paginationSchema schemas.Pagination

	copier.Copy(&paginationSchema, &pagination)

	search := c.Query("search")

	filters := map[string]any{"search": "%" + search + "%"}

	utils.AddStationsFilter(user, filters)

	payments, err := pc.repository.List(&paginationSchema, filters)
	if err != nil {
		// Logging error in sentry
		opts := &utils.TrackErrorOpts{
			Admin: user,
			Tags:  map[string]string{"auth_type": "admin"},
		}
		utils.TrackError(c, err, opts)
		c.JSON(http.StatusInternalServerError, dto.GeneralMessage{Detail: lang.InternalServerError})
		return
	}

	paymentsResponse := make([]dto.PaymentListResponse, 0)

	copier.Copy(&paymentsResponse, &payments)

	// Making the response
	var paginationResponse dto.PaginationResponse

	copier.Copy(&paginationResponse, &paginationSchema)

	paginationResponse.Data = paymentsResponse

	c.JSON(http.StatusOK, paginationResponse)
}

// @Summary Create Payment intent from operation app
// @Description Creaate Payment intent from opertion app with employee credentials
// @Tags Payments
// @Produce json
// @Router /api/v1/payments/create-intent-operation [POST]
// @Param body body dto.CreatePaymentIntentOperationRequest true "Create payment intent to charge fuel from operation app"
// @Param X-GAS-STATION-ID header string true "Gas Station ID"
// @Param X-EMPLOYEE-ID header string true "Employee ID"
// @Param X-EMPLOYEE-NIP header string true "Employee NIP"
// @Success 201 {object} dto.PaymentCrateIntentResponse "Payment intent information"
// @Failure 400 {array} dto.BadRequestMessage "Bad Request, failed on body, form, query..."
// @Failure 401 {object} dto.GeneralMessage "Unauthorized"
// @Failure 402 {object} dto.GeneralMessage "Payment Required, Unsufficient funds"
// @Failure 404 {object} dto.GeneralMessage "Whether gas station or pump not found"
// @Failure 406 {object} dto.GeneralMessage "Not fuel type in gas pump"
// @Failure 500 {object} dto.GeneralMessage "Internal server error"
func (pc *paymentController) CreateIntentOperation(c *gin.Context) {
	var body dto.CreatePaymentIntentOperationRequest
	if err := c.ShouldBind(&body); err != nil {
		c.JSON(
			http.StatusBadRequest,
			utils.MapValidatorError[dto.CreatePaymentIntentOperationRequest](err),
		)
		return
	}

	var customer *models.Customer

	if body.ChargeType == "customer" {
		cus, err := pc.customerRepo.GetCustomerByExternalID(body.ExternalCustomerID)
		customer = cus
		if err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, dto.GeneralMessage{Detail: "Customer not found"})
				return

			}
			// Logging error in sentry
			opts := &utils.TrackErrorOpts{
				Tags: map[string]string{"auth_type": "employee_authentication"},
			}
			utils.TrackError(c, err, opts)
			c.JSON(
				http.StatusInternalServerError,
				dto.GeneralMessage{Detail: lang.InternalServerError},
			)
			return
		}
	}

	gasStation := c.MustGet("gas_station").(*models.GasStation)
	employeeID := c.MustGet("employee_id").(string)

	gasPump, err := pc.gasPumpRepository.GetByGasStationAndNumber(gasStation.ID, body.PumpNumber)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.GeneralMessage{Detail: lang.NotFoundRecord})
			return

		}
		// Logging error in sentry
		opts := &utils.TrackErrorOpts{
			Tags: map[string]string{"auth_type": "employee_authentication"},
		}
		utils.TrackError(c, err, opts)
		c.JSON(http.StatusInternalServerError, dto.GeneralMessage{Detail: lang.InternalServerError})
		return
	}

	var price float64

	if body.FuelType == "regular" {
		price = *gasPump.RegularPrice
	} else if body.FuelType == "premium" {
		price = *gasPump.PremiumPrice
	} else if body.FuelType == "diesel" {
		price = *gasPump.DieselPrice
	}

	if price <= 0 {
		c.JSON(http.StatusNotAcceptable, dto.GeneralMessage{Detail: lang.NotFuelInPump})
		return
	}

	status := "pending"
	var opts services.DebitReserveFundsOpts
	if body.ChargeType == "customer" {
		opts = services.DebitReserveFundsOpts{
			Amount: body.Amount,
			// ExternalCustomerID:  customer.ExternalID,
			ExternalCustomerID:  body.ExternalCustomerID,
			ExternalLegalNameID: gasPump.GasStation.LegalNameID,
		}
	} else {
		opts = services.DebitReserveFundsOpts{
			Amount:              body.Amount,
			ExternalLegalNameID: gasPump.GasStation.LegalNameID,
			CardKey:             body.CardKey,
		}
	}
	transID, err := pc.debitService.ReserveFunds(opts)
	if err != nil {
		if errors.Is(err, services.DebitUnsufficientFunds) {
			c.JSON(
				http.StatusPaymentRequired,
				dto.GeneralMessage{Detail: "Unsufficient funds or invalid card data"},
			)
			return
		} else if errors.Is(err, services.DebitErrNotFound) {
			c.JSON(http.StatusNotFound, dto.GeneralMessage{Detail: "Customer not found nor gift card"})
			return
		} else if errors.Is(err, services.DebitGiftCardInUse) {
			c.JSON(http.StatusConflict, dto.GeneralMessage{Detail: "Gift Card in use"})
			return
		}

		opts := &utils.TrackErrorOpts{
			// Customer: customer,
			Tags: map[string]string{"auth_type": "employee_authentication"},
		}
		utils.TrackError(c, err, opts)
		c.JSON(http.StatusInternalServerError, dto.GeneralMessage{Detail: lang.InternalServerError})
		return
	}
	// Logging error in sentry
	status = "paid"
	liters := body.Amount / float32(price)
	// TODO: Save record in DB here
	payment := models.Payment{
		ExternalTransactionID: transID,
		FuelType:              body.FuelType,
		Amount:                float32(body.Amount),
		TotalLiter:            liters,
		Price:                 price,
		ChargeType:            "by_total",
		GasPump:               gasPump,
		Customer:              customer,
		PaymentProvider:       "debit",
		Events: []models.PaymentEvent{
			{Type: "funds_reserved"},
		}, // Creating pending event
		Status:          status,
		SetByEmployeeID: &employeeID,
		FromOperations:  utils.BoolAddr(true),
		GiftCardKey:     utils.StringAddr(body.CardKey),
	}
	err = pc.repository.CreatePaymentIntent(&payment)
	if err != nil {
		// TODO: Log in sentry as well as the stripe cancelation error
		pc.debitService.CancelReservation(transID)
		// Logging error in sentry
		opts := &utils.TrackErrorOpts{
			Tags: map[string]string{"auth_type": "employee_authentication"},
		}
		utils.TrackError(c, err, opts)
		c.JSON(http.StatusInternalServerError, dto.GeneralMessage{Detail: lang.InternalServerError})
		return
	}

	setting, err := pc.settingsRepo.GetByName("gas_pump_status")
	if err != nil {
		var csmErr error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			csmErr = errors.New("Gas Pump Status not setted")
		} else {
			csmErr = err
		}
		// Logging error in sentry
		opts := &utils.TrackErrorOpts{
			Customer: customer,
			Tags:     map[string]string{"auth_type": "customer"},
		}
		utils.TrackError(c, csmErr, opts)

	}

	gasPumpEnabled := false

	if setting != nil && setting.Value == "enabled" {
		gasPumpEnabled = true
	}

	if gasPumpEnabled {
		// PRE-SET gas pump
		opts := services.SetGasPumpOptions{
			Number:    payment.GasPump.Number,
			Ip:        payment.GasPump.GasStation.Ip,
			FuelType:  payment.FuelType,
			Amount:    payment.Amount,
			PaymentID: payment.ID,
			Discount:  0,
		}
		data, err := pc.socioSmartService.SetGasPump(opts)

		status := "Error"

		if data != nil && data.Status == 0 {
			status = "This pump has already a preset"
		}

		if err != nil || data.Status == 0 {
			// Logging error in sentry
			optsTE := &utils.TrackErrorOpts{
				Customer: customer,
				Tags:     map[string]string{"auth_type": "customer"},
				Context: map[string]map[string]any{
					"Response": {"Status": status},
					"Data": {
						"gas_station":     payment.GasPump.GasStation.Name,
						"gas_station_id":  payment.GasPump.GasStationID,
						"gas_pump_number": payment.GasPump.Number,
						"gas_pump_id":     payment.GasPumpID,
					},
				},
			}
			utils.TrackError(c, err, optsTE)
			// -1 means refund entire transaction
			// TODO: Log error in sentry
			pc.debitService.CancelReservation(payment.ExternalTransactionID)

			event := &models.PaymentEvent{
				PaymentID: payment.ID,
				// TODO: add new status
				Type: "internal_cancellation",
			}

			pc.repository.CreateEvent(event)
			c.JSON(
				http.StatusInternalServerError,
				dto.GeneralMessage{Detail: "not able to preset pump"},
			)
			return
		} else {
			event := &models.PaymentEvent{
				PaymentID: payment.ID,
				Type:      "pump_ready",
			}

			err = pc.repository.CreateEvent(event)
			if err != nil {
				// TODO: Log error in sentry
				// Logging error in sentry
				opts := &utils.TrackErrorOpts{
					Customer: customer,
					Tags:     map[string]string{"auth_type": "customer"},
				}
				utils.TrackError(c, err, opts)
			}
		}
	}

	c.JSON(http.StatusCreated, dto.PaymentCrateIntentOperationResponse{
		Amount:     float64(body.Amount),
		TotalLiter: float64(liters),
		ID:         payment.ID,
	})
}

// @Summary Create Payment intent
// @Description Create Payment intent in order to refuel gas
// @Tags Payments
// @Produce json
// @Router /api/v1/payments/create-intent [POST]
// @Param Authorization header string true "Token"
// @Param body body dto.CreatePaymentIntentRequest true "Create payment intent to charge fuel"
// @Success 201 {object} dto.PaymentCrateIntentResponse "Payment intent information"
// @Failure 400 {array} dto.BadRequestMessage "Bad Request, failed on body, form, query..."
// @Failure 401 {object} dto.GeneralMessage "Unauthorized"
// @Failure 402 {object} dto.GeneralMessage "Payment Required, Unsufficient funds"
// @Failure 406 {object} dto.GeneralMessage "Not fuel type in gas pump"
// @Failure 500 {object} dto.GeneralMessage "Internal server error"
func (pc *paymentController) CreateIntent(c *gin.Context) {
	// Logic here!

	var body dto.CreatePaymentIntentRequest
	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, utils.MapValidatorError[dto.CreatePaymentIntentRequest](err))
		return
	}

	gasPumpID, _ := uuid.Parse(body.GasPumpID)

	gasPump, err := pc.gasPumpRepository.GetByID(gasPumpID)

	customer := c.MustGet("customer").(*models.Customer)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.GeneralMessage{Detail: lang.NotFoundRecord})
			return

		}
		// Logging error in sentry
		opts := &utils.TrackErrorOpts{
			Customer: customer,
			Tags:     map[string]string{"auth_type": "customer"},
		}
		utils.TrackError(c, err, opts)
		c.JSON(http.StatusInternalServerError, dto.GeneralMessage{Detail: lang.InternalServerError})
		return

	}

	// APPLY DISCOUNT
	var discount float64
	discountType := "none"

	applicablePromotionType, err := pc.settingsRepo.GetByName("applicable_promotion_type")
	if err != nil {
		var csmErr error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			csmErr = errors.New("Applicable promotion type not found")
		} else {
			csmErr = err
		}
		// Logging error in sentry
		opts := &utils.TrackErrorOpts{
			Customer: customer,
			Tags:     map[string]string{"auth_type": "customer"},
		}
		utils.TrackError(c, csmErr, opts)
	}

	promotionTypeInDB := "none"
	if applicablePromotionType != nil {
		promotionTypeInDB = applicablePromotionType.Value
	}
	discountType = promotionTypeInDB

	// Checking valid Campaign
	var campaign *models.Campaign
	var cusLevel *models.CustomerLevel
	if promotionTypeInDB == "campaign" {
		campaign, err = pc.campaignRepo.GetApplicableCampaign(time.Now(), *gasPump.GasStationID)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				// Logging error in sentry
				opts := &utils.TrackErrorOpts{
					Customer: customer,
					Tags:     map[string]string{"auth_type": "customer"},
				}
				utils.TrackError(c, err, opts)
				c.JSON(
					http.StatusInternalServerError,
					dto.GeneralMessage{Detail: lang.InternalServerError},
				)
				return
			}
		}

		// Check applicable promotion
		if campaign != nil {
			discount = *campaign.Discount
		} else {
			discountType = "none"
		}
	} else if promotionTypeInDB == "elegibility" {
		// TODO: check if applicable time zone should be UTC or America/Mazatlan
		loc, _ := time.LoadLocation("America/Mazatlan")
		n := time.Now().In(loc)
		cusLevel, err = pc.elegibilityRepo.GetCustomerLevelByCriterias(map[string]any{
			"customer_id":    customer.ID,
			"validity_month": n.Month(),
			"validity_year":  n.Year(),
		})
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				// Logging error in sentry
				opts := &utils.TrackErrorOpts{
					Customer: customer,
					Tags:     map[string]string{"auth_type": "customer"},
				}
				utils.TrackError(c, err, opts)
				c.JSON(
					http.StatusInternalServerError,
					dto.GeneralMessage{Detail: lang.InternalServerError},
				)
				return
			}
		}
		if cusLevel != nil {
			discount = *cusLevel.Level.Discount
		} else {
			discountType = "none"
		}

	}

	var price float64

	if body.FuelType == "regular" {
		price = *gasPump.RegularPrice - discount
	} else if body.FuelType == "premium" {
		price = *gasPump.PremiumPrice - discount
	} else {
		price = *gasPump.DieselPrice - discount
	}

	if price <= 0 {
		c.JSON(http.StatusNotAcceptable, dto.GeneralMessage{Detail: lang.NotFuelInPump})
		return
	}

	var amount float64
	var liters float64

	if body.ChargeType == "by_liter" {
		liters = float64(body.TotalLiter)
		amount = float64(body.TotalLiter) * price
	} else {
		liters = float64(body.Amount) / price
		amount = float64(body.Amount)
	}

	if amount < float64(minChargeAmount) {
		amount = float64(minChargeAmount)
	}

	transID := ""
	var pi *stripe.PaymentIntent
	status := "pending"
	if body.PaymentProvider == "stripe" {
		pi, err = pc.stripeService.CreatePaymentIntent(amount, customer.StripeCustomerID)
		if err != nil {
			// Logging error in sentry
			opts := &utils.TrackErrorOpts{
				Customer: customer,
				Tags:     map[string]string{"auth_type": "customer"},
			}
			utils.TrackError(c, err, opts)
			c.JSON(
				http.StatusInternalServerError,
				dto.GeneralMessage{Detail: lang.InternalServerError},
			)
			return
		}
		transID = pi.ID
	} else if body.PaymentProvider == "swit" {
		opts := services.ReserveFundsOpts{
			CustomerID: customer.SwitCustomerID,
			SourceID:   body.SourceID,
			Cvv:        body.Cvv,
			Last4:      body.Last4,
			Amount:     float32(amount),
		}
		transID, err = pc.switService.ReserveFunds(opts)
		if err != nil {
			if err == services.ErrProccesingPayment {
				c.JSON(http.StatusPaymentRequired, dto.GeneralMessage{Detail: "Unsufficient funds or invalid card data"})
				return
			}
			// Logging error in sentry
			opts := &utils.TrackErrorOpts{
				Customer: customer,
				Tags:     map[string]string{"auth_type": "customer"},
			}
			utils.TrackError(c, err, opts)
			c.JSON(http.StatusInternalServerError, dto.GeneralMessage{Detail: lang.InternalServerError})
			return

		}
		status = "paid"
	} else if body.PaymentProvider == "debit" {
		opts := services.DebitReserveFundsOpts{
			Amount:              float32(amount),
			ExternalCustomerID:  customer.ExternalID,
			ExternalLegalNameID: gasPump.GasStation.LegalNameID,
		}
		transID, err = pc.debitService.ReserveFunds(opts)
		if err != nil {
			if errors.Is(err, services.DebitUnsufficientFunds) {
				c.JSON(http.StatusPaymentRequired, dto.GeneralMessage{Detail: "Unsufficient funds or invalid card data"})
				return
			}
			opts := &utils.TrackErrorOpts{
				Customer: customer,
				Tags:     map[string]string{"auth_type": "customer"},
			}
			utils.TrackError(c, err, opts)
			c.JSON(http.StatusInternalServerError, dto.GeneralMessage{Detail: lang.InternalServerError})
			return
		}
		// Logging error in sentry
		status = "paid"
	}

	// TODO: Save record in DB here
	payment := models.Payment{
		ExternalTransactionID: transID,
		FuelType:              body.FuelType,
		Amount:                float32(amount),
		TotalLiter:            float32(liters),
		Price:                 price,
		ChargeType:            body.ChargeType,
		GasPump:               gasPump,
		Customer:              customer,
		PaymentProvider:       body.PaymentProvider,
		// Events:                []models.PaymentEvent{{}}, // Creating pending event
		Status:           status,
		DiscountPerLiter: discount,
		DiscountType:     discountType,
	}

	if body.PaymentProvider == "stripe" {
		payment.Events = []models.PaymentEvent{{}}
	} else if body.PaymentProvider == "swit" || body.PaymentProvider == "debit" {
		payment.Events = []models.PaymentEvent{{Type: "funds_reserved"}}
	}

	// Check Discount type i n order to save it in DB
	if discountType == "campaign" && campaign != nil {
		payment.CampaignID = &campaign.ID
	} else if discountType == "elegibility" && cusLevel != nil {
		// Elegibility
		payment.LevelID = cusLevel.LevelID
	}

	err = pc.repository.CreatePaymentIntent(&payment)
	if err != nil {
		// TODO: Log in sentry as well as the stripe cancelation error
		if body.PaymentProvider == "stripe" {
			pc.stripeService.CancelPaymentIntent(transID)
		} else if body.PaymentProvider == "swit" {
			// TODO: apply cancelation
			pc.switService.CancelFundReservation(transID)
		} else if body.PaymentProvider == "debit" {
			pc.debitService.CancelReservation(transID)
		}
		// Logging error in sentry
		opts := &utils.TrackErrorOpts{
			Customer: customer,
			Tags:     map[string]string{"auth_type": "customer"},
		}
		utils.TrackError(c, err, opts)
		c.JSON(http.StatusInternalServerError, dto.GeneralMessage{Detail: lang.InternalServerError})
		return
	}

	response := dto.PaymentCrateIntentResponse{
		// ClientSecret: pi.ClientSecret,
		Amount:     amount,
		TotalLiter: liters,
		ID:         payment.ID,
	}

	if body.PaymentProvider == "stripe" {
		response.ClientSecret = pi.ClientSecret
	}

	// Setup pump in swit
	setting, err := pc.settingsRepo.GetByName("gas_pump_status")
	if err != nil {
		var csmErr error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			csmErr = errors.New("Gas Pump Status not setted")
		} else {
			csmErr = err
		}
		// Logging error in sentry
		opts := &utils.TrackErrorOpts{
			Customer: customer,
			Tags:     map[string]string{"auth_type": "customer"},
		}
		utils.TrackError(c, csmErr, opts)

	}

	gasPumpEnabled := false

	if setting != nil && setting.Value == "enabled" {
		gasPumpEnabled = true
	}

	if gasPumpEnabled && (body.PaymentProvider == "swit" || body.PaymentProvider == "debit") {
		// PRE-SET gas pump
		opts := services.SetGasPumpOptions{
			Number:    payment.GasPump.Number,
			Ip:        payment.GasPump.GasStation.Ip,
			FuelType:  payment.FuelType,
			Amount:    payment.Amount,
			PaymentID: payment.ID,
			Discount:  discount,
		}
		data, err := pc.socioSmartService.SetGasPump(opts)

		status := "Error"

		if data != nil && data.Status == 0 {
			status = "This pump has already a preset"
		}

		if err != nil || data.Status == 0 {
			// Logging error in sentry
			optsTE := &utils.TrackErrorOpts{
				Customer: customer,
				Tags:     map[string]string{"auth_type": "customer"},
				Context: map[string]map[string]any{
					"Response": {"Status": status},
					"Data": {
						"gas_station":     payment.GasPump.GasStation.Name,
						"gas_station_id":  payment.GasPump.GasStationID,
						"gas_pump_number": payment.GasPump.Number,
						"gas_pump_id":     payment.GasPumpID,
						"customer":        payment.Customer.FirstName + " " + payment.Customer.FirstLastName + " " + payment.Customer.SecondLastName,
						"customer_id":     payment.CustomerID,
						"customer_email":  payment.Customer.Email,
						"discount":        opts.Discount,
					},
				},
			}
			utils.TrackError(c, err, optsTE)
			// -1 means refund entire transaction
			// TODO: Log error in sentry
			if body.PaymentProvider == "swit" {
				pc.switService.CancelFundReservation(payment.ExternalTransactionID)
			} else if body.PaymentProvider == "debit" {
				pc.debitService.CancelReservation(payment.ExternalTransactionID)
			}

			event := &models.PaymentEvent{
				PaymentID: payment.ID,
				// TODO: add new status
				Type: "internal_cancellation",
			}

			pc.repository.CreateEvent(event)
			// Send email

			requestFuelSchemaMail := &schemas.FuelRequest{}
			requestFuelSchemaMail.FillData(&payment)
			requestFuelSchemaMail.RefundedAmount = payment.Amount
			requestFuelSchemaMail.Error = true

			opts := services.SendMailOpts{
				Data:         requestFuelSchemaMail,
				Description:  "Hubo un error al intentar hacer tu carga",
				TemplatePath: "fuel_request.html",
				To:           payment.Customer.Email,
			}

			pc.mailService.SendMail(opts)
		} else {
			event := &models.PaymentEvent{
				PaymentID: payment.ID,
				Type:      "pump_ready",
			}

			err = pc.repository.CreateEvent(event)
			if err != nil {
				// TODO: Log error in sentry
				// Logging error in sentry
				opts := &utils.TrackErrorOpts{
					Customer: customer,
					Tags:     map[string]string{"auth_type": "customer"},
				}
				utils.TrackError(c, err, opts)
			}

			channel := paymentWebsocket.GetChannel(payment.ID.String())
			channel.BroadcastJson(dto.PaymentWebsocketNotification{Status: "pump_ready"})
		}
	}

	c.JSON(http.StatusCreated, response)
}

func (pc *paymentController) StripeWebhook(c *gin.Context) {
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		// Logging error in sentry
		opts := &utils.TrackErrorOpts{
			Tags: map[string]string{"webhook": "stripe"},
		}
		utils.TrackError(c, err, opts)
		c.Status(http.StatusServiceUnavailable)
		return
	}
	signatureHeader := c.GetHeader("Stripe-Signature")

	event, err := webhook.ConstructEvent(payload, signatureHeader, pc.config.StripeWebhookSecret)
	if err != nil {
		// Logging error in sentry
		opts := &utils.TrackErrorOpts{
			Tags: map[string]string{"webhook": "stripe"},
		}
		utils.TrackError(c, err, opts)
		c.Status(http.StatusBadRequest)
		return
	}

	switch event.Type {
	case "payment_intent.payment_failed":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			// Logging error in sentry
			opts := &utils.TrackErrorOpts{
				Tags: map[string]string{"webhook": "stripe"},
			}
			utils.TrackError(c, err, opts)
			// TODO: Log in sentry
			c.Status(http.StatusBadRequest)
			return
		}

		payment, err := pc.repository.GetPaymentByStripePaymentIntentID(paymentIntent.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// TODO: Log in sentry that the payment has not been found
				c.Status(http.StatusNotFound)
				return
			}
			// TODO: Log in sentry that there is an internal server error
			// Logging error in sentry
			opts := &utils.TrackErrorOpts{
				Tags: map[string]string{"webhook": "stripe"},
			}
			utils.TrackError(c, err, opts)
			c.Status(http.StatusInternalServerError)
			return
		}

		payment.Status = "failed"

		updated, err := pc.repository.UpdateByID(payment.ID, payment)
		if err != nil {
			// TODO: Log in sentry that there is an internal server error
			// Logging error in sentry
			opts := &utils.TrackErrorOpts{
				Tags: map[string]string{"webhook": "stripe"},
			}
			utils.TrackError(c, err, opts)
			c.Status(http.StatusInternalServerError)
			return
		}

		if !updated {
			// TODO: Log in sentry that there is not payment with that id
			// c.Status(http.StatusNotFound)
			// return
			// Logging error in sentry
			opts := &utils.TrackErrorOpts{
				Tags: map[string]string{"webhook": "stripe"},
			}
			utils.TrackError(c, errors.New("There is not payment with this id to update"), opts)
		}

		event := &models.PaymentEvent{
			PaymentID: payment.ID,
			Type:      "failed",
		}

		err = pc.repository.CreateEvent(event)
		if err != nil {
			// TODO: LOg in sentryerr
			// Logging error in sentry
			opts := &utils.TrackErrorOpts{
				Tags: map[string]string{"webhook": "stripe"},
			}
			utils.TrackError(c, err, opts)
		}

	case "payment_intent.succeeded":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		// activate pump
		if err != nil {
			// Logging error in sentry
			opts := &utils.TrackErrorOpts{
				Tags: map[string]string{"webhook": "stripe"},
			}
			utils.TrackError(c, err, opts)
			// TODO: Log in sentry
			c.Status(http.StatusBadRequest)
			return
		}

		payment, err := pc.repository.GetPaymentByStripePaymentIntentID(paymentIntent.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// TODO: Log in sentry that the payment has not been found
				c.Status(http.StatusNotFound)
				return
			}
			// Logging error in sentry
			opts := &utils.TrackErrorOpts{
				Tags: map[string]string{"webhook": "stripe"},
			}
			utils.TrackError(c, err, opts)
			// TODO: Log in sentry that there is an internal server error
			c.Status(http.StatusInternalServerError)
			return
		}

		payment.Status = "paid"

		updated, err := pc.repository.UpdateByID(payment.ID, payment)
		if err != nil {
			// TODO: Should refund
			// Logging error in sentry
			opts := &utils.TrackErrorOpts{
				Tags: map[string]string{"webhook": "stripe"},
			}
			utils.TrackError(c, err, opts)
			c.Status(http.StatusInternalServerError)
			return
		}

		if !updated {
			// Logging error in sentry
			opts := &utils.TrackErrorOpts{
				Tags: map[string]string{"webhook": "stripe"},
			}
			utils.TrackError(c, err, opts)
			// TODO: Should refund
			c.Status(http.StatusNotFound)
			return
		}

		// logging event
		event := &models.PaymentEvent{
			PaymentID: payment.ID,
			Type:      "paid",
		}

		err = pc.repository.CreateEvent(event)
		if err != nil {
			// Logging error in sentry
			opts := &utils.TrackErrorOpts{
				Tags: map[string]string{"webhook": "stripe"},
			}
			utils.TrackError(c, err, opts)
			// TODO: LOg in sentry
		}

		channel := paymentWebsocket.GetChannel(payment.ID.String())

		channel.BroadcastJson(dto.PaymentWebsocketNotification{Status: "paid"})

		// Setup
		setting, err := pc.settingsRepo.GetByName("gas_pump_status")
		if err != nil {
			var csmErr error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				csmErr = errors.New("Gas Pump Status not setted")
			} else {
				csmErr = err
			}
			// Logging error in sentry
			opts := &utils.TrackErrorOpts{
				Tags: map[string]string{"webhook": "stripe"},
			}
			utils.TrackError(c, csmErr, opts)

		}

		gasPumpEnabled := false

		if setting != nil && setting.Value == "enabled" {
			gasPumpEnabled = true
		}

		if gasPumpEnabled { // PRE-SET gas pump
			opts := services.SetGasPumpOptions{
				Number:    payment.GasPump.Number,
				Ip:        payment.GasPump.GasStation.Ip,
				FuelType:  payment.FuelType,
				Amount:    payment.Amount,
				PaymentID: payment.ID,
				Discount:  payment.DiscountPerLiter,
			}
			data, err := pc.socioSmartService.SetGasPump(opts)
			status := "Error"

			if data != nil && data.Status == 0 {
				status = "This pump has already a preset"
			}
			// Logging error in sentry
			optsTE := &utils.TrackErrorOpts{
				Tags: map[string]string{"webhook": "stripe"},
				Context: map[string]map[string]any{
					"Response": {"Status": status},
					"Data": {
						"gas_station":     payment.GasPump.GasStation.Name,
						"gas_station_id":  payment.GasPump.GasStationID,
						"gas_pump_number": payment.GasPump.Number,
						"gas_pump_id":     payment.GasPumpID,
						"customer":        payment.Customer.FirstName + " " + payment.Customer.FirstLastName + " " + payment.Customer.SecondLastName,
						"customer_id":     payment.CustomerID,
						"customer_email":  payment.Customer.Email,
						"discount":        opts.Discount,
					},
				},
			}
			utils.TrackError(c, err, optsTE)

			if err != nil || data.Status == 0 {
				// -1 means refund entire transaction
				// TODO: Log error in sentry
				pc.stripeService.MakeARefund(payment.ExternalTransactionID, -2)

				requestFuelSchemaMail := &schemas.FuelRequest{}
				requestFuelSchemaMail.FillData(payment)
				requestFuelSchemaMail.RefundedAmount = payment.Amount
				requestFuelSchemaMail.Error = true

				opts := services.SendMailOpts{
					Data:         requestFuelSchemaMail,
					Description:  "Hubo un error al intentar hacer tu carga",
					TemplatePath: "fuel_request.html",
					To:           payment.Customer.Email,
				}

				pc.mailService.SendMail(opts)
			} else {
				event := &models.PaymentEvent{
					PaymentID: payment.ID,
					Type:      "pump_ready",
				}

				err = pc.repository.CreateEvent(event)
				if err != nil {
					// TODO: Log error in sentry
					opts := &utils.TrackErrorOpts{
						Tags: map[string]string{"webhook": "stripe"},
					}
					utils.TrackError(c, err, opts)
				}
				channel.BroadcastJson(dto.PaymentWebsocketNotification{Status: "pump_ready"})
			}
		}

		// Notify user

	case "charge.refunded":
		var chargeBody stripe.Charge
		err := json.Unmarshal(event.Data.Raw, &chargeBody)
		if err != nil {
			// TODO: Log in sentry
			opts := &utils.TrackErrorOpts{
				Tags: map[string]string{"webhook": "stripe"},
			}
			utils.TrackError(c, err, opts)
			c.Status(http.StatusBadRequest)
			return
		}

		params := stripe.ChargeParams{}
		params.AddExpand("refunds")

		cha, err := charge.Get(chargeBody.ID, &params)
		if err != nil {
			opts := &utils.TrackErrorOpts{
				Tags: map[string]string{"webhook": "stripe"},
			}
			utils.TrackError(c, err, opts)
			c.Status(http.StatusNotFound)
			return
		}

		status := "partial_refund"

		if len(cha.Refunds.Data) > 0 {
			re := cha.Refunds.Data[0]

			if sta, ok := re.Metadata["status"]; ok {
				status = sta
			}
		}

		payment, err := pc.repository.GetPaymentByStripePaymentIntentID(chargeBody.PaymentIntent.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// TODO: Log in sentry that the payment has not been found
				c.Status(http.StatusNotFound)
				return
			}
			// TODO: Log in sentry that there is an internal server error
			opts := &utils.TrackErrorOpts{
				Tags: map[string]string{"webhook": "stripe"},
			}
			utils.TrackError(c, err, opts)
			c.Status(http.StatusInternalServerError)
			return
		}

		event := &models.PaymentEvent{
			PaymentID: payment.ID,
			Type:      status,
		}

		err = pc.repository.CreateEvent(event)
		if err != nil {
			opts := &utils.TrackErrorOpts{
				Tags: map[string]string{"webhook": "stripe"},
			}
			utils.TrackError(c, err, opts)
			c.JSON(
				http.StatusInternalServerError,
				dto.GeneralMessage{Detail: lang.InternalServerError},
			)
			return
		}

		if status == "manual_action" {
			pc.repository.CreateEvent(&models.PaymentEvent{
				Type:      "partial_refund",
				PaymentID: payment.ID,
			})
		}

		amount := float32(chargeBody.AmountRefunded / 100)
		decimals := float32(chargeBody.AmountRefunded%100) / 100

		payment.RefundedAmount = float32(amount + decimals)

		_, err = pc.repository.UpdateByID(payment.ID, payment)
		if err != nil {
			// TODO: Log error in sentry
			opts := &utils.TrackErrorOpts{
				Tags: map[string]string{"webhook": "stripe"},
			}
			utils.TrackError(c, err, opts)
			c.JSON(
				http.StatusInternalServerError,
				dto.GeneralMessage{Detail: lang.InternalServerError},
			)
			return
		}

		// TODO: Notify user that money were refunded

		channel := paymentWebsocket.GetChannel(payment.ID.String())

		channel.BroadcastJson(dto.PaymentWebsocketNotification{Status: status})

		paymentWebsocket.DeleteChannel(payment.ID.String())

	}

	c.Status(http.StatusOK)
}

// @Summary Add event to payment
// @Description Add event to payment
// @Tags Payments
// @Produce json
// @Router /api/v1/payments/{id}/events [POST]
// @Param id path string true "uuid4 id" minLength(36) maxLength(39)
// @Param body body dto.AddEventBodyRequest true "Add event for payment"
// @Param APP-KEY header string true "App Key"
// @Param API-KEY header string true "Api Key"
// @Success 201 {object} dto.GeneralMessage "OK"
// @Failure 400 {array} dto.BadRequestMessage "Bad Request, failed on body, form, query..."
// @Failure 401 {object} dto.GeneralMessage "Unauthorized"
// @Failure 406 {object} dto.GeneralMessage "Already Served"
// @Failure 402 {object} dto.GeneralMessage "Payment Required"
// @Failure 404 {object} dto.GeneralMessage "Payment id not found in db"
// @Failure 500 {object} dto.GeneralMessage "Internal server error"
func (pc *paymentController) AddEvent(c *gin.Context) {
	authorizedApp := c.MustGet("application").(*models.AuthorizedApplication)

	var path dto.AddEventPathRequest
	if err := c.ShouldBindUri(&path); err != nil {
		// Logging error in sentry
		opts := &utils.TrackErrorOpts{
			Application: authorizedApp,
			Tags:        map[string]string{"auth_type": "application"},
			Level:       sentry.LevelInfo,
		}
		utils.TrackError(c, err, opts)
		c.JSON(http.StatusBadRequest, utils.MapValidatorError[dto.AddEventPathRequest](err))
		return
	}

	var body dto.AddEventBodyRequest
	if err := c.ShouldBind(&body); err != nil {
		// Logging error in sentry
		opts := &utils.TrackErrorOpts{
			Application: authorizedApp,
			Tags:        map[string]string{"auth_type": "application"},
			Level:       sentry.LevelInfo,
		}
		utils.TrackError(c, err, opts)
		c.JSON(http.StatusBadRequest, utils.MapValidatorError[dto.AddEventBodyRequest](err))
		return
	}

	newId, _ := strings.CutPrefix(strings.ToLower(path.ID), "ap_")

	id, _ := uuid.Parse(newId)
	payment, err := pc.repository.GetByIDPreloaded(id)
	if err != nil {
		// Logging error in sentry
		opts := &utils.TrackErrorOpts{
			Application: authorizedApp,
			Tags:        map[string]string{"auth_type": "application"},
			Level:       sentry.LevelInfo,
		}
		utils.TrackError(c, err, opts)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.GeneralMessage{Detail: lang.NotFoundRecord})
			return

		}
		c.JSON(http.StatusInternalServerError, dto.GeneralMessage{Detail: lang.InternalServerError})
		return

	}

	if payment.Status != "paid" {
		// Logging error in sentry
		opts := &utils.TrackErrorOpts{
			Application: authorizedApp,
			Tags:        map[string]string{"auth_type": "application"},
			Level:       sentry.LevelInfo,
		}
		utils.TrackError(c, errors.New(lang.PaymentRequired), opts)
		c.JSON(http.StatusPaymentRequired, dto.GeneralMessage{Detail: lang.PaymentRequired})
		return
	}

	e, err := pc.repository.GetLastEventByPaymentID(payment.ID)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			// Logging error in sentry
			opts := &utils.TrackErrorOpts{
				Application: authorizedApp,
				Tags:        map[string]string{"auth_type": "application"},
				Level:       sentry.LevelInfo,
			}
			utils.TrackError(c, err, opts)
			c.JSON(
				http.StatusInternalServerError,
				dto.GeneralMessage{Detail: lang.InternalServerError},
			)
			return
		}
	}

	if e != nil {
		if e.Type == "served" || e.Type == "partial_refund" {
			// Logging error in sentry
			opts := &utils.TrackErrorOpts{
				Application: authorizedApp,
				Tags:        map[string]string{"auth_type": "application"},
				Level:       sentry.LevelInfo,
			}
			utils.TrackError(c, errors.New(lang.PaymentAlreadyServed), opts)
			c.JSON(http.StatusNotAcceptable, dto.GeneralMessage{Detail: lang.PaymentAlreadyServed})
			return
		}

		if e.Type == "internal_cancellation" {
			msg := "Impossible preset gas pump"

			// Logging error in sentry
			opts := &utils.TrackErrorOpts{
				Application: authorizedApp,
				Tags:        map[string]string{"auth_type": "application"},
				Level:       sentry.LevelInfo,
			}
			utils.TrackError(c, errors.New(msg), opts)
			c.JSON(http.StatusNotAcceptable, dto.GeneralMessage{Detail: msg})
			return
		}
	}
	realAmountCharged := body.AmountCharged
	// Charge that will be necessary in order to complete $10 MXN
	var chargeFee float32

	if body.Type == "served" {
		if realAmountCharged > payment.Amount {
			// Logging error in sentry
			opts := &utils.TrackErrorOpts{
				Application: authorizedApp,
				Tags:        map[string]string{"auth_type": "application"},
				Level:       sentry.LevelInfo,
			}
			utils.TrackError(c, errors.New(lang.AmountChargedGratherThanPaid), opts)
			c.JSON(
				http.StatusNotAcceptable,
				dto.GeneralMessage{Detail: lang.AmountChargedGratherThanPaid},
			)
			return
		}

		if realAmountCharged < minChargeAmount {
			chargeFee = minChargeAmount - realAmountCharged
			realAmountCharged = minChargeAmount
		}

		payment.RealAmountReported = realAmountCharged
		payment.ChargeFee = chargeFee

		_, err = pc.repository.UpdateByID(payment.ID, payment)
		if err != nil {
			// TODO: Log error in sentry
			// Logging error in sentry
			opts := &utils.TrackErrorOpts{
				Application: authorizedApp,
				Tags:        map[string]string{"auth_type": "application"},
				Level:       sentry.LevelInfo,
			}
			utils.TrackError(c, err, opts)
			c.JSON(
				http.StatusInternalServerError,
				dto.GeneralMessage{Detail: lang.InternalServerError},
			)
			return
		}

	}

	event := &models.PaymentEvent{
		PaymentID:               payment.ID,
		Type:                    body.Type,
		AuthorizedApplicationID: &authorizedApp.ID,
	}

	err = pc.repository.CreateEvent(event)
	if err != nil {
		// Logging error in sentry
		opts := &utils.TrackErrorOpts{
			Application: authorizedApp,
			Tags:        map[string]string{"auth_type": "application"},
			Level:       sentry.LevelInfo,
		}
		utils.TrackError(c, err, opts)
		c.JSON(http.StatusInternalServerError, dto.GeneralMessage{Detail: lang.InternalServerError})
		return
	}

	// TODO: check totals to refund
	difference := payment.Amount - realAmountCharged

	closeChannel := true
	// checking if some money should be returned
	if body.Type == "served" && difference > 0 {
		closeChannel = false
		//  Refund money stripe

		if payment.PaymentProvider == "stripe" {

			_, err = pc.stripeService.MakeARefund(
				payment.ExternalTransactionID,
				float64(difference),
			)
			if err != nil {
				// TODO: Log error in sentry
				// Logging error in sentry
				opts := &utils.TrackErrorOpts{
					Application: authorizedApp,
					Tags:        map[string]string{"auth_type": "application"},
					Level:       sentry.LevelInfo,
				}
				utils.TrackError(c, err, opts)
				c.JSON(
					http.StatusInternalServerError,
					dto.GeneralMessage{Detail: lang.InternalServerError},
				)
				return
			}
		}
	}

	if (payment.PaymentProvider == "swit" || payment.PaymentProvider == "debit") &&
		body.Type == "served" {
		if difference > 0 {
			event := &models.PaymentEvent{
				PaymentID: payment.ID,
				Type:      "partial_refund",
			}
			pc.repository.CreateEvent(event)

			payment.RefundedAmount = difference

			pc.repository.UpdateByID(payment.ID, payment)

		}
		var err error
		if payment.PaymentProvider == "swit" {
			err = pc.switService.ConfirmFundReservation(
				payment.ExternalTransactionID,
				realAmountCharged,
			)
		} else if payment.PaymentProvider == "debit" {
			err = pc.debitService.PaymentConfirmation(payment.ExternalTransactionID, realAmountCharged)
		}
		if err != nil {
			// Logging error in sentry
			opts := &utils.TrackErrorOpts{
				Application: authorizedApp,
				Tags:        map[string]string{"auth_type": "application"},
			}
			utils.TrackError(c, err, opts)
		}
	}

	// POints in GM

	if body.Type == "served" {
		if !*payment.FromOperations {
			points, err := pc.socioSmartService.AccumPoints(payment)
			if err != nil {
				// TODO: Report in sentry
				// Logging error in sentry
				opts := &utils.TrackErrorOpts{
					Application: authorizedApp,
					Tags:        map[string]string{"auth_type": "application"},
				}
				utils.TrackError(c, err, opts)
			}

			payment.GMPoints = points.Amount
			payment.GMID = points.Id

			pc.repository.UpdateByID(payment.ID, payment)

			// Send email

			requestFuelSchemaMail := &schemas.FuelRequest{}
			requestFuelSchemaMail.FillData(payment)
			requestFuelSchemaMail.RefundedAmount = difference

			opts := services.SendMailOpts{
				Data:         requestFuelSchemaMail,
				Description:  "Tu carga ha sido completada",
				TemplatePath: "fuel_request.html",
				To:           payment.Customer.Email,
			}

			pc.mailService.SendMail(opts)
		}
		// Setup
		setting, err := pc.settingsRepo.GetByName("gas_pump_status")
		if err != nil {
			var csmErr error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				csmErr = errors.New("Gas Pump Status not setted")
			} else {
				csmErr = err
			}
			// Logging error in sentry
			opts := &utils.TrackErrorOpts{
				Tags: map[string]string{"webhook": "stripe"},
			}
			utils.TrackError(c, csmErr, opts)

		}

		gasPumpEnabled := false

		if setting != nil && setting.Value == "enabled" {
			gasPumpEnabled = true
		}

		if gasPumpEnabled {
			// Post ticket to GM

			optsTrans := services.ReportTransactionOpts{
				Ip:     payment.GasPump.GasStation.Ip,
				Number: payment.GasPump.Number,
			}
			err := pc.socioSmartService.ReportTransaction(optsTrans)
			if err != nil {
				// Logging error in sentry
				opts := &utils.TrackErrorOpts{
					Application: authorizedApp,
					Tags:        map[string]string{"auth_type": "application"},
					Level:       sentry.LevelError,
				}
				utils.TrackError(c, err, opts)

			}
		}

	}

	// Real time notification
	channel := paymentWebsocket.GetChannel(payment.ID.String())
	channel.BroadcastJson(dto.PaymentWebsocketNotification{Status: body.Type})

	if body.Type == "served" && closeChannel {
		paymentWebsocket.DeleteChannel(payment.ID.String())
	}

	// Implement logic for type in order to notify

	c.JSON(http.StatusCreated, dto.GeneralMessage{Detail: lang.RecordCreated})
}

// @Summary Payment Detail Customer
// @Description Payment detail for customer
// @Tags Payments
// @Produce json
// @Router /api/v1/payments/{id}/customer-detail [GET]
// @Param Authorization header string true "Token"
// @Param id path string true "uuid4 id" minLength(36) maxLength(36)
// @Success 200 {object} dto.PaymentDetailCustomer "Payment Detail for customer"
// @Failure 400 {array} dto.BadRequestMessage "Bad Request, failed on body, form, query..."
// @Failure 401 {object} dto.GeneralMessage "Unauthorized"
// @Failure 404 {object} dto.GeneralMessage "Payment id not found in db, payment not paid yet"
// @Failure 500 {object} dto.GeneralMessage "Internal server error"
func (pc *paymentController) GetByIDForCustomer(c *gin.Context) {
	var path dto.PaymentDetailCustomerPath
	if err := c.ShouldBindUri(&path); err != nil {
		c.JSON(http.StatusBadRequest, utils.MapValidatorError[dto.PaymentDetailCustomerPath](err))
		return
	}
	id, _ := uuid.Parse(path.ID)

	customer := c.MustGet("customer").(*models.Customer)

	payment, err := pc.repository.GetByIDForCustomer(id, customer.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.GeneralMessage{Detail: lang.NotFoundRecord})
			return

		}
		// Logging error in sentry
		opts := &utils.TrackErrorOpts{
			Customer: customer,
			Tags:     map[string]string{"auth_type": "customer"},
		}
		utils.TrackError(c, err, opts)
		c.JSON(http.StatusInternalServerError, dto.GeneralMessage{Detail: lang.InternalServerError})
		return
	}

	var response dto.PaymentDetailCustomer

	copier.Copy(&response, &payment)

	c.JSON(http.StatusOK, response)
}

func (pc *paymentController) PaymentNotifierWS(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	var path dto.PaymentDetailCustomerPathWebsocket

	if err := c.ShouldBindUri(&path); err != nil {
		conn.WriteMessage(
			websocket.ClosePolicyViolation,
			[]byte("not valid form"),
		)
		conn.Close()
		return
	}

	// TODO: Replace for customer one
	id, _ := uuid.Parse(path.ID)
	payment, err := pc.repository.GetByID(id)
	if err != nil {
		conn.WriteMessage(
			websocket.ClosePolicyViolation,
			[]byte("not valid payment or not found"),
		)
		conn.Close()
		return
	}

	channel := paymentWebsocket.GetChannel(payment.ID.String())

	channel.AddClient(conn)

	// Stablishing timeout in order to renew websocket
	// Kepping alive connection for a couple of minutes
	after := time.After(time.Minute * 5)

	go func() {
		select {
		case <-after:
			channel.CloseClientConnection(conn)
		}
	}()
}

// @Summary Payment provider
// @Description Payment provider for stripe or swit
// @Tags Payments
// @Produce json
// @Router /api/v1/payments/provider [GET]
// @Param Authorization header string true "Token"
// @Success 200 {object} dto.GetPaymentProviderResponse "Payment Provider"
// @Failure 401 {object} dto.GeneralMessage "Unauthorized"
func (pc *paymentController) GetPaymentProvider(c *gin.Context) {
	customer := c.MustGet("customer").(*models.Customer)

	setting, err := pc.settingsRepo.GetByName("payment_provider")
	if err != nil {
		var csmErr error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			csmErr = errors.New("Payment Provider not setted")
		} else {
			csmErr = err
		}
		// Logging error in sentry
		opts := &utils.TrackErrorOpts{
			Customer: customer,
			Tags:     map[string]string{"auth_type": "customer"},
		}
		utils.TrackError(c, csmErr, opts)
		c.JSON(http.StatusInternalServerError, dto.GeneralMessage{Detail: lang.InternalServerError})
		return

	}

	paymentProvider := setting.Value

	setup := dto.GetPaymentProviderResponse{
		PaymentProvider: paymentProvider,
	}

	if paymentProvider == "swit" {

		setup.Business = pc.config.SwitBusiness
		setup.Token = pc.config.SwitToken
		setup.CustomerID = customer.SwitCustomerID
	}
	c.JSON(http.StatusOK, setup)
}

// @Summary Invoice payment
// @Description Make invoicing for a customer payment
// @Tags Payments
// @Produce json
// @Router /api/v1/payments/invoicing/{id} [POST]
// @Param id path string true "uuid4 id" minLength(36) maxLength(36)
// @Param Authorization header string true "Token"
// @Param body body dto.SignInvoiceRequest true "Body for invocing"
// @Success 200 {object} dto.SignInvoiceResponse "Response with uuid"
// @Failure 404 {object} dto.GeneralMessage "Not found"
// @Failure 406 {object} dto.GeneralMessage "Not accepted"
// @Failure 409 {object} dto.GeneralMessage "Already invoiced"
// @Failure 412 {object} dto.GeneralMessage "Data sent bad"
// @Failure 401 {object} dto.GeneralMessage "Unauthorized"
func (pc *paymentController) SignInvoice(c *gin.Context) {
	var body dto.SignInvoiceRequest
	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, utils.MapValidatorError[dto.SignInvoiceRequest](err))
		return
	}

	customer := c.MustGet("customer").(*models.Customer)

	id, _ := uuid.Parse(c.Param("id"))

	payment, err := pc.repository.GetByIDForCustomer(id, customer.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, dto.GeneralMessage{Detail: lang.NotFoundRecord})
			return
		}
		// Logging error in sentry
		opts := &utils.TrackErrorOpts{
			Customer: customer,
			Tags:     map[string]string{"auth_type": "customer"},
		}
		utils.TrackError(c, err, opts)
		c.JSON(http.StatusInternalServerError, dto.GeneralMessage{Detail: lang.InternalServerError})
		return
	}

	if payment.Invoiced {
		c.JSON(http.StatusConflict, dto.GeneralMessage{Detail: "Payment already invoiced"})
		return
	}

	if len(payment.Events) == 0 {
		c.JSON(
			http.StatusNotAcceptable,
			dto.GeneralMessage{Detail: "Load not finished or payment required"},
		)
		return
	}

	if payment.Events[0].Type != "served" && payment.Events[0].Type != "partial_refund" {
		c.JSON(
			http.StatusNotAcceptable,
			dto.GeneralMessage{Detail: "Load not finished or payment required"},
		)
		return
	}

	opts := services.SignInvoiceOpts{
		Params:  body,
		Payment: payment,
	}
	data, err := pc.invoicingService.SignInvoice(opts)
	if err != nil {
		// Logging error in sentry
		opts := &utils.TrackErrorOpts{
			Customer: customer,
			Tags:     map[string]string{"auth_type": "customer"},
		}
		if err == services.ErrData {
			utils.TrackError(c, errors.New(data), opts)
			c.JSON(http.StatusPreconditionFailed, dto.GeneralMessage{Detail: data})
			return
		}
		utils.TrackError(c, err, opts)
		c.JSON(http.StatusInternalServerError, dto.GeneralMessage{Detail: lang.InternalServerError})
		return

	}

	payment.Invoiced = true
	payment.InvoiceID = data

	pc.repository.UpdateByID(payment.ID, payment)

	c.JSON(http.StatusOK, dto.SignInvoiceResponse{UUID: data})
}

// @Summary Invoice payment Resend
// @Description Resend invoice previously invoiced
// @Tags Payments
// @Produce json
// @Router /api/v1/payments/invoicing/{id}/resend [POST]
// @Param id path string true "uuid4 id" minLength(36) maxLength(36)
// @Param Authorization header string true "Token"
// @Param body body dto.ResendInvoiceRequest true "Body for invocing"
// @Success 200 {object} dto.GeneralMessage "Created"
// @Failure 404 {object} dto.GeneralMessage "Not found"
// @Failure 401 {object} dto.GeneralMessage "Unauthorized"
// @Failure 500 {object} dto.GeneralMessage "Internal Server Error"
func (pc *paymentController) ResendInvoice(c *gin.Context) {
	var body dto.ResendInvoiceRequest
	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, utils.MapValidatorError[dto.ResendInvoiceRequest](err))
		return
	}

	customer := c.MustGet("customer").(*models.Customer)

	id, _ := uuid.Parse(c.Param("id"))

	payment, err := pc.repository.GetByIDForCustomer(id, customer.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, dto.GeneralMessage{Detail: lang.NotFoundRecord})
			return
		}
		// Logging error in sentry
		opts := &utils.TrackErrorOpts{
			Customer: customer,
			Tags:     map[string]string{"auth_type": "customer"},
		}
		utils.TrackError(c, err, opts)
		c.JSON(http.StatusInternalServerError, dto.GeneralMessage{Detail: lang.InternalServerError})
		return
	}

	if !payment.Invoiced {
		c.JSON(http.StatusConflict, dto.GeneralMessage{Detail: "Payment not invoiced yet"})
		return
	}

	err = pc.invoicingService.ResendInvoice(payment.InvoiceID, body.Email)

	c.JSON(http.StatusOK, dto.GeneralMessage{Detail: "OK"})
}

// @Summary Get PDF from invoice
// @Description Get PDF File from already generated invoice
// @Tags Payments
// @Produce json
// @Router /api/v1/payments/invoicing/{id}/pdf [GET]
// @Param id path string true "uuid4 id" minLength(36) maxLength(36)
// @Param Authorization header string true "Token"
// @Success 200 {object} dto.GetInvoicePDFResponse "PDF"
// @Failure 404 {object} dto.GeneralMessage "Not found"
// @Failure 401 {object} dto.GeneralMessage "Unauthorized"
// @Failure 500 {object} dto.GeneralMessage "Internal Server Error"
func (pc *paymentController) GetInvoicePDF(c *gin.Context) {
	customer := c.MustGet("customer").(*models.Customer)

	id, _ := uuid.Parse(c.Param("id"))

	payment, err := pc.repository.GetByIDForCustomer(id, customer.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, dto.GeneralMessage{Detail: lang.NotFoundRecord})
			return
		}
		// Logging error in sentry
		opts := &utils.TrackErrorOpts{
			Customer: customer,
			Tags:     map[string]string{"auth_type": "customer"},
		}
		utils.TrackError(c, err, opts)
		c.JSON(http.StatusInternalServerError, dto.GeneralMessage{Detail: lang.InternalServerError})
		return
	}

	if !payment.Invoiced {
		c.JSON(http.StatusConflict, dto.GeneralMessage{Detail: "Payment not invoiced yet"})
		return
	}

	url, err := pc.invoicingService.GetInvoicePDF(payment.InvoiceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.GeneralMessage{Detail: lang.InternalServerError})
		return
	}

	c.JSON(http.StatusOK, dto.GetInvoicePDFResponse{UrlPDF: url})
}

// @Summary Do payment action
// @Description Manual action for payment
// @Tags Payments
// @Produce json
// @Router /api/v1/payments/actions/{id} [POST]
// @Param id path string true "uuid4 id" minLength(36) maxLength(36)
// @Security Bearer
// @Param body body dto.DoPaymentActionRequest true "Do action for payment"
// @Success 200 {object} dto.GeneralMessage "done"
// @Failure 404 {object} dto.GeneralMessage "Not found"
// @Failure 401 {object} dto.GeneralMessage "Unauthorized"
// @Failure 500 {object} dto.GeneralMessage "Internal Server Error"
func (pc *paymentController) DoPaymentAction(c *gin.Context) {
	var path dto.DoPaymentActionRequestPath
	if err := c.ShouldBindUri(&path); err != nil {
		c.JSON(http.StatusBadRequest, utils.MapValidatorError[dto.DoPaymentActionRequestPath](err))
		return
	}

	id, _ := uuid.Parse(path.ID)

	var body dto.DoPaymentActionRequest
	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, utils.MapValidatorError[dto.DoPaymentActionRequest](err))
		return
	}

	user := c.MustGet("user").(*models.User)

	payment, err := pc.repository.GetByIDPreloaded(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, dto.GeneralMessage{Detail: lang.NotFoundRecord})
			return
		}
		// Logging error in sentry
		opts := &utils.TrackErrorOpts{
			Admin: user,
			Tags:  map[string]string{"auth_type": "admin"},
		}
		utils.TrackError(c, err, opts)
		c.JSON(http.StatusInternalServerError, dto.GeneralMessage{Detail: lang.InternalServerError})
		return
	}

	e, err := pc.repository.GetLastEventByPaymentID(payment.ID)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			// Logging error in sentry
			opts := &utils.TrackErrorOpts{
				Admin: user,
				Tags:  map[string]string{"auth_type": "admin"},
				Level: sentry.LevelInfo,
			}
			utils.TrackError(c, err, opts)
			c.JSON(
				http.StatusInternalServerError,
				dto.GeneralMessage{Detail: lang.InternalServerError},
			)
			return
		}
	}

	if e == nil {
		c.JSON(http.StatusPaymentRequired, dto.GeneralMessage{Detail: lang.PaymentRequired})
		return
	}

	if e.Type != "paid" && e.Type != "funds_reserved" {
		c.JSON(
			http.StatusPreconditionFailed,
			dto.GeneralMessage{Detail: "payment status not valid"},
		)
		return
	}

	if body.Action == "refund" {
		var err error
		if payment.PaymentProvider == "stripe" {
			_, err = pc.stripeService.MakeARefund(payment.ExternalTransactionID, -3)
		} else if payment.PaymentProvider == "swit" {
			err = pc.switService.CancelFundReservation(payment.ExternalTransactionID)
		} else if payment.PaymentProvider == "debit" {
			err = pc.debitService.CancelReservation(payment.ExternalTransactionID)
		}

		// TODO: Log error in sentry

		if err != nil {
			// Logging error in sentry
			opts := &utils.TrackErrorOpts{
				Admin: user,
				Tags: map[string]string{
					"auth_type":        "admin",
					"payment_provider": payment.PaymentProvider,
				},
				Level: sentry.LevelInfo,
			}
			utils.TrackError(c, err, opts)
			c.JSON(
				http.StatusInternalServerError,
				dto.GeneralMessage{Detail: lang.InternalServerError},
			)
			return
		}

		if !*payment.FromOperations {
			requestFuelSchemaMail := &schemas.FuelRequest{}
			requestFuelSchemaMail.FillData(payment)
			requestFuelSchemaMail.RefundedAmount = payment.Amount
			requestFuelSchemaMail.Error = true

			opts := services.SendMailOpts{
				Data:         requestFuelSchemaMail,
				Description:  "Hubo un error al intentar hacer tu carga (manual)",
				TemplatePath: "fuel_request.html",
				To:           payment.Customer.Email,
			}

			pc.mailService.SendMail(opts)

		}

		if payment.PaymentProvider != "stripe" {
			pc.repository.CreateEvent(
				&models.PaymentEvent{PaymentID: payment.ID, Type: "manual_action"},
			)
			pc.repository.CreateEvent(
				&models.PaymentEvent{PaymentID: payment.ID, Type: "partial_refund"},
			)
		}

		payment.RefundedAmount = payment.Amount

		pc.repository.UpdateByID(payment.ID, payment)

		c.JSON(http.StatusOK, dto.GeneralMessage{Detail: "ok"})
		return
	} else {
		channel := paymentWebsocket.GetChannel(payment.ID.String())

		channel.BroadcastJson(dto.PaymentWebsocketNotification{Status: "paid"})
		// Setup
		setting, err := pc.settingsRepo.GetByName("gas_pump_status")
		if err != nil {
			var csmErr error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				csmErr = errors.New("Gas Pump Status not setted")
			} else {
				csmErr = err
			}
			// Logging error in sentry
			opts := &utils.TrackErrorOpts{
				Tags:  map[string]string{"auth_type": "admin"},
				Admin: user,
			}
			utils.TrackError(c, csmErr, opts)
		}

		gasPumpEnabled := false

		if setting != nil && setting.Value == "enabled" {
			gasPumpEnabled = true
		}

		if gasPumpEnabled { // PRE-SET gas pump
			opts := services.SetGasPumpOptions{
				Number:    payment.GasPump.Number,
				Ip:        payment.GasPump.GasStation.Ip,
				FuelType:  payment.FuelType,
				Amount:    payment.Amount,
				PaymentID: payment.ID,
				Discount:  payment.DiscountPerLiter,
			}
			data, err := pc.socioSmartService.SetGasPump(opts)
			status := "Error"

			if data != nil && data.Status == 0 {
				status = "This pump has already a preset"
			}

			if err != nil || data.Status == 0 {
				// Logging error in sentry
				optsTE := &utils.TrackErrorOpts{
					Tags:  map[string]string{"auth_type": "admin"},
					Admin: user,
					Context: map[string]map[string]any{
						"Response": {"Status": status},
						"Data": {
							"gas_station":     payment.GasPump.GasStation.Name,
							"gas_station_id":  payment.GasPump.GasStationID,
							"gas_pump_number": payment.GasPump.Number,
							"gas_pump_id":     payment.GasPumpID,
							"discount":        opts.Discount,
						},
					},
				}
				utils.TrackError(c, err, optsTE)
				if err != nil {
					c.JSON(http.StatusInternalServerError, dto.GeneralMessage{Detail: lang.InternalServerError})
					return
				}

				c.JSON(http.StatusServiceUnavailable, dto.GeneralMessage{Detail: status})
				return
			} else {
				event := &models.PaymentEvent{
					PaymentID: payment.ID,
					Type:      "pump_ready",
				}

				err = pc.repository.CreateEvent(event)
				if err != nil {
					// TODO: Log error in sentry
					opts := &utils.TrackErrorOpts{
						Tags: map[string]string{"auth_type": "admin"},
					}
					utils.TrackError(c, err, opts)
				}
				channel.BroadcastJson(dto.PaymentWebsocketNotification{Status: "pump_ready"})
				channel.BroadcastJson(dto.PaymentWebsocketNotification{Status: "pump_ready"})
				c.JSON(http.StatusOK, dto.GeneralMessage{Detail: "ok"})
				return
			}
		} else {
			c.JSON(http.StatusPreconditionRequired, dto.GeneralMessage{Detail: "gas pump not enabled, please contact with admin"})
			return
		}
	}
}
