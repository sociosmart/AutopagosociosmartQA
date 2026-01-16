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

type UserController interface {
	Me(*gin.Context)
	List(*gin.Context)
	Create(*gin.Context)
	Update(*gin.Context)
	GetUserDetail(*gin.Context)
}

type userController struct {
	repository repository.UserRepository
}

func ProvideUserController(userRepository repository.UserRepository) *userController {
	return &userController{
		repository: userRepository,
	}
}

// @Summary User Me
// @Description Get user's information
// @Tags Users
// @Produce json
// @Router /api/v1/users/me [GET]
// @Security Bearer
// @Success 200 {object} dto.UserMeResponse "User's profile"
// @Failure 401 {object} dto.GeneralMessage "Unauthorized"
// @Failure 500 {object} dto.GeneralMessage "Internal server error"
func (uc *userController) Me(c *gin.Context) {
	user := c.MustGet("user").(*models.User)

	var me dto.UserMeResponse
	_ = copier.Copy(&me, user)

	c.JSON(http.StatusOK, &me)

}

// @Summary List users
// @Description Paginate users
// @Tags Users
// @Produce json
// @Router /api/v1/users [GET]
// @Security Bearer
// @Param param query dto.PaginateRequest true "Detalles"
// @Param search query string false "Lookup in users"
// @Success 200 {object} dto.PaginationResponse{data=[]dto.UserListResponse} "Users Paginated"
// @Failure 401 {object} dto.GeneralMessage "Unauthorized"
// @Failure 400 {array} dto.BadRequestMessage "Bad Request, failed on body, form, query..."
// @Failure 500 {object} dto.GeneralMessage "Internal server error"
func (uc *userController) List(c *gin.Context) {
	var pagination dto.PaginateRequest

	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, utils.MapValidatorError[dto.PaginateRequest](err))
		return
	}

	var paginationSchema schemas.Pagination

	copier.Copy(&paginationSchema, &pagination)

	search := c.Query("search")

	users, err := uc.repository.List(&paginationSchema, map[string]any{"search": "%" + search + "%"})

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

	usersResponse := make([]dto.UserListResponse, 0)

	copier.Copy(&usersResponse, &users)

	var paginationResponse dto.PaginationResponse

	copier.Copy(&paginationResponse, &paginationSchema)

	paginationResponse.Data = usersResponse

	c.JSON(http.StatusOK, paginationResponse)
}

// @Summary Create User
// @Description Create user with given parameters
// @Tags Users
// @Router /api/v1/users [POST]
// @Produce json
// @Security Bearer
// @Param data body dto.UserCreateRequest true "User body"
// @Success 201 {object} dto.GeneralMessage "Created"
// @Failure 401 {object} dto.GeneralMessage "Unauthorized"
// @Failure 409 {object} dto.GeneralMessage "Email already taken"
// @Failure 406 {object} dto.GeneralMessage "Foreign key for permission, group, station not exists"
// @Failure 400 {array} dto.BadRequestMessage "Bad Request, failed on body, form, query..."
// @Failure 500 {object} dto.GeneralMessage "Internal server error"
func (uc *userController) Create(c *gin.Context) {
	var body dto.UserCreateRequest

	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, utils.MapValidatorError[dto.UserCreateRequest](err))
		return
	}

	var user models.User

	copier.Copy(&user, &body)

	err := uc.repository.CreateUser(&user)

	user_ := c.MustGet("user").(*models.User)

	if err != nil {
		if utils.CheckDuplicatedEntry(err) {
			c.JSON(http.StatusConflict, dto.GeneralMessage{Detail: lang.DuplicatedEntry + "email"})
			return
		} else if utils.CheckMysqlErrCode(err, 1452) {
			c.JSON(http.StatusNotAcceptable, dto.GeneralMessage{Detail: lang.NotAcceptable})
			return
		}
		// Logging error in sentry
		opts := &utils.TrackErrorOpts{
			Admin: user_,
			Tags:  map[string]string{"auth_type": "admin"},
		}
		utils.TrackError(c, err, opts)
		c.JSON(http.StatusInternalServerError, dto.GeneralMessage{Detail: lang.InternalServerError})
		return
	}

	c.JSON(http.StatusCreated, dto.GeneralMessage{Detail: lang.RecordCreated})

}

// @Summary Update user
// @Description Update a user by id
// @Tags Users
// @Router /api/v1/users/{id} [PUT]
// @Produce json
// @Security Bearer
// @Param id path string true "uuid4 id" minLength(36) maxLength(36)
// @Param data body dto.UserUpdateRequest true "User body"
// @Success 200 {object} dto.GeneralMessage "Updated"
// @Failure 401 {object} dto.GeneralMessage "Unauthorized"
// @Failure 406 {object} dto.GeneralMessage "Foreign key for permission, group, station not exists"
// @Failure 400 {array} dto.BadRequestMessage "Bad Request, failed on body, form, query..."
// @Failure 500 {object} dto.GeneralMessage "Internal server error"
func (uc *userController) Update(c *gin.Context) {
	var body dto.UserUpdateRequest

	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, utils.MapValidatorError[dto.UserUpdateRequest](err))
		return
	}

	var pathParams dto.UserUpdatePathRequest

	if err := c.ShouldBindUri(&pathParams); err != nil {
		c.JSON(http.StatusBadRequest, utils.MapValidatorError[dto.UserUpdatePathRequest](err))
		return
	}

	var user models.User

	copier.Copy(&user, &body)

	id, _ := uuid.Parse(pathParams.ID)

	user_ := c.MustGet("user").(*models.User)

	if err := uc.repository.UpdateByID(id, &user); err != nil {
		if utils.CheckMysqlErrCode(err, 1452) {
			c.JSON(http.StatusNotAcceptable, dto.GeneralMessage{Detail: lang.NotAcceptable})
			return
		}
		// Logging error in sentry
		opts := &utils.TrackErrorOpts{
			Admin: user_,
			Tags:  map[string]string{"auth_type": "admin"},
		}
		utils.TrackError(c, err, opts)
		c.JSON(http.StatusInternalServerError, dto.GeneralMessage{Detail: lang.InternalServerError})
		return
	}

	c.JSON(http.StatusOK, dto.GeneralMessage{Detail: lang.RecordUpdated})
}

// @Summary Get user detail
// @Description Get user detail by id
// @Tags Users
// @Router /api/v1/users/{id} [GET]
// @Produce json
// @Security Bearer
// @Param id path string true "uuid4 id" minLength(36) maxLength(36)
// @Success 200 {object} dto.UserDetailResponse "User detail"
// @Failure 401 {object} dto.GeneralMessage "Unauthorized"
// @Failure 500 {object} dto.GeneralMessage "Internal server error"
func (uc *userController) GetUserDetail(c *gin.Context) {
	var pathParams dto.UserGetDetailPathRequest

	if err := c.ShouldBindUri(&pathParams); err != nil {
		c.JSON(http.StatusBadRequest, utils.MapValidatorError[dto.UserGetDetailPathRequest](err))
		return
	}
	id, _ := uuid.Parse(pathParams.ID)

	user, err := uc.repository.GetUserDetailByID(id)

	if !utils.CheckNotFoundRecord(c, err) {
		return
	}

	var userDetailResponse dto.UserDetailResponse

	copier.Copy(&userDetailResponse, &user)

	c.JSON(http.StatusOK, &userDetailResponse)
}
