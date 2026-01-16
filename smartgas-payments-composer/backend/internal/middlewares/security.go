package middlewares

import (
	"errors"
	"net/http"
	"smartgas-payment/internal/dto"
	"smartgas-payment/internal/lang"
	"smartgas-payment/internal/repository"
	"smartgas-payment/internal/services"
	"smartgas-payment/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SecurityMiddleware struct {
	repository           repository.SecurityRepository
	gasStationRepository repository.GasStationRepository
	smartGasServicee     services.SocioSmartService
}

func ProvideSecurityMiddleware(
	repository repository.SecurityRepository,
	gasStationRepository repository.GasStationRepository,
	smartGasService services.SocioSmartService,
) *SecurityMiddleware {
	return &SecurityMiddleware{
		repository:           repository,
		gasStationRepository: gasStationRepository,
		smartGasServicee:     smartGasService,
	}
}

func (sm *SecurityMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var headers dto.SecurityHeadersRequest

		if err := c.ShouldBindHeader(&headers); err != nil {
			c.JSON(http.StatusBadRequest, utils.MapValidatorError[dto.SecurityHeadersRequest](err))
			c.Abort()
			return
		}

		appKey, _ := uuid.Parse(headers.AppKey)
		apiKey, _ := uuid.Parse(headers.ApiKey)

		authorizedApp, err := sm.repository.GetByKeys(appKey, apiKey)
		if err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(
					http.StatusUnauthorized,
					dto.GeneralMessage{Detail: lang.ApplicationUnauthorized},
				)
				c.Abort()
				return
			}

			c.JSON(
				http.StatusInternalServerError,
				dto.GeneralMessage{Detail: lang.InternalServerError},
			)
			c.Abort()
			return
		}

		c.Set("application", authorizedApp)

		c.Next()
	}
}

func (sm *SecurityMiddleware) SmartGasEmployeeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var headers dto.EmployeeValidationHeadersRequest

		if err := c.ShouldBindHeader(&headers); err != nil {
			c.JSON(http.StatusBadRequest, utils.MapValidatorError[dto.SecurityHeadersRequest](err))
			c.Abort()
			return
		}
		gasStationID := headers.ExternalGasStationID
		station, err := sm.gasStationRepository.GetByExternalID(gasStationID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, dto.GeneralMessage{Detail: lang.GasStationNotFound})
				c.Abort()
				return
			}
			opts := &utils.TrackErrorOpts{
				Tags: map[string]string{"scope": "employee_authentication"},
			}
			utils.TrackError(c, err, opts)
			c.JSON(
				http.StatusInternalServerError,
				dto.GeneralMessage{Detail: lang.InternalServerError},
			)
			c.Abort()
			return
		}
		opts := services.ValidateEmployeeOpts{
			ExternalGasStationID: station.ExternalID,
			EmployeeID:           headers.EmployeeID,
			EmployeeNIP:          headers.EmployeeNIP,
			GasPumpCRE:           station.CrePermission,
		}
		authorized, err := sm.smartGasServicee.ValidateEmployee(opts)
		if err != nil {
			opts := &utils.TrackErrorOpts{
				Tags: map[string]string{"scope": "employee_authentication"},
			}
			utils.TrackError(c, err, opts)
			c.JSON(
				http.StatusInternalServerError,
				dto.GeneralMessage{Detail: lang.InternalServerError},
			)
			c.Abort()
			return
		}

		if !authorized {
			c.JSON(
				http.StatusUnauthorized,
				dto.GeneralMessage{Detail: lang.UnauthorizedEmployee},
			)
			c.Abort()
			return
		}
		c.Set("gas_station", station)
		c.Set("employee_id", headers.EmployeeID)

		c.Next()
	}
}
