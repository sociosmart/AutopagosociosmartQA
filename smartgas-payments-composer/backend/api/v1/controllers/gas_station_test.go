package controllers_test

import (
	"encoding/json"
	"errors"
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

var validate = validator.New()

type gasStationCtrlTest struct {
	suite.Suite
	repository             *repository.MockGasStationRepository
	userRepository         *repository.MockUserRepository
	userID                 uuid.UUID
	gasStationID           uuid.UUID
	gasStationErrorID      uuid.UUID
	gasStationDuplicatedID uuid.UUID
	gasStation             *models.GasStation
	invalidToken           string
	validToken             string
	testRequest            *utils.TestRequest
	createGasStation       *models.GasStation
	createDuplicatedEntry  *models.GasStation
	gasStations            []*models.GasStation
}

func (suite *gasStationCtrlTest) SetupSuite() {
	setup, _ := injectors.InitializeServerWithMocks()

	suite.repository = setup.GasStationRepositoryMock
	suite.userRepository = setup.UserRepositoryMock

	suite.testRequest = &utils.TestRequest{
		Router: setup.Router,
	}

	suite.userID = uuid.New()
	suite.gasStationID = uuid.New()
	suite.gasStationErrorID = uuid.New()
	suite.gasStationDuplicatedID = uuid.New()

	suite.gasStations = make([]*models.GasStation, 0)

	// map gas stations
	for i := 1; i <= 22; i++ {
		suite.gasStations = append(suite.gasStations, &models.GasStation{
			Name: "Gasolinera " + strconv.Itoa(i),
			Ip:   "192.168.100." + strconv.Itoa(i),
		})
	}

	user := &models.User{
		ID:      suite.userID,
		IsAdmin: utils.BoolAddr(true),
	}

	suite.gasStation = &models.GasStation{
		ID:     suite.gasStationID,
		Name:   "Estacon don",
		Ip:     "192.168.100.100",
		Active: utils.BoolAddr(true),
	}

	mysqlErr := &mysql.MySQLError{
		Number: 1062,
	}

	suite.userRepository.On("GetUserByID", suite.userID).Return(user, nil)
	suite.userRepository.On("GetUserByID", mock.AnythingOfType("uuid.UUID")).Return(nil, gorm.ErrRecordNotFound)

	// Get by id
	suite.repository.On("GetByID", suite.gasStationID).Return(suite.gasStation, nil)
	suite.repository.On("GetByID", suite.gasStationErrorID).Return(nil, errors.New(lang.InternalServerError))
	suite.repository.On("GetByID", mock.AnythingOfType("uuid.UUID")).Return(nil, gorm.ErrRecordNotFound)

	// List

	suite.repository.On("List", mock.MatchedBy(func(pagination *schemas.Pagination) bool {
		if pagination.Limit == 10 && pagination.Page == 1 {
			pagination.TotalPages = 3
			pagination.TotalRows = int64(len(suite.gasStations))
			return true
		}

		return false
	}), mock.Anything).Return(suite.gasStations[:10], nil).Once()

	suite.repository.On("List", mock.MatchedBy(func(pagination *schemas.Pagination) bool {
		if pagination.Limit == 10 && pagination.Page == 3 {
			pagination.TotalPages = 3
			pagination.TotalRows = int64(len(suite.gasStations))
			return true
		}
		return false
	}), mock.Anything).Return(suite.gasStations[20:21], nil).Once()

	suite.repository.On("List", mock.MatchedBy(func(pagination *schemas.Pagination) bool {
		if pagination.Limit >= 100 && pagination.Page == 1 {
			pagination.TotalPages = 1
			pagination.TotalRows = int64(len(suite.gasStations))
			return true
		}
		return false
	}), mock.Anything).Return(suite.gasStations[:len(suite.gasStations)], nil)

	suite.repository.On("List", mock.MatchedBy(func(pagination *schemas.Pagination) bool {
		if pagination.Limit == 10 && pagination.Page == 10 {
			return true
		}
		return false
	}), mock.Anything).Return([]*models.GasStation{}, nil).Once()

	// Simulate internal server error
	suite.repository.On("List", mock.MatchedBy(func(pagination *schemas.Pagination) bool {
		if pagination.Limit == 100 && pagination.Page == 1000 {
			return true
		}
		return false
	}), mock.Anything).Return([]*models.GasStation{}, errors.New(lang.InternalServerError))

	// Create

	suite.createGasStation = &models.GasStation{ID: suite.gasStationID, Name: "Guerrero", Ip: "111.111.111.111", ExternalID: "11", CrePermission: "PL/21"}
	suite.repository.On("Create", mock.MatchedBy(func(gs *models.GasStation) bool {
		if gs.Name == suite.createGasStation.Name && gs.Ip == suite.createGasStation.Ip {
			gs.ID = suite.gasStationID
			return true
		}
		return false
	})).Return(nil)
	suite.createDuplicatedEntry = &models.GasStation{Name: "Degollado", Ip: "222.222.222.222", ExternalID: "17", CrePermission: "PL/21"}
	suite.repository.On("Create", mock.MatchedBy(func(gs *models.GasStation) bool {
		if gs.Name == suite.createDuplicatedEntry.Name && gs.Ip == suite.createDuplicatedEntry.Ip {
			return true
		}
		return false
	})).Return(mysqlErr)
	suite.repository.On("Create", mock.MatchedBy(func(gs *models.GasStation) bool {
		if gs.Name == "Error" {
			return true
		}
		return false
	})).Return(errors.New(lang.InternalServerError))

	// Update
	suite.repository.On("UpdateByID", suite.gasStationID, mock.AnythingOfType("*models.GasStation")).Return(true, nil)
	suite.repository.On("UpdateByID", suite.gasStationErrorID, mock.AnythingOfType("*models.GasStation")).Return(false, errors.New(lang.InternalServerError))
	suite.repository.On("UpdateByID", suite.gasStationDuplicatedID, mock.AnythingOfType("*models.GasStation")).Return(false, mysqlErr)
	suite.repository.On("UpdateByID", mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("*models.GasStation")).Return(false, nil)

	// ListAll
	suite.repository.On("ListAll").Return([]*models.GasStation{}, errors.New(lang.InternalServerError)).Once()
	suite.repository.On("ListAll").Return([]*models.GasStation{}, nil).Once()
	suite.repository.On("ListAll").Return(suite.gasStations, nil).Once()

	claims := &schemas.JwtClaims{
		Sub: suite.userID,
	}
	token, _ := claims.ClaimToken()
	suite.validToken = token

	claims.Sub = uuid.New()
	token, _ = claims.ClaimToken()
	suite.invalidToken = token
}

func (suite *gasStationCtrlTest) TestListAll() {
	url := "/api/v1/gas-stations/all"

	expectedResponse := make([]dto.GasStationListAllResponse, 0)

	copier.Copy(&expectedResponse, &suite.gasStations)

	testcases := []struct {
		Url                string
		Name               string
		ExpectedResponse   any
		ExpectedStatusCode int
		Pre                func()
		Pos                func()
	}{
		{
			Name:               "GasStationControllerTest_ListAllUnauthorized",
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
			Name:               "GasStationControllerTest_ListAllInternalServerError",
			Url:                url,
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
			Name:               "GasStationControllerTest_ListAllEmpty",
			Url:                url,
			ExpectedStatusCode: http.StatusOK,
			ExpectedResponse:   []dto.GasStationListAllResponse{},
			Pre: func() {
				suite.testRequest.SetBearerToken("Bearer " + suite.validToken)
			},
			Pos: func() {
				suite.testRequest.SetBearerToken("")
			},
		},
		{
			Name:               "GasStationControllerTest_ListAll",
			Url:                url,
			ExpectedStatusCode: http.StatusOK,
			ExpectedResponse:   expectedResponse,
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

func (suite *gasStationCtrlTest) TestGasStationList() {
	url := "/api/v1/gas-stations"

	gasStationsLimit := make([]dto.GasStationListResponse, 0)

	copier.Copy(&gasStationsLimit, suite.gasStations[:len(suite.gasStations)])

	gasStationsRest := make([]dto.GasStationListResponse, 0)

	copier.Copy(&gasStationsRest, suite.gasStations[20:21])

	gasStations := make([]dto.GasStationListResponse, 0)

	copier.Copy(&gasStations, suite.gasStations[:10])

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
			Name:               "GasStationControllerTest_ListUnauthorized",
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
			Name:               "GasStationControllerTest_ListBadRequest",
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
			Name:               "GasStationControllerTest_ListInternalServerError",
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
			Name:               "GasStationControllerTest_ListEmpty",
			Url:                url + "?page=10&limit=10",
			ExpectedStatusCode: http.StatusOK,
			ExpectedResponse: dto.PaginationResponse{
				Page:       10,
				Limit:      10,
				TotalRows:  0,
				TotalPages: 0,
				Data:       []dto.GasStationListResponse{},
			},
			Pre: func() {
				suite.testRequest.SetBearerToken("Bearer " + suite.validToken)
			},
			Pos: func() {
				suite.testRequest.SetBearerToken("")
			},
		},
		{
			Name:               "GasStationControllerTest_ListLimit",
			Url:                url + "?page=1&limit=100",
			ExpectedStatusCode: http.StatusOK,
			ExpectedResponse: dto.PaginationResponse{
				Page:       1,
				Limit:      100,
				TotalRows:  int64(len(suite.gasStations)),
				TotalPages: 1,
				Data:       gasStationsLimit,
			},
			Pre: func() {
				suite.testRequest.SetBearerToken("Bearer " + suite.validToken)
			},
			Pos: func() {
				suite.testRequest.SetBearerToken("")
			},
		},
		{
			Name:               "GasStationControllerTest_ListRestElements",
			Url:                url + "?page=3&limit=10",
			ExpectedStatusCode: http.StatusOK,
			ExpectedResponse: dto.PaginationResponse{
				Page:       3,
				Limit:      10,
				TotalRows:  int64(len(suite.gasStations)),
				TotalPages: 3,
				Data:       gasStationsRest,
			},
			Pre: func() {
				suite.testRequest.SetBearerToken("Bearer " + suite.validToken)
			},
			Pos: func() {
				suite.testRequest.SetBearerToken("")
			},
		},
		{
			Name:               "GasStationControllerTest_List",
			Url:                url + "?page=1&limit=10",
			ExpectedStatusCode: http.StatusOK,
			ExpectedResponse: dto.PaginationResponse{
				Page:       1,
				Limit:      10,
				TotalRows:  int64(len(suite.gasStations)),
				TotalPages: 3,
				Data:       gasStations,
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

func (suite *gasStationCtrlTest) TestGetGasStationDetail() {
	url := "/api/v1/gas-stations/"

	err := validate.Struct(dto.GasStationGetPathRequest{ID: "wrongformat"})

	badRequestExpectedStruct := utils.MapValidatorError[dto.GasStationGetPathRequest](err)

	var gasStation dto.GasStationGetResponse

	copier.Copy(&gasStation, &suite.gasStation)
	testcases := []struct {
		Url                string
		Name               string
		ExpectedResponse   any
		ExpectedStatusCode int
		Pre                func()
		Pos                func()
	}{
		{
			Name:               "GasStationControllerTest_GetDetailUnauthorized",
			Url:                url + suite.gasStationID.String(),
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
			Name:               "GasStationControllerTest_GetDetailBadRequest",
			Url:                url + "wrongformat",
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedResponse:   badRequestExpectedStruct,
			Pre: func() {
				suite.testRequest.SetBearerToken("Bearer " + suite.validToken)
			},
			Pos: func() {
				suite.testRequest.SetBearerToken("")
			},
		},
		{
			Name:               "GasStationControllerTest_GetDetailNotFound",
			Url:                url + uuid.NewString(),
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
			Name:               "GasStationControllerTest_GetDetailInternalServerError",
			Url:                url + suite.gasStationErrorID.String(),
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
			Name:               "GasStationControllerTest_GetDetail",
			Url:                url + suite.gasStationID.String(),
			ExpectedResponse:   gasStation,
			ExpectedStatusCode: http.StatusOK,
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

func (suite *gasStationCtrlTest) TestGasStationUpdate() {
	url := "/api/v1/gas-stations/"

	err := validate.Struct(dto.GasStationUpdatePathRequest{ID: "wrongformat"})

	expectedBadResponse := utils.MapValidatorError[dto.GasStationUpdatePathRequest](err)

	err = validate.Struct(dto.GasStationUpdateRequest{Ip: "wrongformat"})

	expectedBadResponseBody := utils.MapValidatorError[dto.GasStationUpdateRequest](err)

	var gasStationRequest dto.GasStationUpdateRequest

	copier.Copy(&gasStationRequest, &suite.gasStation)

	testcases := []struct {
		Url                string
		Name               string
		Body               any
		ExpectedResponse   any
		ExpectedStatusCode int
		Pre                func()
		Pos                func()
	}{
		{
			Name:               "GasStationControllerTest_UpdateUnauthorized",
			Url:                url + suite.gasStationID.String(),
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
			Name:               "GasStationControllerTest_UpdatePathBadRequest",
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
			Name:               "GasStationControllerTest_UpdateBodyBadRequest",
			Url:                url + suite.gasStationID.String(),
			Body:               dto.GasStationUpdateRequest{Ip: "wrongformat"},
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedResponse:   expectedBadResponseBody,
			Pre: func() {
				suite.testRequest.SetBearerToken("Bearer " + suite.validToken)
			},
			Pos: func() {
				suite.testRequest.SetBearerToken("")
			},
		},
		{
			Name:               "GasStationControllerTest_UpdateInternalServerError",
			Url:                url + suite.gasStationErrorID.String(),
			Body:               dto.GasStationUpdateRequest{Ip: "192.168.100.100", Name: "Cardenas", ExternalID: "17"},
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
			Name:               "GasStationControllerTest_UpdateNotFound",
			Url:                url + uuid.NewString(),
			Body:               dto.GasStationUpdateRequest{Ip: "192.168.100.100", Name: "Cardenas", ExternalID: "18"},
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
			Name:               "GasStationControllerTest_UpdateDupliatedEntry",
			Url:                url + suite.gasStationDuplicatedID.String(),
			Body:               dto.GasStationUpdateRequest{Ip: "192.168.100.100", Name: "Cardenas", ExternalID: "19"},
			ExpectedStatusCode: http.StatusConflict,
			ExpectedResponse:   dto.GeneralMessage{Detail: lang.DuplicatedEntry + "name, ip"},
			Pre: func() {
				suite.testRequest.SetBearerToken("Bearer " + suite.validToken)
			},
			Pos: func() {
				suite.testRequest.SetBearerToken("")
			},
		},
		{
			Name:               "GasStationControllerTest_Update",
			Url:                url + suite.gasStationID.String(),
			Body:               gasStationRequest,
			ExpectedStatusCode: http.StatusOK,
			ExpectedResponse:   dto.GeneralMessage{Detail: lang.RecordUpdated},
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

func (suite *gasStationCtrlTest) TestGasStationCreate() {

	url := "/api/v1/gas-stations"

	err := validate.Struct(dto.GasStationCreateRequest{})

	expectedBadResponse := utils.MapValidatorError[dto.GasStationCreateRequest](err)

	var duplicatedEntryRequest dto.GasStationCreateRequest

	copier.Copy(&duplicatedEntryRequest, &suite.createDuplicatedEntry)

	var gasStationRequest dto.GasStationCreateRequest

	copier.Copy(&gasStationRequest, &suite.createGasStation)

	testcases := []struct {
		Url                string
		Name               string
		Body               any
		ExpectedResponse   any
		ExpectedStatusCode int
		Pre                func()
		Pos                func()
	}{
		{
			Name:               "GasStationControllerTest_CreateUnauthorized",
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
			Name:               "GasStationControllerTest_CreateBadRequest",
			Url:                url,
			Body:               dto.GasStationCreateRequest{},
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
			Name:               "GasStationControllerTest_CreateDuplicatedEntry",
			Url:                url,
			Body:               duplicatedEntryRequest,
			ExpectedStatusCode: http.StatusConflict,
			ExpectedResponse:   dto.GeneralMessage{Detail: lang.DuplicatedEntry + "name, ip"},
			Pre: func() {
				suite.testRequest.SetBearerToken("Bearer " + suite.validToken)
			},
			Pos: func() {
				suite.testRequest.SetBearerToken("")
			},
		},
		{
			Name: "GasStationControllerTest_CreateInternalServerError",
			Url:  url,
			Body: dto.GasStationCreateRequest{
				Name:          "Error",
				Ip:            "123.123.123.123",
				ExternalID:    "13",
				CrePermission: "PL/21",
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
			Name:               "GasStationControllerTest_Create",
			Url:                url,
			Body:               gasStationRequest,
			ExpectedStatusCode: http.StatusCreated,
			ExpectedResponse:   dto.GasStationCreateResponse{ID: suite.gasStationID},
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

func TestGasStationController(t *testing.T) {
	suite.Run(t, new(gasStationCtrlTest))
}
