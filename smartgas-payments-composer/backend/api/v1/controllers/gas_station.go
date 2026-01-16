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

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type GasStationController interface {
	List(*gin.Context)
	Get(*gin.Context)
	Update(*gin.Context)
	Create(*gin.Context)
	ListAll(*gin.Context)
}

type gasStationController struct {
	repository          repository.GasStationRepository
	synchronizationTask tasks.SynchronizationTask
}

func ProvideGasStationProvider(
	repository repository.GasStationRepository,
	synchronizationTask tasks.SynchronizationTask,
) *gasStationController {
	return &gasStationController{
		repository:          repository,
		synchronizationTask: synchronizationTask,
	}
}

// @Summary Gas Station List
// @Description Get paginated gas stations
// @Tags Gas Stations
// @Produce json
// @Router /api/v1/gas-stations [GET]
// @Security Bearer
// @Param param query dto.PaginateRequest true "Detalles"
// @Param search query string false "Search in gas stations, it looks up through id, name ip"
// @Success 200 {object} dto.PaginationResponse{data=[]dto.GasStationListResponse} "Gas station paginated"
// @Failure 400 {array} dto.BadRequestMessage "Bad Request, failed on body, form, query..."
// @Failure 401 {object} dto.GeneralMessage "Unauthorized"
// @Failure 500 {object} dto.GeneralMessage "Internal server error"
func (gs *gasStationController) List(c *gin.Context) {
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

	stations, err := gs.repository.List(&paginationSchema, filters)
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

	stationsResponse := make([]dto.GasStationListResponse, 0)

	copier.Copy(&stationsResponse, &stations)

	// Making the response
	var paginationResponse dto.PaginationResponse

	copier.Copy(&paginationResponse, &paginationSchema)

	paginationResponse.Data = stationsResponse

	c.JSON(http.StatusOK, paginationResponse)
}

// @Summary Gas Station List All
// @Description Get all gas stations
// @Tags Gas Stations
// @Produce json
// @Router /api/v1/gas-stations/all [GET]
// @Security Bearer
// @Success 200 {array} dto.GasStationListAllResponse "Gas stations"
// @Failure 401 {object} dto.GeneralMessage "Unauthorized"
// @Failure 500 {object} dto.GeneralMessage "Internal server error"
func (gs *gasStationController) ListAll(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	filters := make(map[string]any)

	utils.AddStationsFilter(user, filters)

	stations, err := gs.repository.ListAll(filters)
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

	stationsResponse := make([]dto.GasStationListAllResponse, 0)

	copier.Copy(&stationsResponse, stations)

	c.JSON(http.StatusOK, stationsResponse)
}

// @Summary Gas Station Create
// @Description Create Gas Station
// @Tags Gas Stations
// @Produce json
// @Router /api/v1/gas-stations [POST]
// @Security Bearer
// @Param Data body dto.GasStationCreateRequest true "Gas Station Body"
// @Success 201 {object} dto.GasStationCreateResponse "Gas station create response"
// @Failure 401 {object} dto.GeneralMessage "Unauthorized"
// @Failure 409 {object} dto.GeneralMessage "Duplicated entry"
// @Failure 400 {array} dto.BadRequestMessage "Bad Request, failed on body, form, query..."
// @Failure 500 {object} dto.GeneralMessage "Internal server error"
func (gs *gasStationController) Create(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	var body dto.GasStationCreateRequest
	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, utils.MapValidatorError[dto.GasStationCreateRequest](err))
		return
	}

	var gasStation models.GasStation
	gasStation.CreatedBy = user

	copier.Copy(&gasStation, &body)

	if err := gs.repository.Create(&gasStation); err != nil {
		if utils.CheckDuplicatedEntry(err) {
			c.JSON(
				http.StatusConflict,
				dto.GeneralMessage{Detail: lang.DuplicatedEntry + "name, ip"},
			)
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

	var response dto.GasStationCreateResponse

	copier.Copy(&response, &gasStation)

	c.JSON(http.StatusCreated, response)
}

// @Summary Gas Station Detail
// @Description Get Gas station detail
// @Tags Gas Stations
// @Produce json
// @Router /api/v1/gas-stations/{id} [GET]
// @Security Bearer
// @Param id path string true "uuid4 id" minLength(36) maxLength(36)
// @Success 200 {object} dto.GasStationGetResponse "Gas station detail"
// @Failure 400 {array} dto.BadRequestMessage "Bad Request, failed on body, form, query..."
// @Failure 404 {object} dto.GeneralMessage "Not Found"
// @Failure 401 {object} dto.GeneralMessage "Unauthorized"
// @Failure 500 {object} dto.GeneralMessage "Internal server error"
func (gs *gasStationController) Get(c *gin.Context) {
	var pathParams dto.GasStationGetPathRequest
	if err := c.ShouldBindUri(&pathParams); err != nil {
		c.JSON(http.StatusBadRequest, utils.MapValidatorError[dto.GasStationGetPathRequest](err))
		return
	}

	id, _ := uuid.Parse(pathParams.ID)
	station, err := gs.repository.GetByID(id)

	user := c.MustGet("user").(*models.User)

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

	var gasStation dto.GasStationGetResponse

	copier.Copy(&gasStation, &station)

	c.JSON(http.StatusOK, gasStation)
}

// @Summary Gas Station Update
// @Description Update Gas station
// @Tags Gas Stations
// @Produce json
// @Router /api/v1/gas-stations/{id} [PUT]
// @Security Bearer
// @Param id path string true "uuid4 id" minLength(36) maxLength(36)
// @Param Data body dto.GasStationUpdateRequest true "Gas station values"
// @Success 200 {object} dto.GeneralMessage "Record Updated"
// @Failure 400 {array} dto.BadRequestMessage "Bad Request, failed on body, form, query..."
// @Failure 404 {object} dto.GeneralMessage "Not Found"
// @Failure 401 {object} dto.GeneralMessage "Unauthorized"
// @Failure 409 {object} dto.GeneralMessage "Duplicated entry"
// @Failure 500 {object} dto.GeneralMessage "Internal server error"
func (gs *gasStationController) Update(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	var path dto.GasStationUpdatePathRequest

	if err := c.ShouldBindUri(&path); err != nil {
		c.JSON(http.StatusBadRequest, utils.MapValidatorError[dto.GasStationUpdatePathRequest](err))
		return
	}

	var body dto.GasStationUpdateRequest

	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, utils.MapValidatorError[dto.GasStationUpdateRequest](err))
		return
	}

	var gasStationValues models.GasStation
	gasStationValues.UpdatedByID = &user.ID

	copier.Copy(&gasStationValues, body)

	id, _ := uuid.Parse(path.ID)
	updated, err := gs.repository.UpdateByID(id, &gasStationValues)
	if err != nil {
		if utils.CheckDuplicatedEntry(err) {
			c.JSON(
				http.StatusConflict,
				dto.GeneralMessage{Detail: lang.DuplicatedEntry + "name, ip"},
			)
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
