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

type PermissionController interface {
	ListAllGroups(*gin.Context)
	ListAll(*gin.Context)
}

type permissionController struct {
	repository repository.PermissionRepository
}

func ProvidePermissionController(repository repository.PermissionRepository) *permissionController {
	return &permissionController{
		repository: repository,
	}
}

// @Summary Get all permissions
// @Description Get all available permissions
// @Tags Permissions
// @Produce json
// @Router /api/v1/permissions/all [GET]
// @Security Bearer
// @Success 200 {array} dto.PermissionListAllResponse "Gas stations"
// @Failure 401 {object} dto.GeneralMessage "Unauthorized"
// @Failure 500 {object} dto.GeneralMessage "Internal server error"
func (pc *permissionController) ListAll(c *gin.Context) {
	permissions, err := pc.repository.ListAll()

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

	var permissionsResponse []dto.PermissionListAllResponse

	copier.Copy(&permissionsResponse, &permissions)

	c.JSON(http.StatusOK, permissionsResponse)
}

// @Summary Get all permission groups
// @Description Get all available permissions groups
// @Tags Permissions
// @Produce json
// @Router /api/v1/permissions/all-groups [GET]
// @Security Bearer
// @Success 200 {array} dto.GroupListAllResponse "Gas stations"
// @Failure 401 {object} dto.GeneralMessage "Unauthorized"
// @Failure 500 {object} dto.GeneralMessage "Internal server error"
func (pc *permissionController) ListAllGroups(c *gin.Context) {
	groups, err := pc.repository.ListAllGroups()

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

	var groupsResponse []dto.GroupListAllResponse

	copier.Copy(&groupsResponse, &groups)

	c.JSON(http.StatusOK, groupsResponse)
}
