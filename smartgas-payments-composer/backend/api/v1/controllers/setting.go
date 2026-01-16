package controllers

import (
	"net/http"
	"smartgas-payment/internal/dto"
	"smartgas-payment/internal/lang"
	"smartgas-payment/internal/models"
	"smartgas-payment/internal/repository"
	"smartgas-payment/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type SettingController interface {
	List(*gin.Context)
	CreateOrUpdate(*gin.Context)
}

type settingController struct {
	repo repository.SettingRepository
}

func ProvideSettingController(repo repository.SettingRepository) *settingController {
	return &settingController{
		repo: repo,
	}
}

// @Summary Get All settings
// @Description Endpoint to gather all settings
// @Tags Settings
// @Produce json
// @Accept json
// @Router /api/v1/settings [GET]
// @Security Bearer
// @Success 200 {array} dto.SettinGetAllResponse "All settings"
// @Failure 401 {object} dto.GeneralMessage "Unauthorized, token malformed, token invalid, expired, etc"
// @Failure 500 {object} dto.GeneralMessage "Internal Server Error"
func (sc *settingController) List(c *gin.Context) {
	settings, err := sc.repo.GetAll()

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

	settingsResponse := make([]dto.SettinGetAllResponse, 0)

	copier.Copy(&settingsResponse, settings)

	c.JSON(http.StatusOK, settingsResponse)
}

// @Summary Update setting
// @Description Update or create specific setting
// @Tags Settings
// @Router /api/v1/settings/ [POST]
// @Produce json
// @Security Bearer
// @Param data body dto.SettingCreateBody true "Setting Create or Update Body"
// @Success 200 {object} dto.GeneralMessage "Updated or created"
// @Failure 401 {object} dto.GeneralMessage "Unauthorized"
// @Failure 400 {array} dto.BadRequestMessage "Bad Request, failed on body, form, query..."
// @Failure 500 {object} dto.GeneralMessage "Internal server error"
func (sc *settingController) CreateOrUpdate(c *gin.Context) {
	var body dto.SettingCreateBody
	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, utils.MapValidatorError[dto.SettingCreateBody](err))
		return
	}

	user := c.MustGet("user").(*models.User)

	created, err := sc.repo.GetOrCreate(&models.Setting{Name: body.Name, Value: body.Value})

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

	if created {
		c.JSON(http.StatusOK, dto.GeneralMessage{Detail: "created"})
		return
	}

	_, err = sc.repo.Update(body.Name, body.Value)

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
