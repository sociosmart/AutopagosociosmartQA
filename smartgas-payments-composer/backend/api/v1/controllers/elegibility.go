package controllers

import (
	"net/http"
	"smartgas-payment/internal/dto"
	"smartgas-payment/internal/lang"
	"smartgas-payment/internal/models"
	"smartgas-payment/internal/repository"
	"smartgas-payment/internal/schemas"
	"smartgas-payment/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type ElegibilityController interface {
	LevelList(*gin.Context)
	CreateLevel(*gin.Context)
	UpdateLevel(*gin.Context)
	CustomerLevelList(*gin.Context)
	LevelListAll(*gin.Context)
	UpdateCustomerLevel(*gin.Context)
	CreateCustomerLevel(*gin.Context)
}

type elegibilityController struct {
	repository repository.ElegibilityRepository
}

func ProvideElegibityController(
	repository repository.ElegibilityRepository,
) *elegibilityController {
	return &elegibilityController{
		repository: repository,
	}
}

// @Summary List Levels
// @Description Paginate Levels
// @Tags Elegibility
// @Produce json
// @Router /api/v1/elegibility/levels [GET]
// @Security Bearer
// @Param param query dto.PaginateRequest true "Detalles"
// @Param search query string false "Lookup in levels"
// @Success 200 {object} dto.PaginationResponse{data=[]dto.LevelListResponse} "Levels Paginated"
// @Failure 401 {object} dto.GeneralMessage "Unauthorized"
// @Failure 400 {array} dto.BadRequestMessage "Bad Request, failed on body, form, query..."
// @Failure 500 {object} dto.GeneralMessage "Internal server error"
func (ec *elegibilityController) LevelList(c *gin.Context) {
	var pagination dto.PaginateRequest

	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, utils.MapValidatorError[dto.PaginateRequest](err))
		return
	}

	var paginationSchema schemas.Pagination

	copier.Copy(&paginationSchema, &pagination)

	search := c.Query("search")

	levels, err := ec.repository.LevelList(
		&paginationSchema,
		map[string]any{"search": "%" + search + "%"},
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

	levelsResponse := make([]dto.LevelListResponse, 0)

	copier.Copy(&levelsResponse, &levels)

	var paginationResponse dto.PaginationResponse

	copier.Copy(&paginationResponse, &paginationSchema)

	paginationResponse.Data = levelsResponse

	c.JSON(http.StatusOK, paginationResponse)
}

// @Summary Create Elegibility Level
// @Description Create Elegibility level with given parameters
// @Tags Elegibility
// @Router /api/v1/elegibility/levels [POST]
// @Produce json
// @Security Bearer
// @Param data body dto.EelegibilityLevelCreateRequest true "Elegibility body"
// @Success 201 {object} dto.LevelCreateResponse "Created"
// @Failure 401 {object} dto.GeneralMessage "Unauthorized"
// @Failure 409 {object} dto.GeneralMessage "Duplicated entry"
// @Failure 400 {array} dto.BadRequestMessage "Bad Request, failed on body, form, query..."
// @Failure 500 {object} dto.GeneralMessage "Internal server error"
func (ec *elegibilityController) CreateLevel(c *gin.Context) {
	var body dto.EelegibilityLevelCreateRequest

	if err := c.ShouldBind(&body); err != nil {
		c.JSON(
			http.StatusBadRequest,
			utils.MapValidatorError[dto.EelegibilityLevelCreateRequest](err),
		)
		return
	}

	user := c.MustGet("user").(*models.User)

	var level models.Level

	copier.Copy(&level, &body)

	level.CreatedByID = &user.ID

	err := ec.repository.CreateLevel(&level)
	if err != nil {
		if utils.CheckDuplicatedEntry(err) {
			c.JSON(
				http.StatusConflict,
				dto.GeneralMessage{Detail: lang.DuplicatedEntry + "name"},
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

	c.JSON(http.StatusCreated, dto.LevelCreateResponse{ID: level.ID.String()})
}

// @Summary Update Elegibility Level
// @Description Update Elegibility Level with given parameters
// @Tags Elegibility
// @Router /api/v1/elegibility/levels/{id} [PUT]
// @Produce json
// @Security Bearer
// @Param data body dto.ElegibilityLevelUpdatesRequest true "Campaign body"
// @Param id path string true "uuid4 id" minLength(36) maxLength(36)
// @Success 200 {object} dto.GeneralMessage "Updated"
// @Failure 401 {object} dto.GeneralMessage "Unauthorized"
// @Failure 409 {object} dto.GeneralMessage "Duplicated entry"
// @Failure 404 {object} dto.GeneralMessage "Not found"
// @Failure 400 {array} dto.BadRequestMessage "Bad Request, failed on body, form, query..."
// @Failure 500 {object} dto.GeneralMessage "Internal server error"
func (ec *elegibilityController) UpdateLevel(c *gin.Context) {
	var body dto.ElegibilityLevelUpdatesRequest

	if err := c.ShouldBind(&body); err != nil {
		c.JSON(
			http.StatusBadRequest,
			utils.MapValidatorError[dto.ElegibilityLevelUpdatesRequest](err),
		)
		return
	}

	var pathParams dto.ElegibilityLevelUpdatePathRequest

	if err := c.ShouldBindUri(&pathParams); err != nil {
		c.JSON(
			http.StatusBadRequest,
			utils.MapValidatorError[dto.ElegibilityLevelUpdatePathRequest](err),
		)
		return
	}

	id, _ := uuid.Parse(pathParams.ID)

	var level models.Level

	copier.Copy(&level, &body)

	user := c.MustGet("user").(*models.User)

	level.UpdatedByID = &user.ID

	updated, err := ec.repository.UpdateLevelByID(id, &level)
	if err != nil {
		if utils.CheckDuplicatedEntry(err) {
			c.JSON(
				http.StatusConflict,
				dto.GeneralMessage{Detail: lang.DuplicatedEntry + "name"},
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

	c.JSON(http.StatusOK, dto.GeneralMessage{Detail: "updated"})
}

// @Summary Customer Levels List
// @Description Paginate Customer Levels
// @Tags Elegibility
// @Produce json
// @Router /api/v1/elegibility/customers/levels [GET]
// @Security Bearer
// @Param param query dto.PaginateRequest true "Detalles"
// @Param search query string false "Lookup in customer levels"
// @Success 200 {object} dto.PaginationResponse{data=[]dto.CustomerLevelListResponse} "Customer Levels Paginated"
// @Failure 401 {object} dto.GeneralMessage "Unauthorized"
// @Failure 400 {array} dto.BadRequestMessage "Bad Request, failed on body, form, query..."
// @Failure 500 {object} dto.GeneralMessage "Internal server error"
func (ec *elegibilityController) CustomerLevelList(c *gin.Context) {
	var pagination dto.PaginateRequest

	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, utils.MapValidatorError[dto.PaginateRequest](err))
		return
	}

	var paginationSchema schemas.Pagination

	copier.Copy(&paginationSchema, &pagination)

	search := c.Query("search")

	customerLevels, err := ec.repository.CustomerLevelList(
		&paginationSchema,
		map[string]any{"search": "%" + search + "%"},
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

	customerLevelsResponse := make([]dto.CustomerLevelListResponse, 0)

	copier.Copy(&customerLevelsResponse, &customerLevels)

	var paginationResponse dto.PaginationResponse

	copier.Copy(&paginationResponse, &paginationSchema)

	paginationResponse.Data = customerLevelsResponse

	c.JSON(http.StatusOK, paginationResponse)
}

// @Summary List all levels
// @Description List of all levels
// @Tags Elegibility
// @Produce json
// @Router /api/v1/elegibility/levels/all [GET]
// @Security Bearer
// @Success 200 {array} dto.LevelListAllResponse "Levels"
// @Failure 401 {object} dto.GeneralMessage "Unauthorized"
// @Failure 500 {object} dto.GeneralMessage "Internal server error"
func (ec *elegibilityController) LevelListAll(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	levels, err := ec.repository.LevelListAll()
	if err != nil {
		// Logging error for when swit list cards fails
		opts := &utils.TrackErrorOpts{
			Admin: user,
			Tags:  map[string]string{"auth_type": "customer"},
		}
		utils.TrackError(c, err, opts)
		c.JSON(http.StatusInternalServerError, dto.GeneralMessage{Detail: lang.InternalServerError})
		return
	}

	levelsResponse := make([]dto.LevelListAllResponse, 0)

	copier.Copy(&levelsResponse, levels)

	c.JSON(http.StatusOK, levelsResponse)
}

// @Summary Update customer level
// @Description Update customer level with given parameters
// @Tags Elegibility
// @Router /api/v1/elegibility/customers/levels/{id} [PUT]
// @Produce json
// @Security Bearer
// @Param data body dto.ElegibilityUpdateCustomerLevelRequest true "Customer Level Update body"
// @Param id path string true "uuid4 id" minLength(36) maxLength(36)
// @Success 200 {object} dto.GeneralMessage "Updated"
// @Failure 401 {object} dto.GeneralMessage "Unauthorized"
// @Failure 406 {object} dto.GeneralMessage "Foreign key for level, customer not exists"
// @Failure 404 {object} dto.GeneralMessage "Not found"
// @Failure 409 {object} dto.GeneralMessage "Duplicated entry"
// @Failure 400 {array} dto.BadRequestMessage "Bad Request, failed on body, form, query..."
// @Failure 500 {object} dto.GeneralMessage "Internal server error"
func (ec *elegibilityController) UpdateCustomerLevel(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	var path dto.ElegibilityCustomerLevelUpdatePathRequest

	if err := c.ShouldBindUri(&path); err != nil {
		c.JSON(
			http.StatusBadRequest,
			utils.MapValidatorError[dto.ElegibilityCustomerLevelUpdatePathRequest](err),
		)
		return
	}

	var body dto.ElegibilityUpdateCustomerLevelRequest

	if err := c.ShouldBind(&body); err != nil {
		c.JSON(
			http.StatusBadRequest,
			utils.MapValidatorError[dto.ElegibilityUpdateCustomerLevelRequest](err),
		)
		return
	}

	id, _ := uuid.Parse(path.ID)

	var cusLevel models.CustomerLevel

	copier.Copy(&cusLevel, &body)

	cusLevel.UpdatedByID = &user.ID
	cusLevel.ManuallyTouched = utils.BoolAddr(true)

	updated, err := ec.repository.UpdateCustomerLevelByID(id, &cusLevel)
	if err != nil {
		if utils.CheckDuplicatedEntry(err) {
			c.JSON(
				http.StatusConflict,
				dto.GeneralMessage{
					Detail: lang.DuplicatedEntry + "customer_id, validity_month, validity_year",
				},
			)
			return
		} else if utils.CheckMysqlErrCode(err, 1452) {
			c.JSON(http.StatusNotAcceptable, dto.GeneralMessage{Detail: lang.NotAcceptable + "elegibility_level_id or customer_id"})
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

// @Summary Create customer level
// @Description Create customer level with given parameters
// @Tags Elegibility
// @Router /api/v1/elegibility/customers/levels/ [POST]
// @Produce json
// @Security Bearer
// @Param data body dto.ElegibilityCreateCustomerLevelRequest true "Customer Level create body"
// @Success 201 {object} dto.CustomerLevelCreateResponse "Created"
// @Failure 401 {object} dto.GeneralMessage "Unauthorized"
// @Failure 406 {object} dto.GeneralMessage "Foreign key for level, customer not exists"
// @Failure 409 {object} dto.GeneralMessage "Duplicated entry"
// @Failure 400 {array} dto.BadRequestMessage "Bad Request, failed on body, form, query..."
// @Failure 500 {object} dto.GeneralMessage "Internal server error"
func (ec *elegibilityController) CreateCustomerLevel(c *gin.Context) {
	var body dto.ElegibilityCreateCustomerLevelRequest

	if err := c.ShouldBind(&body); err != nil {
		c.JSON(
			http.StatusBadRequest,
			utils.MapValidatorError[dto.ElegibilityCreateCustomerLevelRequest](err),
		)
		return
	}

	user := c.MustGet("user").(*models.User)

	var cusLevel models.CustomerLevel

	copier.Copy(&cusLevel, &body)

	cusLevel.CreatedByID = &user.ID
	cusLevel.ManuallyTouched = utils.BoolAddr(true)

	err := ec.repository.CreateCustomerLevel(&cusLevel)
	if err != nil {
		if utils.CheckDuplicatedEntry(err) {
			c.JSON(
				http.StatusConflict,
				dto.GeneralMessage{
					Detail: lang.DuplicatedEntry + "customer_id, validity_month, validity_year",
				},
			)
			return
		} else if utils.CheckMysqlErrCode(err, 1452) {
			c.JSON(http.StatusNotAcceptable, dto.GeneralMessage{Detail: lang.NotAcceptable + "elegibility_level_id or customer_id"})
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

	c.JSON(http.StatusCreated, dto.LevelCreateResponse{ID: cusLevel.ID.String()})
}
