package api

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go-sample/member/app"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockRegisterService struct {
	mock.Mock
}

func (m *mockRegisterService) Register(req app.MemberRegistRequest) (int64, error) {
	args := m.Called(req)
	return args.Get(0).(int64), args.Error(1)
}

func TestMemberApiTestSuite(t *testing.T) {
	suite.Run(t, new(MemberApiRegisterTestSuite))
}

type MemberApiRegisterTestSuite struct {
	suite.Suite
	api                 *MemberApi
	router              *gin.Engine
	mockRegisterService *mockRegisterService
}

func (suite *MemberApiRegisterTestSuite) SetupTest() {
	router := gin.Default()
	mockRegisterService := new(mockRegisterService)
	suite.mockRegisterService = mockRegisterService
	suite.api = CreateMemberApi(nil, mockRegisterService, nil)
	suite.api.RegisterHandler(router)
	suite.router = router
}

func (suite *MemberApiRegisterTestSuite) TestRegisterRequestPassing() {
	suite.mockRegisterService.On("Register", mock.Anything).Return(int64(1), nil)

	w := httptest.NewRecorder()
	json := `{"name": "name", "email": "a@a.com", "birthdate": "2022-01-01"}`
	req, _ := http.NewRequest(http.MethodPost, "/members", bytes.NewBuffer([]byte(json)))
	suite.router.ServeHTTP(w, req)

	arg := suite.mockRegisterService.Calls[0].Arguments.Get(0)
	request, ok := arg.(app.MemberRegistRequest)
	assert.True(suite.T(), ok)
	if ok {
		assert.Equal(suite.T(), "name", request.Name)
		assert.Equal(suite.T(), "a@a.com", request.Email)
		assert.Equal(suite.T(), "2022-01-01", request.Birthdate.ToDateString())
	}
}
