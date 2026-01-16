package controllers_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"smartgas-payment/internal/dto"
	"smartgas-payment/internal/injectors"
	"smartgas-payment/internal/lang"
	"smartgas-payment/internal/models"
	"smartgas-payment/internal/repository"
	"smartgas-payment/internal/schemas"
	"smartgas-payment/internal/utils"
	"strconv"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type gasPumpCtrlTest struct {
	suite.Suite
	repository          *repository.MockGasPumpRepository
	userRepository      *repository.MockUserRepository
	testRequest         *utils.TestRequest
	validate            *validator.Validate
	userID              uuid.UUID
	gasPumpID           uuid.UUID
	gasPump             *models.GasPump
	validToken          string
	invalidToken        string
	gasPumpErrorID      uuid.UUID
	gasPumpNotFoundID   uuid.UUID
	gasPumpConflictID   uuid.UUID
	gasPumpDuplicatedID uuid.UUID
	gasPumps            []*models.GasPump
}

func (suite *gasPumpCtrlTest) SetupSuite() {
	setup, _ := injectors.InitializeServerWithMocks()

	suite.validate = validator.New()

	suite.repository = setup.GasPumpRepositoryMock
	suite.userRepository = setup.UserRepositoryMock

	suite.testRequest = &utils.TestRequest{
		Router: setup.Router,
	}

	suite.userID = uuid.New()

	user := &models.User{
		ID:      suite.userID,
		IsAdmin: utils.BoolAddr(true),
	}

	suite.userRepository.On("GetUserByID", suite.userID).Return(user, nil)
	suite.userRepository.On("GetUserByID", mock.AnythingOfType("uuid.UUID")).Return(nil, gorm.ErrRecordNotFound)

	claims := &schemas.JwtClaims{
		Sub: suite.userID,
	}

	token, _ := claims.ClaimToken()

	suite.validToken = token

	// Invalid token
	claims.Sub = uuid.New()
	token, _ = claims.ClaimToken()
	suite.invalidToken = token

	// Get by ID setup
	suite.gasPumpID = uuid.New()
	gasStationID := uuid.New()
	suite.gasPumpErrorID = uuid.New()
	suite.gasPumpNotFoundID = uuid.New()

	suite.gasPump = &models.GasPump{
		ID:           suite.gasPumpID,
		Number:       "01",
		GasStationID: &gasStationID,
		GasStation: &models.GasStation{
			ID:   uuid.New(),
			Ip:   "192.168.100.100",
			Name: "Guerrero",
		},
	}
	suite.repository.On("GetByID", suite.gasPumpID).Return(suite.gasPump, nil)
	suite.repository.On("GetByID", suite.gasPumpErrorID).Return(nil, errors.New(lang.InternalServerError))
	suite.repository.On("GetByID", suite.gasPumpNotFoundID).Return(nil, gorm.ErrRecordNotFound)

	// Create
	mysqlErrConflict := &mysql.MySQLError{
		Number: 1452,
	}

	mysqlErrDuplicated := &mysql.MySQLError{
		Number: 1062,
	}

	suite.gasPumpConflictID = uuid.New()
	suite.gasPumpDuplicatedID = uuid.New()

	suite.repository.On("Create", mock.MatchedBy(func(pump *models.GasPump) bool {
		if *pump.GasStationID == suite.gasPumpID {
			pump.ID = suite.gasPumpID
			return true
		}
		return false
	})).Return(nil)

	suite.repository.On("Create", mock.MatchedBy(func(pump *models.GasPump) bool {
		if *pump.GasStationID == suite.gasPumpConflictID {
			return true
		}
		return false
	})).Return(mysqlErrConflict)

	suite.repository.On("Create", mock.MatchedBy(func(pump *models.GasPump) bool {
		if *pump.GasStationID == suite.gasPumpDuplicatedID {
			return true
		}
		return false
	})).Return(mysqlErrDuplicated)

	suite.repository.On("Create", mock.MatchedBy(func(pump *models.GasPump) bool {
		if *pump.GasStationID == suite.gasPumpErrorID {
			return true
		}
		return false
	})).Return(errors.New(lang.InternalServerError))

	// Test list
	suite.gasPumps = make([]*models.GasPump, 0)

	// map gas stations
	for i := 1; i <= 22; i++ {
		gasId := uuid.New()
		suite.gasPumps = append(suite.gasPumps, &models.GasPump{
			ID:           uuid.New(),
			Number:       fmt.Sprintf("%02d", i),
			ExternalID:   "01",
			GasStationID: &gasId,
			GasStation: &models.GasStation{
				ID:   gasId,
				Name: "Guerrero " + strconv.Itoa(i),
			},
		})
	}

	suite.repository.On("List", mock.MatchedBy(func(pagination *schemas.Pagination) bool {
		if pagination.Limit == 10 && pagination.Page == 1 {
			pagination.TotalPages = 3
			pagination.TotalRows = int64(len(suite.gasPumps))
			return true
		}

		return false
	}), mock.Anything).Return(suite.gasPumps[:10], nil).Once()

	suite.repository.On("List", mock.MatchedBy(func(pagination *schemas.Pagination) bool {
		if pagination.Limit == 10 && pagination.Page == 3 {
			pagination.TotalPages = 3
			pagination.TotalRows = int64(len(suite.gasPumps))
			return true
		}
		return false
	}), mock.Anything).Return(suite.gasPumps[20:21], nil).Once()

	suite.repository.On("List", mock.MatchedBy(func(pagination *schemas.Pagination) bool {
		if pagination.Limit >= 100 && pagination.Page == 1 {
			pagination.TotalPages = 1
			pagination.TotalRows = int64(len(suite.gasPumps))
			return true
		}
		return false
	}), mock.Anything).Return(suite.gasPumps[:len(suite.gasPumps)], nil)

	suite.repository.On("List", mock.MatchedBy(func(pagination *schemas.Pagination) bool {
		if pagination.Limit == 10 && pagination.Page == 10 {
			return true
		}
		return false
	}), mock.Anything).Return([]*models.GasPump{}, nil).Once()

	// Simulate internal server error
	suite.repository.On("List", mock.MatchedBy(func(pagination *schemas.Pagination) bool {
		if pagination.Limit == 100 && pagination.Page == 1000 {
			return true
		}
		return false
	}), mock.Anything).Return([]*models.GasPump{}, errors.New(lang.InternalServerError))

	// Update

	suite.repository.On("UpdateByID", suite.gasPumpID, mock.AnythingOfType("*models.GasPump")).Return(true, nil)
	suite.repository.On("UpdateByID", suite.gasPumpNotFoundID, mock.AnythingOfType("*models.GasPump")).Return(false, nil)
	suite.repository.On("UpdateByID", suite.gasPumpDuplicatedID, mock.AnythingOfType("*models.GasPump")).Return(false, mysqlErrDuplicated)
	suite.repository.On("UpdateByID", suite.gasPumpConflictID, mock.AnythingOfType("*models.GasPump")).Return(false, mysqlErrConflict)
	suite.repository.On("UpdateByID", suite.gasPumpErrorID, mock.AnythingOfType("*models.GasPump")).Return(false, errors.New(lang.InternalServerError))

}

func (suite *gasPumpCtrlTest) TestUpdate() {
	url := "/api/v1/gas-pumps/"

	wrongformat := "wrongformat"
	err := suite.validate.Struct(dto.GasPumpUpdateRequest{
		Number:       utils.StringAddr(""),
		GasStationID: &wrongformat,
	})

	expectedBadResponse := utils.MapValidatorError[dto.GasPumpUpdateRequest](err)

	err = suite.validate.Struct(dto.GasPumpUpdatePathRequest{ID: wrongformat})

	expectedPathBadRequest := utils.MapValidatorError[dto.GasPumpUpdatePathRequest](err)

	fakeId := uuid.NewString()

	testcases := []struct {
		Name               string
		Url                string
		Body               any
		ExpectedStatusCode int
		ExpectedResponse   any
		Pre                func()
		Pos                func()
	}{
		{
			Name:               "TestGasPumpController_UpdateUnauthorized",
			Url:                url + suite.gasPumpID.String(),
			ExpectedStatusCode: http.StatusUnauthorized,
			ExpectedResponse:   dto.GeneralMessage{Detail: lang.NotFoundRecord},
			Pre: func() {
				suite.testRequest.SetBearerToken("Bearer " + suite.invalidToken)
			},
			Pos: func() {
				suite.testRequest.SetBearerToken("")
			},
		},
		{
			Name:               "TestGasPumpController_UpdatePathBadRequest",
			Url:                url + wrongformat,
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedResponse:   expectedPathBadRequest,
			Pre: func() {
				suite.testRequest.SetBearerToken("Bearer " + suite.validToken)
			},
			Pos: func() {
				suite.testRequest.SetBearerToken("")
			},
		},
		{
			Name:               "TestGasPumpController_UpdateBodyBadRequest",
			Url:                url + suite.gasPumpID.String(),
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedResponse:   expectedBadResponse,
			Body: dto.GasPumpUpdateRequest{
				GasStationID: &wrongformat,
				Number:       utils.StringAddr("1"),
			},
			Pre: func() {
				suite.testRequest.SetBearerToken("Bearer " + suite.validToken)
			},
			Pos: func() {
				suite.testRequest.SetBearerToken("")
			},
		},
		{
			Name:               "TestGasPumpController_UpdateInternalServerError",
			Url:                url + suite.gasPumpErrorID.String(),
			ExpectedStatusCode: http.StatusInternalServerError,
			ExpectedResponse:   dto.GeneralMessage{Detail: lang.InternalServerError},
			Body:               dto.GasPumpUpdateRequest{},
			Pre: func() {
				suite.testRequest.SetBearerToken("Bearer " + suite.validToken)
			},
			Pos: func() {
				suite.testRequest.SetBearerToken("")
			},
		},
		{
			Name:               "TestGasPumpController_UpdateConflict",
			Url:                url + suite.gasPumpConflictID.String(),
			ExpectedStatusCode: http.StatusNotAcceptable,
			ExpectedResponse:   dto.GeneralMessage{Detail: lang.NotAcceptable + "gas_station_id"},
			Body:               dto.GasPumpUpdateRequest{},
			Pre: func() {
				suite.testRequest.SetBearerToken("Bearer " + suite.validToken)
			},
			Pos: func() {
				suite.testRequest.SetBearerToken("")
			},
		},
		{
			Name:               "TestGasPumpController_UpdateDuplicated",
			Url:                url + suite.gasPumpDuplicatedID.String(),
			ExpectedStatusCode: http.StatusConflict,
			ExpectedResponse:   dto.GeneralMessage{Detail: lang.DuplicatedEntry + "gas_station_id, number"},
			Body:               dto.GasPumpUpdateRequest{},
			Pre: func() {
				suite.testRequest.SetBearerToken("Bearer " + suite.validToken)
			},
			Pos: func() {
				suite.testRequest.SetBearerToken("")
			},
		},
		{
			Name:               "TestGasPumpController_UpdateNotFound",
			Url:                url + suite.gasPumpNotFoundID.String(),
			ExpectedStatusCode: http.StatusNotFound,
			ExpectedResponse:   dto.GeneralMessage{Detail: lang.NotFoundRecord},
			Body:               dto.GasPumpUpdateRequest{},
			Pre: func() {
				suite.testRequest.SetBearerToken("Bearer " + suite.validToken)
			},
			Pos: func() {
				suite.testRequest.SetBearerToken("")
			},
		},
		{
			Name:               "TestGasPumpController_UpdateNotFound",
			Url:                url + suite.gasPumpID.String(),
			ExpectedStatusCode: http.StatusOK,
			ExpectedResponse:   dto.GeneralMessage{Detail: lang.RecordUpdated},
			Body: dto.GasPumpUpdateRequest{
				Number:       utils.StringAddr("10"),
				GasStationID: &fakeId,
				Active:       utils.BoolAddr(false),
			},
			Pre: func() {
				suite.testRequest.SetBearerToken("Bearer " + suite.validToken)
			},
			Pos: func() {
				suite.testRequest.SetBearerToken("")
			},
		},
	}

	t := suite.T()

	for i := range testcases {
		tc := testcases[i]

		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			suite.repository.Test(t)

			tc.Pre()

			res := suite.testRequest.Put(tc.Url, tc.Body)

			if tc.ExpectedStatusCode != 0 {
				suite.Equal(tc.ExpectedStatusCode, res.Code, utils.PrintExpectedValues(tc.ExpectedStatusCode, res.Code))
			}

			if tc.ExpectedResponse != nil {
				expected, _ := json.Marshal(tc.ExpectedResponse)
				suite.Equal(string(expected), res.Body.String(), utils.PrintExpectedValues(string(expected), res.Body.String()))
			}

			tc.Pos()
		})
	}
}

func (suite *gasPumpCtrlTest) TestList() {
	url := "/api/v1/gas-pumps"

	gasPumpsLimit := make([]dto.GasPumpListResponse, 0)

	copier.Copy(&gasPumpsLimit, suite.gasPumps[:len(suite.gasPumps)])

	gasPumpsRest := make([]dto.GasPumpListResponse, 0)

	copier.Copy(&gasPumpsRest, suite.gasPumps[20:21])

	gasPumps := make([]dto.GasPumpListResponse, 0)

	copier.Copy(&gasPumps, suite.gasPumps[:10])

	err := validate.Struct(dto.PaginateRequest{Page: -1, Limit: 101})

	expectedBadResponse := utils.MapValidatorError[dto.PaginateRequest](err)

	testcases := []struct {
		Url                string
		Name               string
		ExpectedResponse   any
		ExpectedStatusCode int
		Pre                func()
		Pos                func()
	}{
		{
			Name:               "TestGasPumpController_ListUnauthorized",
			Url:                url,
			ExpectedStatusCode: http.StatusUnauthorized,
			ExpectedResponse:   dto.GeneralMessage{Detail: lang.NotFoundRecord},
			Pre: func() {
				suite.testRequest.SetBearerToken("Bearer " + suite.invalidToken)
			},
			Pos: func() {
				suite.testRequest.SetBearerToken("")
			},
		},
		{
			Name:               "TestGasPumpController_ListBadRequest",
			Url:                url + "?page=-1&limit=101",
			ExpectedStatusCode: 400,
			ExpectedResponse:   expectedBadResponse,
			Pre: func() {
				suite.testRequest.SetBearerToken("Bearer " + suite.validToken)
			},
			Pos: func() {
				suite.testRequest.SetBearerToken("")
			},
		},
		{
			Name:               "TestGasPumpController_ListInternalServerError",
			Url:                url + "?page=1000&limit=100",
			ExpectedStatusCode: http.StatusInternalServerError,
			ExpectedResponse:   dto.GeneralMessage{Detail: lang.InternalServerError},
			Pre: func() {
				suite.testRequest.SetBearerToken("Bearer " + suite.validToken)
			},
			Pos: func() {
				suite.testRequest.SetBearerToken("")
			},
		},
		{
			Name:               "TestGasPumpController_ListEmpty",
			Url:                url + "?page=10&limit=10",
			ExpectedStatusCode: http.StatusOK,
			ExpectedResponse: dto.PaginationResponse{
				Page:       10,
				Limit:      10,
				TotalRows:  0,
				TotalPages: 0,
				Data:       []dto.GasPumpListResponse{},
			},
			Pre: func() {
				suite.testRequest.SetBearerToken("Bearer " + suite.validToken)
			},
			Pos: func() {
				suite.testRequest.SetBearerToken("")
			},
		},
		{
			Name:               "TestGasPumpController_ListLimit",
			Url:                url + "?page=1&limit=100",
			ExpectedStatusCode: http.StatusOK,
			ExpectedResponse: dto.PaginationResponse{
				Page:       1,
				Limit:      100,
				TotalRows:  int64(len(suite.gasPumps)),
				TotalPages: 1,
				Data:       gasPumpsLimit,
			},
			Pre: func() {
				suite.testRequest.SetBearerToken("Bearer " + suite.validToken)
			},
			Pos: func() {
				suite.testRequest.SetBearerToken("")
			},
		},
		{
			Name:               "TestGasPumpController_ListRestElements",
			Url:                url + "?page=3&limit=10",
			ExpectedStatusCode: http.StatusOK,
			ExpectedResponse: dto.PaginationResponse{
				Page:       3,
				Limit:      10,
				TotalRows:  int64(len(suite.gasPumps)),
				TotalPages: 3,
				Data:       gasPumpsRest,
			},
			Pre: func() {
				suite.testRequest.SetBearerToken("Bearer " + suite.validToken)
			},
			Pos: func() {
				suite.testRequest.SetBearerToken("")
			},
		},
		{
			Name:               "TestGasPumpController_List",
			Url:                url + "?page=1&limit=10",
			ExpectedStatusCode: http.StatusOK,
			ExpectedResponse: dto.PaginationResponse{
				Page:       1,
				Limit:      10,
				TotalRows:  int64(len(suite.gasPumps)),
				TotalPages: 3,
				Data:       gasPumps,
			},
			Pre: func() {
				suite.testRequest.SetBearerToken("Bearer " + suite.validToken)
			},
			Pos: func() {
				suite.testRequest.SetBearerToken("")
			},
		},
	}

	t := suite.T()

	for i := range testcases {
		tc := testcases[i]

		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			suite.repository.Test(t)

			tc.Pre()

			res := suite.testRequest.Get(tc.Url, nil)

			if tc.ExpectedStatusCode != 0 {
				suite.Equal(tc.ExpectedStatusCode, res.Code, utils.PrintExpectedValues(tc.ExpectedStatusCode, res.Code))
			}

			if tc.ExpectedResponse != nil {
				expected, _ := json.Marshal(tc.ExpectedResponse)
				suite.Equal(string(expected), res.Body.String(), utils.PrintExpectedValues(string(expected), res.Body.String()))
			}

			tc.Pos()
		})
	}
}

func (suite *gasPumpCtrlTest) TestCreate() {
	url := "/api/v1/gas-pumps"

	err := suite.validate.Struct(dto.GasPumpCreateRequest{})

	expectedBadResponse := utils.MapValidatorError[dto.GasPumpCreateRequest](err)

	testcases := []struct {
		Name               string
		Url                string
		Body               any
		ExpectedStatusCode int
		ExpectedResponse   any
		Pre                func()
		Pos                func()
	}{
		{
			Name:               "TestGasPumpController_CreateUnauthorized",
			Url:                url,
			ExpectedStatusCode: http.StatusUnauthorized,
			ExpectedResponse:   dto.GeneralMessage{Detail: lang.NotFoundRecord},
			Pre: func() {
				suite.testRequest.SetBearerToken("Bearer " + suite.invalidToken)
			},
			Pos: func() {
				suite.testRequest.SetBearerToken("")
			},
		},
		{
			Name:               "TestGasPumpController_CreateBadRequest",
			Url:                url,
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedResponse:   expectedBadResponse,
			Body:               dto.GasPumpCreateRequest{},
			Pre: func() {
				suite.testRequest.SetBearerToken("Bearer " + suite.validToken)
			},
			Pos: func() {
				suite.testRequest.SetBearerToken("")
			},
		},
		{
			Name: "TestGasPumpController_CreateDuplicated",
			Url:  url,
			Body: dto.GasPumpCreateRequest{
				Number:       "01",
				GasStationID: suite.gasPumpDuplicatedID.String(),
				ExternalID:   "01",
			},
			ExpectedStatusCode: http.StatusConflict,
			ExpectedResponse:   dto.GeneralMessage{Detail: lang.DuplicatedEntry + "gas_station_id, number"},
			Pre: func() {
				suite.testRequest.SetBearerToken("Bearer " + suite.validToken)
			},
			Pos: func() {
				suite.testRequest.SetBearerToken("")
			},
		},
		{
			Name: "TestGasPumpController_CreateNotAcceptable",
			Url:  url,
			Body: dto.GasPumpCreateRequest{
				Number:       "01",
				GasStationID: suite.gasPumpConflictID.String(),
				ExternalID:   "01",
			},
			ExpectedStatusCode: http.StatusNotAcceptable,
			ExpectedResponse:   dto.GeneralMessage{Detail: lang.NotAcceptable + "gas_station_id"},
			Pre: func() {
				suite.testRequest.SetBearerToken("Bearer " + suite.validToken)
			},
			Pos: func() {
				suite.testRequest.SetBearerToken("")
			},
		},
		{
			Name: "TestGasPumpController_CreateInternalServerError",
			Url:  url,
			Body: dto.GasPumpCreateRequest{
				Number:       "02",
				GasStationID: suite.gasPumpErrorID.String(),
				ExternalID:   "01",
			},
			ExpectedStatusCode: http.StatusInternalServerError,
			ExpectedResponse:   dto.GeneralMessage{Detail: lang.InternalServerError},
			Pre: func() {
				suite.testRequest.SetBearerToken("Bearer " + suite.validToken)
			},
			Pos: func() {
				suite.testRequest.SetBearerToken("")
			},
		},
		{
			Name: "TestGasPumpController_Create",
			Url:  url,
			Body: dto.GasPumpCreateRequest{
				Number:       "01",
				GasStationID: suite.gasPumpID.String(),
				ExternalID:   "01",
			},
			ExpectedStatusCode: http.StatusCreated,
			ExpectedResponse: dto.GasPumpCreateResponse{
				ID: suite.gasPumpID,
			},
			Pre: func() {
				suite.testRequest.SetBearerToken("Bearer " + suite.validToken)
			},
			Pos: func() {
				suite.testRequest.SetBearerToken("")
			},
		},
	}
	t := suite.T()

	for i := range testcases {
		tc := testcases[i]

		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			suite.repository.Test(t)

			tc.Pre()

			res := suite.testRequest.Post(tc.Url, tc.Body)

			if tc.ExpectedStatusCode != 0 {
				suite.Equal(tc.ExpectedStatusCode, res.Code, utils.PrintExpectedValues(tc.ExpectedStatusCode, res.Code))
			}

			if tc.ExpectedResponse != nil {
				expected, _ := json.Marshal(tc.ExpectedResponse)
				suite.Equal(string(expected), res.Body.String(), utils.PrintExpectedValues(string(expected), res.Body.String()))
			}

			tc.Pos()
		})
	}
}

func (suite *gasPumpCtrlTest) TestGetDetail() {
	url := "/api/v1/gas-pumps/"

	err := suite.validate.Struct(dto.GasPumpGetPathRequest{ID: "wrongformat"})

	expectedBadResponse := utils.MapValidatorError[dto.GasPumpGetPathRequest](err)

	var expectedOkResponse dto.GasPumpGetResponse

	copier.Copy(&expectedOkResponse, &suite.gasPump)

	testcases := []struct {
		Name               string
		Url                string
		ExpectedStatusCode int
		ExpectedResponse   any
		Pre                func()
		Pos                func()
	}{
		{
			Name:               "TestGasPumpController_GetUnauthorized",
			Url:                url + suite.gasPumpID.String(),
			ExpectedStatusCode: http.StatusUnauthorized,
			ExpectedResponse:   dto.GeneralMessage{Detail: lang.NotFoundRecord},
			Pre: func() {
				suite.testRequest.SetBearerToken("Bearer " + suite.invalidToken)
			},
			Pos: func() {
				suite.testRequest.SetBearerToken("")
			},
		},
		{
			Name:               "TestGasPumpController_GetBadRequest",
			Url:                url + "wrongformat",
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedResponse:   expectedBadResponse,
			Pre: func() {
				suite.testRequest.SetBearerToken("Bearer " + suite.validToken)
			},
			Pos: func() {
				suite.testRequest.SetBearerToken("")
			},
		},
		{
			Name:               "TestGasPumpController_GetNotFound",
			Url:                url + suite.gasPumpNotFoundID.String(),
			ExpectedStatusCode: http.StatusNotFound,
			ExpectedResponse:   dto.GeneralMessage{Detail: lang.NotFoundRecord},
			Pre: func() {
				suite.testRequest.SetBearerToken("Bearer " + suite.validToken)
			},
			Pos: func() {
				suite.testRequest.SetBearerToken("")
			},
		},
		{
			Name:               "TestGasPumpController_GetInternalServerError",
			Url:                url + suite.gasPumpErrorID.String(),
			ExpectedStatusCode: http.StatusInternalServerError,
			ExpectedResponse:   dto.GeneralMessage{Detail: lang.InternalServerError},
			Pre: func() {
				suite.testRequest.SetBearerToken("Bearer " + suite.validToken)
			},
			Pos: func() {
				suite.testRequest.SetBearerToken("")
			},
		},
		{
			Name:               "TestGasPumpController_Get",
			Url:                url + suite.gasPumpID.String(),
			ExpectedStatusCode: http.StatusOK,
			ExpectedResponse:   expectedOkResponse,
			Pre: func() {
				suite.testRequest.SetBearerToken("Bearer " + suite.validToken)
			},
			Pos: func() {
				suite.testRequest.SetBearerToken("")
			},
		},
	}

	t := suite.T()

	for i := range testcases {
		tc := testcases[i]

		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			suite.repository.Test(t)

			tc.Pre()

			res := suite.testRequest.Get(tc.Url, nil)

			if tc.ExpectedStatusCode != 0 {
				suite.Equal(tc.ExpectedStatusCode, res.Code, utils.PrintExpectedValues(tc.ExpectedStatusCode, res.Code))
			}

			if tc.ExpectedResponse != nil {
				expected, _ := json.Marshal(tc.ExpectedResponse)
				suite.Equal(string(expected), res.Body.String(), utils.PrintExpectedValues(string(expected), res.Body.String()))
			}

			tc.Pos()
		})
	}
}

func TestGasPumpController(t *testing.T) {
	suite.Run(t, new(gasPumpCtrlTest))
}
