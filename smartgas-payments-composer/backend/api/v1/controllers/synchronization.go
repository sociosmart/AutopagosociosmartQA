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
)

type SynchronizationController interface {
	GetLastSync(*gin.Context)
	SyncNow(*gin.Context)
	List(*gin.Context)
	ListDetails(*gin.Context)
}

type synchronizationController struct {
	repository repository.SynchronizationRepository
	task       tasks.SynchronizationTask
}

func ProvideSynchronizationController(
	repository repository.SynchronizationRepository,
	task tasks.SynchronizationTask,
) *synchronizationController {
	return &synchronizationController{
		task:       task,
		repository: repository,
	}
}

// @Summary Last Synchronization
// @Description Get last synchronization
// @Tags Synchronization
// @Produce json
// @Accept json
// @Router /api/v1/synchronizations/last [GET]
// @Security Bearer
// @Param query query dto.SynchronizationGetLastSyncQueryRequest true "Sync type"
// @Success 200 {object} dto.SynchronizationGetLastSyncResponse "Access token & Refresh Token"
// @Failure 401 {object} dto.GeneralMessage "Unauthorized"
// @Failure 500 {object} dto.GeneralMessage "Internal server error"
// @Failure 400 {array} dto.BadRequestMessage "Bad Request, failed on body, form, query..."
func (sc *synchronizationController) GetLastSync(c *gin.Context) {
	var params dto.SynchronizationGetLastSyncQueryRequest

	if err := c.ShouldBind(&params); err != nil {
		c.JSON(
			http.StatusBadRequest,
			utils.MapValidatorError[dto.SynchronizationGetLastSyncQueryRequest](err),
		)
		return
	}

	last, err := sc.repository.GetLastByType(params.Type)

	if !utils.CheckNotFoundRecord(c, err) {
		return
	}

	var response dto.SynchronizationGetLastSyncResponse

	copier.Copy(&response, &last)

	c.JSON(http.StatusOK, response)
}

// @Summary Sync Now
// @Description Synchronize gas pumps, gas stations or all
// @Tags Synchronization
// @Produce json
// @Accept json
// @Router /api/v1/synchronizations/now [POST]
// @Security Bearer
// @Param params body dto.SynchronizationNowQueryRequest true "Sync type"
// @Success 200 {object} dto.GeneralMessage "Synchronized"
// @Failure 401 {object} dto.GeneralMessage "Unauthorized"
// @Failure 500 {object} dto.GeneralMessage "Internal server error"
// @Failure 423 {object} dto.GeneralMessage "Internal server error"
// @Failure 400 {array} dto.BadRequestMessage "Bad Request, failed on body, form, query..."
func (sc *synchronizationController) SyncNow(c *gin.Context) {
	var params dto.SynchronizationNowQueryRequest

	if err := c.ShouldBind(&params); err != nil {
		c.JSON(
			http.StatusBadRequest,
			utils.MapValidatorError[dto.SynchronizationNowQueryRequest](err),
		)
		return
	}

	user := c.MustGet("user").(*models.User)

	if params.Type == "gas_pumps" {
		if err := sc.task.SyncGasPumps(); err != nil {
			if errors.Is(err, tasks.RunningError) {
				c.JSON(http.StatusLocked, dto.GeneralMessage{Detail: err.Error()})
				return
			}
			// Logging error in sentry
			opts := &utils.TrackErrorOpts{
				Admin: user,
				Tags:  map[string]string{"auth_type": "admin"},
			}
			utils.TrackError(c, err, opts)
			c.JSON(
				http.StatusInternalServerError,
				dto.GeneralMessage{Detail: lang.InternalServerError},
			)
			return
		}
	} else if params.Type == "gas_stations" {
		if err := sc.task.SyncGasStations(); err != nil {
			if errors.Is(err, tasks.RunningError) {
				c.JSON(http.StatusLocked, dto.GeneralMessage{Detail: err.Error()})
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
	} else if params.Type == "customer_levels" {
		if err := sc.task.GenerateElegibilityCustomers(); err != nil {
			if errors.Is(err, tasks.RunningError) {
				c.JSON(http.StatusLocked, dto.GeneralMessage{Detail: err.Error()})
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
	}
	c.JSON(http.StatusOK, dto.GeneralMessage{Detail: lang.Synchronized})
}

// @Summary Synchronization List
// @Description Get paginated synchronizations
// @Tags Synchronization
// @Produce json
// @Router /api/v1/synchronizations [GET]
// @Security Bearer
// @Param param query dto.PaginateRequest true "Detalles"
// @Param param query dto.SynchronizationListQueryRequest true "Criterias"
// @Success 200 {object} dto.PaginationResponse{data=[]dto.SynchronizationListResponse} "Synchronizations List"
// @Failure 400 {array} dto.BadRequestMessage "Bad Request, failed on body, form, query..."
// @Failure 401 {object} dto.GeneralMessage "Unauthorized"
// @Failure 500 {object} dto.GeneralMessage "Internal server error"
func (sc *synchronizationController) List(c *gin.Context) {
	var pagination dto.PaginateRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, utils.MapValidatorError[dto.PaginateRequest](err))
		return
	}

	var params dto.SynchronizationListQueryRequest
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(
			http.StatusBadRequest,
			utils.MapValidatorError[dto.SynchronizationListQueryRequest](err),
		)
		return
	}

	var paginationSchema schemas.Pagination

	copier.Copy(&paginationSchema, &pagination)

	filters := make(map[string]any)

	if params.Type != "" {
		filters["type"] = params.Type
	}

	syncs, err := sc.repository.List(&paginationSchema, filters)

	user := c.MustGet("user").(*models.User)

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

	syncsResponse := make([]dto.SynchronizationListResponse, 0)

	copier.Copy(&syncsResponse, &syncs)

	var paginationResponse dto.PaginationResponse

	copier.Copy(&paginationResponse, &paginationSchema)

	paginationResponse.Data = syncsResponse

	c.JSON(http.StatusOK, paginationResponse)
}

// @Summary Synchronization Detail List
// @Description Get paginated synchronizations detail List
// @Tags Synchronization
// @Produce json
// @Router /api/v1/synchronizations/{id}/details [GET]
// @Security Bearer
// @Param param query dto.PaginateRequest true "Detalles"
// @Param id path string true "uuid4 id" minLength(36) maxLength(36)
// @Success 200 {object} dto.PaginationResponse{data=[]dto.SynchronizationListDetailResponse} "Synchronizations Detail List"
// @Failure 400 {array} dto.BadRequestMessage "Bad Request, failed on body, form, query..."
// @Failure 401 {object} dto.GeneralMessage "Unauthorized"
// @Failure 500 {object} dto.GeneralMessage "Internal server error"
func (sc *synchronizationController) ListDetails(c *gin.Context) {
	var pagination dto.PaginateRequest
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, utils.MapValidatorError[dto.PaginateRequest](err))
		return
	}

	var pathParams dto.SynchronizationListDetailPathRequest
	if err := c.ShouldBindUri(&pathParams); err != nil {
		c.JSON(
			http.StatusBadRequest,
			utils.MapValidatorError[dto.SynchronizationListDetailPathRequest](err),
		)
		return
	}

	var paginationSchema schemas.Pagination

	copier.Copy(&paginationSchema, &pagination)

	id, _ := uuid.Parse(pathParams.ID)

	details, err := sc.repository.ListDetails(
		&paginationSchema,
		map[string]any{"synchronization_id": id},
	)

	user := c.MustGet("user").(*models.User)

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

	syncsResponse := make([]dto.SynchronizationListDetailResponse, 0)

	copier.Copy(&syncsResponse, &details)

	var paginationResponse dto.PaginationResponse

	copier.Copy(&paginationResponse, &paginationSchema)

	paginationResponse.Data = syncsResponse

	c.JSON(http.StatusOK, paginationResponse)
}
