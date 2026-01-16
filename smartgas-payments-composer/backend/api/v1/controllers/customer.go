package controllers

import (
	"errors"
	"net/http"
	"smartgas-payment/config"
	"smartgas-payment/internal/dto"
	"smartgas-payment/internal/lang"
	"smartgas-payment/internal/models"
	"smartgas-payment/internal/repository"
	"smartgas-payment/internal/schemas"
	"smartgas-payment/internal/services"
	"smartgas-payment/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type CustomerController interface {
	ListPaymenthMethods(*gin.Context)
	ListPaymenthMethodsSwit(*gin.Context)
	DeleteCard(*gin.Context)
	ListAll(*gin.Context)
	GetElegibilityLevel(*gin.Context)
}

type customerController struct {
	stripeService   services.StripeService
	switService     services.SwitService
	config          config.Config
	settingsRepo    repository.SettingRepository
	repository      repository.CustomerRepository
	elegibilityRepo repository.ElegibilityRepository
}

func ProvideCustomerController(
	stripeService services.StripeService,
	switService services.SwitService,
	config config.Config,
	settingsRepo repository.SettingRepository,
	repository repository.CustomerRepository,
	elegibilityRepo repository.ElegibilityRepository,
) *customerController {
	return &customerController{
		stripeService:   stripeService,
		switService:     switService,
		config:          config,
		settingsRepo:    settingsRepo,
		repository:      repository,
		elegibilityRepo: elegibilityRepo,
	}
}

// @Summary Customer cards SWIT
// @Description Get customer's payment methods swit
// @Tags Customers
// @Produce json
// @Router /api/v1/customers/payment-methods-swit [GET]
// @Param Authorization header string true "Token"
// @Failure 401 {object} dto.GeneralMessage "Unauthorized"
// @Success 200 {object} dto.ListCustomerPaymentMethodResponse "User's profile"
func (cc *customerController) ListPaymenthMethodsSwit(c *gin.Context) {
	customer := c.MustGet("customer").(*models.Customer)

	paymentMethods, err := cc.switService.ListCardsByCustomer(customer.SwitCustomerID)
	if err != nil {
		// Logging error for when swit list cards fails
		opts := &utils.TrackErrorOpts{
			Context:  map[string]map[string]any{"Payment Provider": {"name": "swit"}},
			Tags:     map[string]string{"auth_type": "customer"},
			Customer: customer,
		}
		utils.TrackError(c, err, opts)
	}

	paymentMethodsResponse := make(dto.ListCustomerPaymentMethodResponse, 0)

	for _, pm := range paymentMethods {
		var paymentMethod schemas.PaymentMethod

		paymentMethod.ID = pm.ID

		paymentMethod.Card.Brand = "default"
		paymentMethod.Card.Last4 = pm.Last4
		paymentMethod.IsLastUsed = pm.IsLastUsed

		paymentMethodsResponse = append(paymentMethodsResponse, paymentMethod)

	}

	c.JSON(http.StatusOK, paymentMethodsResponse)
}

// @Summary Delete a customer card
// @Description Delete a card from customer
// @Tags Customers
// @Produce json
// @Router /api/v1/customers/payment-methods/{card_id} [DELETE]
// @Param card_id path string true "id" minLength(1)
// @Param Authorization header string true "Token"
// @Failure 401 {object} dto.GeneralMessage "Unauthorized"
// @Success 200 {object} dto.GeneralMessage "OK if deleted"
func (ss *customerController) DeleteCard(c *gin.Context) {
	customer := c.MustGet("customer").(*models.Customer)

	cardID := c.Param("card_id")

	setting, err := ss.settingsRepo.GetByName("payment_provider")
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

	if paymentProvider == "stripe" {
		err = ss.stripeService.DeletePaymentMethod(cardID)
	} else {
		err = ss.switService.DeleteCard(customer.SwitCustomerID, cardID)
	}

	if err != nil {
		// Logging error for when swit list cards fails
		opts := &utils.TrackErrorOpts{
			Context:  map[string]map[string]any{"Payment Provider": {"name": paymentProvider}},
			Customer: customer,
			Tags:     map[string]string{"auth_type": "customer"},
		}
		utils.TrackError(c, err, opts)
		c.JSON(http.StatusInternalServerError, dto.GeneralMessage{Detail: lang.InternalServerError})
		return
	}

	c.JSON(http.StatusOK, dto.GeneralMessage{Detail: "Deleted"})
}

// @Summary Customer cards
// @Description Get customer's payment methods
// @Tags Customers
// @Produce json
// @Router /api/v1/customers/payment-methods [GET]
// @Param Authorization header string true "Token"
// @Failure 401 {object} dto.GeneralMessage "Unauthorized"
// @Success 200 {object} dto.ListCustomerPaymentMethodResponse "User's profile"
func (cc *customerController) ListPaymenthMethods(c *gin.Context) {
	customer := c.MustGet("customer").(*models.Customer)

	paymentMethods := cc.stripeService.ListPaymenthMethodsByCustomer(customer.StripeCustomerID)

	paymentMethodsResponse := make(dto.ListCustomerPaymentMethodResponse, 0)

	copier.Copy(&paymentMethodsResponse, paymentMethods)

	c.JSON(http.StatusOK, paymentMethodsResponse)
}

// @Summary List all customers
// @Description List of all customers
// @Tags Customers
// @Produce json
// @Router /api/v1/customers/all [GET]
// @Security Bearer
// @Success 200 {array} dto.ListAllCustomersResponse "Customers"
// @Failure 401 {object} dto.GeneralMessage "Unauthorized"
// @Failure 500 {object} dto.GeneralMessage "Internal server error"
func (cc *customerController) ListAll(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	customers, err := cc.repository.ListAll()
	if err != nil {
		// Logging error for when swit list cards fails
		opts := &utils.TrackErrorOpts{
			Admin: user,
			Tags:  map[string]string{"auth_type": "admin"},
		}
		utils.TrackError(c, err, opts)
		c.JSON(http.StatusInternalServerError, dto.GeneralMessage{Detail: lang.InternalServerError})
		return
	}

	customersResponse := make([]dto.ListAllCustomersResponse, 0)

	copier.Copy(&customersResponse, customers)

	c.JSON(http.StatusOK, customersResponse)
}

// @Summary List all customers
// @Description List of all customers
// @Tags Customers
// @Produce json
// @Router /api/v1/customers/level [GET]
// @Param Authorization header string true "Token"
// @Success 200 {object} dto.CustomerLevelAssignedResponse "Customer Level"
// @Failure 401 {object} dto.GeneralMessage "Unauthorized"
// @Failure 404 {object} dto.GeneralMessage "Not level in this month"
// @Failure 500 {object} dto.GeneralMessage "Internal server error"
func (cc *customerController) GetElegibilityLevel(c *gin.Context) {
	customer := c.MustGet("customer").(*models.Customer)

	now := time.Now()

	cusLevel, err := cc.elegibilityRepo.GetCustomerLevelByCriterias(map[string]any{
		"customer_id":    customer.ID,
		"validity_month": now.Month(),
		"validity_year":  now.Year(),
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.GeneralMessage{Detail: lang.NotFoundRecord})
			return
		}
		// Logging error for when swit list cards fails
		opts := &utils.TrackErrorOpts{
			Customer: customer,
			Tags:     map[string]string{"auth_type": "customer"},
		}
		utils.TrackError(c, err, opts)
		c.JSON(http.StatusInternalServerError, dto.GeneralMessage{Detail: lang.InternalServerError})
		return
	}

	var levelResponse dto.CustomerLevelAssignedResponse

	copier.Copy(&levelResponse, cusLevel.Level)

	applicablePromotionType, err := cc.settingsRepo.GetByName("applicable_promotion_type")
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

	if promotionTypeInDB == "elegibility" {
		levelResponse.LevelsEnabled = true
	}

	c.JSON(http.StatusOK, levelResponse)
}
