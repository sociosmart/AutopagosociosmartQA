package middlewares

import (
	"errors"
	"net/http"
	"smartgas-payment/internal/dto"
	"smartgas-payment/internal/lang"
	"smartgas-payment/internal/models"
	"smartgas-payment/internal/repository"
	"smartgas-payment/internal/services"
	"smartgas-payment/internal/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type CustomerAuthMiddleware struct {
	customerRepository    repository.CustomerRepository
	extCustomerService    services.CustomerService
	stripeService         services.StripeService
	switService           services.SwitService
	elegibilityRepository repository.ElegibilityRepository
}

func ProvideCustomerAUthMiddleware(
	customerRepository repository.CustomerRepository,
	extCustomerService services.CustomerService,
	stripeService services.StripeService,
	switService services.SwitService,
	elegibilityRepository repository.ElegibilityRepository,
) *CustomerAuthMiddleware {
	return &CustomerAuthMiddleware{
		customerRepository:    customerRepository,
		extCustomerService:    extCustomerService,
		stripeService:         stripeService,
		switService:           switService,
		elegibilityRepository: elegibilityRepository,
	}
}

func (cm *CustomerAuthMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var customerAuthHeader dto.CustomerAuthorizationHeader

		if err := c.ShouldBindHeader(&customerAuthHeader); err != nil {
			c.JSON(http.StatusUnauthorized, dto.GeneralMessage{Detail: lang.NoAuthorizationHeader})
			c.Abort()
			return
		}

		splittedHeader := strings.Split(customerAuthHeader.Authorization, " ")

		if len(splittedHeader) < 2 || splittedHeader[0] != "Token" {
			c.JSON(
				http.StatusUnauthorized,
				dto.GeneralMessage{Detail: lang.AuthorizationHeaderMalformed},
			)
			c.Abort()
			return
		}

		token := splittedHeader[1]

		customer, err := cm.extCustomerService.Verify(token)
		if err != nil {
			if err.Error() == lang.InvalidOrExpiredToken {
				c.JSON(
					http.StatusUnauthorized,
					dto.GeneralMessage{Detail: lang.InvalidOrExpiredToken},
				)
				c.Abort()
				return
			}

			// Logging error for when swit list cards fails
			opts := &utils.TrackErrorOpts{
				Tags: map[string]string{"scope": "customer_auth_middleware"},
			}
			utils.TrackError(c, err, opts)
			c.JSON(
				http.StatusInternalServerError,
				dto.GeneralMessage{Detail: lang.InternalServerError},
			)
			c.Abort()
			return
		}

		customerInDB, created, err := cm.customerRepository.GetCustomerOrCreate(customer)
		if err != nil {
			// Logging error for when swit list cards fails
			opts := &utils.TrackErrorOpts{
				Tags: map[string]string{"scope": "customer_auth_middleware"},
			}
			utils.TrackError(c, err, opts)
			c.JSON(
				http.StatusInternalServerError,
				dto.GeneralMessage{Detail: lang.InternalServerError},
			)
			c.Abort()
			return
		}

		if created {
			filters := map[string]any{
				"min_amount":  0,
				"min_charges": 0,
				"active":      true,
			}

			level, err := cm.elegibilityRepository.GetLevelByCriterias(filters)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					opts := &utils.TrackErrorOpts{
						Tags:     map[string]string{"scope": "customer_auth_middleware"},
						Customer: customerInDB,
					}
					utils.TrackError(
						c,
						errors.New(
							"You should create a level record in the DB with min_amount in 0 and min_charges in 0 in order to assign to new customers",
						),
						opts,
					)
					c.JSON(
						http.StatusInternalServerError,
						dto.GeneralMessage{Detail: lang.InternalServerError},
					)
				}
			}

			if level != nil {
				loc, _ := time.LoadLocation("America/Mazatlan")
				now := time.Now().In(loc)
				cusLevelM := models.CustomerLevel{
					CustomerID:    &customerInDB.ID,
					LevelID:       &level.ID,
					ValidityMonth: utils.IntAddr(int(now.Month())),
					ValidityYear:  utils.IntAddr(now.Year()),
				}
				cm.elegibilityRepository.CreateCustomerLevel(&cusLevelM)
			}
		}

		// Update user to sync with GM data
		var toUpdate models.Customer
		copier.Copy(&toUpdate, &customer)
		_, err = cm.customerRepository.UpdateByID(customerInDB.ID, &toUpdate)
		if err != nil {
			// Logging error for when swit list cards fails
			opts := &utils.TrackErrorOpts{
				Tags: map[string]string{"scope": "customer_auth_middleware"},
			}
			utils.TrackError(c, err, opts)
			c.JSON(
				http.StatusInternalServerError,
				dto.GeneralMessage{Detail: lang.InternalServerError},
			)
			c.Abort()
			return
		}

		// create stripe's customer
		if customerInDB.StripeCustomerID == "" {
			cusId, err := cm.stripeService.CreateCustomer(customer)

			// TODO: Check what to do with the error
			if err == nil {
				customerInDB.StripeCustomerID = cusId

				cm.customerRepository.UpdateByID(customerInDB.ID, customerInDB)
			} else {
				// Logging error for when swit list cards fails
				opts := &utils.TrackErrorOpts{
					Tags: map[string]string{"scope": "customer_auth_middleware"},
				}
				utils.TrackError(c, err, opts)
			}

		}

		if customerInDB.SwitCustomerID == "" {
			cusId, err := cm.switService.CreateCustomer(customer)

			if err == nil {
				customerInDB.SwitCustomerID = cusId

				cm.customerRepository.UpdateByID(customerInDB.ID, customerInDB)
			} else {
				// Logging error for when swit list cards fails
				opts := &utils.TrackErrorOpts{
					Tags: map[string]string{"scope": "customer_auth_middleware"},
				}
				utils.TrackError(c, err, opts)

			}
		}

		c.Set("customer", customerInDB)
		c.Next()
	}
}
