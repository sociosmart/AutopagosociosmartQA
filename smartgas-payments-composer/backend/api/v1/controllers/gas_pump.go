package controllers

import (
	"errors"
	"net/http"
	"smartgas-payment/internal/dto"
	"smartgas-payment/internal/lang"
	"smartgas-payment/internal/models"
	"smartgas-payment/internal/repository"
	"smartgas-payment/internal/schemas"
	"smartgas-payment/internal/tasks"
	"smartgas-payment/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type GasPumpController interface {
	List(*gin.Context)
	Get(*gin.Context)
	Update(*gin.Context)
	Create(*gin.Context)
	GetDetailForCustomer(*gin.Context)
}

type gasPumpController struct {
	repository          repository.GasPumpRepository
	synchronizationTask tasks.SynchronizationTask
	campaignRepository  repository.CampaignRepository
	settingsRepo        repository.SettingRepository
}

func ProvideGasPumpProvider(
	repository repository.GasPumpRepository,
	synchronizationTask tasks.SynchronizationTask,
	campaignRepository repository.CampaignRepository,
	settingsRepo repository.SettingRepository,
) *gasPumpController {
	return &gasPumpController{
		repository:          repository,
		synchronizationTask: synchronizationTask,
		campaignRepository:  campaignRepository,
		settingsRepo:        settingsRepo,
	}
}

// @Summary Gas Pump List
// @Description Get paginated gas pumps
// @Tags Gas Pumps
// @Produce json
// @Router /api/v1/gas-pumps [GET]
// @Security Bearer
// @Param param query dto.PaginateRequest true "Detalles"
// @Param search query string false "Search in gas stations name, ip it looks up through id, name ip"
// @Success 200 {object} dto.PaginationResponse{data=[]dto.GasPumpListResponse} "Gas station paginated"
// @Failure 400 {array} dto.BadRequestMessage "Bad Request, failed on body, form, query..."
// @Failure 401 {object} dto.GeneralMessage "Unauthorized"
// @Failure 500 {object} dto.GeneralMessage "Internal server error"
func (gp *gasPumpController) List(c *gin.Context) {
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

	stations, err := gp.repository.List(&paginationSchema, filters)
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

	pumpsResponse := make([]dto.GasPumpListResponse, 0)

	copier.Copy(&pumpsResponse, &stations)

	// Making the response
	var paginationResponse dto.PaginationResponse

	copier.Copy(&paginationResponse, &paginationSchema)

	paginationResponse.Data = pumpsResponse

	c.JSON(http.StatusOK, paginationResponse)
}

// @Summary Gas Pump Create
// @Description Create Gas Pump
// @Tags Gas Pumps
// @Produce json
// @Router /api/v1/gas-pumps [POST]
// @Security Bearer
// @Param Data body dto.GasPumpCreateRequest true "Gas Pump Body"
// @Success 201 {object} dto.GasPumpCreateResponse "Gas station create response"
// @Failure 401 {object} dto.GeneralMessage "Unauthorized"
// @Failure 409 {object} dto.GeneralMessage "Duplicated entry"
// @Failure 406 {object} dto.GeneralMessage "Not acceptable value for gas_station_id, id does not exist in DB"
// @Failure 400 {array} dto.BadRequestMessage "Bad Request, failed on body, form, query..."
// @Failure 500 {object} dto.GeneralMessage "Internal server error"
func (gp *gasPumpController) Create(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	var body dto.GasPumpCreateRequest
	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, utils.MapValidatorError[dto.GasPumpCreateRequest](err))
		return
	}

	var gasPump models.GasPump
	gasPump.CreatedBy = user

	copier.Copy(&gasPump, &body)

	if err := gp.repository.Create(&gasPump); err != nil {
		if utils.CheckDuplicatedEntry(err) {
			c.JSON(
				http.StatusConflict,
				dto.GeneralMessage{Detail: lang.DuplicatedEntry + "gas_station_id, number"},
			)
			return
		} else if utils.CheckMysqlErrCode(err, 1452) {
			c.JSON(http.StatusNotAcceptable, dto.GeneralMessage{Detail: lang.NotAcceptable + "gas_station_id"})
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

	var response dto.GasPumpCreateResponse

	copier.Copy(&response, &gasPump)

	c.JSON(http.StatusCreated, response)
}

// @Summary Gas Pump Detail
// @Description Get Gas station detail
// @Tags Gas Pumps
// @Produce json
// @Router /api/v1/gas-pumps/{id} [GET]
// @Security Bearer
// @Param id path string true "uuid4 id" minLength(36) maxLength(36)
// @Success 200 {object} dto.GasPumpGetResponse "Gas station detail"
// @Failure 400 {array} dto.BadRequestMessage "Bad Request, failed on body, form, query..."
// @Failure 404 {object} dto.GeneralMessage "Not Found"
// @Failure 401 {object} dto.GeneralMessage "Unauthorized"
// @Failure 500 {object} dto.GeneralMessage "Internal server error"
func (gp *gasPumpController) Get(c *gin.Context) {
	var pathParams dto.GasPumpGetPathRequest
	if err := c.ShouldBindUri(&pathParams); err != nil {
		c.JSON(http.StatusBadRequest, utils.MapValidatorError[dto.GasPumpGetPathRequest](err))
		return
	}

	user := c.MustGet("user").(*models.User)

	id, _ := uuid.Parse(pathParams.ID)
	station, err := gp.repository.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
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

	var gasPump dto.GasPumpGetResponse

	copier.Copy(&gasPump, &station)

	c.JSON(http.StatusOK, gasPump)
}

// @Summary Gas Pump Update
// @Description Update Gas station
// @Tags Gas Pumps
// @Produce json
// @Router /api/v1/gas-pumps/{id} [PUT]
// @Security Bearer
// @Param id path string true "uuid4 id" minLength(36) maxLength(36)
// @Param Data body dto.GasPumpUpdateRequest true "Gas station values"
// @Success 200 {object} dto.GeneralMessage "Record Updated"
// @Failure 400 {array} dto.BadRequestMessage "Bad Request, failed on body, form, query..."
// @Failure 404 {object} dto.GeneralMessage "Not Found"
// @Failure 401 {object} dto.GeneralMessage "Unauthorized"
// @Failure 406 {object} dto.GeneralMessage "Not acceptable value for gas_station_id, id does not exist in DB"
// @Failure 409 {object} dto.GeneralMessage "Duplicated entry"
// @Failure 500 {object} dto.GeneralMessage "Internal server error"
func (gp *gasPumpController) Update(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	var path dto.GasPumpUpdatePathRequest

	if err := c.ShouldBindUri(&path); err != nil {
		c.JSON(http.StatusBadRequest, utils.MapValidatorError[dto.GasPumpUpdatePathRequest](err))
		return
	}

	var body dto.GasPumpUpdateRequest

	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, utils.MapValidatorError[dto.GasPumpUpdateRequest](err))
		return
	}

	var gasPumpValues models.GasPump

	gasPumpValues.UpdatedByID = &user.ID

	copier.Copy(&gasPumpValues, body)

	id, _ := uuid.Parse(path.ID)
	updated, err := gp.repository.UpdateByID(id, &gasPumpValues)
	if err != nil {
		if utils.CheckDuplicatedEntry(err) {
			c.JSON(
				http.StatusConflict,
				dto.GeneralMessage{Detail: lang.DuplicatedEntry + "gas_station_id, number"},
			)
			return
		} else if utils.CheckMysqlErrCode(err, 1452) {
			c.JSON(http.StatusNotAcceptable, dto.GeneralMessage{Detail: lang.NotAcceptable + "gas_station_id"})
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

	if !updated {
		c.JSON(http.StatusNotFound, dto.GeneralMessage{Detail: lang.NotFoundRecord})
		return
	}

	c.JSON(http.StatusOK, dto.GeneralMessage{Detail: lang.RecordUpdated})
}

// @Summary Gas Pump Detail customer
// @Description Get Gas pump for customer when confirming
// @Tags Gas Pumps
// @Produce json
// @Router /api/v1/gas-pumps/{id}/customer [GET]
// @Param id path string true "uuid4 id" minLength(36) maxLength(36)
// @Param Authorization header string true "Token"
// @Success 200 {object} dto.GasPumpGetDetailForCustomerResponse "Gas station detail"
// @Failure 400 {array} dto.BadRequestMessage "Bad Request, failed on body, form, query..."
// @Failure 404 {object} dto.GeneralMessage "Not Found"
// @Failure 401 {object} dto.GeneralMessage "Unauthorized"
// @Failure 500 {object} dto.GeneralMessage "Internal server error"
func (gp *gasPumpController) GetDetailForCustomer(c *gin.Context) {
	var pathParams dto.GasPumpGetDetailForCustomerPathRequest
	if err := c.ShouldBindUri(&pathParams); err != nil {
		c.JSON(
			http.StatusBadRequest,
			utils.MapValidatorError[dto.GasPumpGetDetailForCustomerPathRequest](err),
		)
		return
	}

	customer := c.MustGet("customer").(*models.Customer)

	id, _ := uuid.Parse(pathParams.ID)
	station, err := gp.repository.GetActiveByID(id)
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

	var gasPump dto.GasPumpGetDetailForCustomerResponse

	applicablePromotionType, err := gp.settingsRepo.GetByName("applicable_promotion_type")
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
	// Check applicable promotion
	gasPump.DiscountType = "none"

	if applicablePromotionType != nil {
		promotionTypeInDB = applicablePromotionType.Value
	}

	if promotionTypeInDB == "campaign" {
		// Checking valid promotions
		campaign, err := gp.campaignRepository.GetApplicableCampaign(
			time.Now(),
			station.GasStation.ID,
		)
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

		if campaign != nil {
			gasPump.DiscountType = "campaign"
			gasPump.Campaign = &struct {
				Name     string  "json:\"name\""
				Discount float64 "json:\"discount\""
			}{
				Discount: *campaign.Discount,
				Name:     campaign.Name,
			}
		}
	} else if promotionTypeInDB == "elegibility" {
		gasPump.DiscountType = "elegibility"
	}

	copier.Copy(&gasPump, &station)

	c.JSON(http.StatusOK, gasPump)
}
