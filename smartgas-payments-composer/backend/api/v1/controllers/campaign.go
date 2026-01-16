package controllers

import (
	"errors"
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
	"gorm.io/gorm"
)

type CampaignController interface {
	List(*gin.Context)
	Create(*gin.Context)
	Update(*gin.Context)
	GetCampaignDetail(*gin.Context)
}

type campaignController struct {
	repository repository.CampaignRepository
}

func ProvideCampaignController(repository repository.CampaignRepository) *campaignController {
	return &campaignController{
		repository: repository,
	}
}

// @Summary Campaign List
// @Description Get Campaigns with pagination
// @Tags Campaigns
// @Produce json
// @Router /api/v1/campaigns [GET]
// @Security Bearer
// @Param param query dto.PaginateRequest true "Detalles"
// @Param search query string false "Search promotion"
// @Success 200 {object} dto.PaginationResponse{data=[]dto.CampaignListResponse} "Paginated Campaigns"
// @Failure 400 {array} dto.BadRequestMessage "Bad Request, failed on body, form, query..."
// @Failure 401 {object} dto.GeneralMessage "Unauthorized"
// @Failure 500 {object} dto.GeneralMessage "Internal server error"
func (pc *campaignController) List(c *gin.Context) {
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

	campaigns, err := pc.repository.List(&paginationSchema, filters)
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

	campaignsResponse := make([]dto.CampaignListResponse, 0)

	copier.Copy(&campaignsResponse, &campaigns)

	var paginationResponse dto.PaginationResponse

	copier.Copy(&paginationResponse, &paginationSchema)

	paginationResponse.Data = campaignsResponse

	c.JSON(http.StatusOK, paginationResponse)
}

// @Summary Create Campaign
// @Description Create campaign with given parameters
// @Tags Campaigns
// @Router /api/v1/campaigns [POST]
// @Produce json
// @Security Bearer
// @Param data body dto.CampaignCreateRequest true "Campaign body"
// @Success 201 {object} dto.GeneralMessage "Created"
// @Failure 401 {object} dto.GeneralMessage "Unauthorized"
// @Failure 406 {object} dto.GeneralMessage "Foreign key for permission, group, station not exists"
// @Failure 400 {array} dto.BadRequestMessage "Bad Request, failed on body, form, query..."
// @Failure 500 {object} dto.GeneralMessage "Internal server error"
func (cc *campaignController) Create(c *gin.Context) {
	var body dto.CampaignCreateRequest

	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, utils.MapValidatorError[dto.CampaignCreateRequest](err))
		return
	}

	user := c.MustGet("user").(*models.User)

	var campaign models.Campaign

	copier.Copy(&campaign, &body)

	campaign.CreatedBy = user

	err := cc.repository.Create(&campaign)
	if err != nil {
		if utils.CheckDuplicatedEntry(err) {
			c.JSON(
				http.StatusConflict,
				dto.GeneralMessage{Detail: lang.DuplicatedEntry + "gas_station"},
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

	c.JSON(http.StatusCreated, dto.GeneralMessage{Detail: "created"})
}

// @Summary Update Campaign
// @Description Update campaign with given parameters
// @Tags Campaigns
// @Router /api/v1/campaigns/{id} [PUT]
// @Produce json
// @Security Bearer
// @Param data body dto.CampaignUpdateRequest true "Campaign body"
// @Param id path string true "uuid4 id" minLength(36) maxLength(36)
// @Success 200 {object} dto.GeneralMessage "Updated"
// @Failure 401 {object} dto.GeneralMessage "Unauthorized"
// @Failure 406 {object} dto.GeneralMessage "Foreign key for permission, group, station not exists"
// @Failure 400 {array} dto.BadRequestMessage "Bad Request, failed on body, form, query..."
// @Failure 500 {object} dto.GeneralMessage "Internal server error"
func (cc *campaignController) Update(c *gin.Context) {
	var body dto.CampaignUpdateRequest

	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, utils.MapValidatorError[dto.CampaignUpdateRequest](err))
		return
	}

	var pathParams dto.CampaignUpdatePathRequest

	if err := c.ShouldBindUri(&pathParams); err != nil {
		c.JSON(http.StatusBadRequest, utils.MapValidatorError[dto.CampaignUpdatePathRequest](err))
		return
	}

	user := c.MustGet("user").(*models.User)

	var campaign models.Campaign

	copier.Copy(&campaign, &body)

	campaign.UpdatedBy = user

	id, _ := uuid.Parse(pathParams.ID)

	err := cc.repository.UpdateByID(id, &campaign)
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

	c.JSON(http.StatusOK, dto.GeneralMessage{Detail: "updated"})
}

// @Summary Get Campaign detail
// @Description Get campaign detail
// @Tags Campaigns
// @Router /api/v1/campaigns/{id} [GET]
// @Produce json
// @Security Bearer
// @Param id path string true "uuid4 id" minLength(36) maxLength(36)
// @Success 200 {object} dto.CampaignDetailResponse "Campaign Detail"
// @Failure 401 {object} dto.GeneralMessage "Unauthorized"
// @Failure 404 {object} dto.GeneralMessage "Not found"
// @Failure 400 {array} dto.BadRequestMessage "Bad Request, failed on body, form, query..."
// @Failure 500 {object} dto.GeneralMessage "Internal server error"
func (cc *campaignController) GetCampaignDetail(c *gin.Context) {
	var pathParams dto.CampaignDetailPathRequest

	if err := c.ShouldBindUri(&pathParams); err != nil {
		c.JSON(http.StatusBadRequest, utils.MapValidatorError[dto.CampaignDetailPathRequest](err))
		return
	}

	user := c.MustGet("user").(*models.User)

	id, _ := uuid.Parse(pathParams.ID)

	campaign, err := cc.repository.GetCampaignByID(id, map[string]any{})
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
	}

	var campaignResponse dto.CampaignDetailResponse

	copier.Copy(&campaignResponse, campaign)

	c.JSON(http.StatusOK, &campaignResponse)
}

