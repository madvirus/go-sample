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

type mockUpdateService struct {
	mock.Mock
}

func (m *mockUpdateService) UpdateMember(req app.UpdateMemberRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func TestMemberApiUpdateTestSuite(t *testing.T) {
	suite.Run(t, new(MemberApiUpdateTestSuite))
}

type MemberApiUpdateTestSuite struct {
	suite.Suite
	api               *MemberApi
	router            *gin.Engine
	mockUpdateService *mockUpdateService
}

func (suite *MemberApiUpdateTestSuite) SetupTest() {
	router := gin.Default()
	mockUpdateService := new(mockUpdateService)
	suite.mockUpdateService = mockUpdateService
	suite.api = CreateMemberApi(nil, nil, mockUpdateService)
	suite.api.RegisterHandler(router)
	suite.router = router
}

func (suite *MemberApiUpdateTestSuite) TestUpdateRequestPassing() {
	suite.mockUpdateService.On("UpdateMember", mock.Anything).Return(nil)

	w := httptest.NewRecorder()
	json := `{"id": 1, "name": "name", "birthdate": "2022-01-01"}`
	req, _ := http.NewRequest(http.MethodPut, "/members", bytes.NewBuffer([]byte(json)))
	suite.router.ServeHTTP(w, req)

	arg := suite.mockUpdateService.Calls[0].Arguments.Get(0)
	request, ok := arg.(app.UpdateMemberRequest)
	assert.True(suite.T(), ok)
	if ok {
		assert.Equal(suite.T(), int64(1), request.Id)
		assert.Equal(suite.T(), "name", request.Name)
		assert.Equal(suite.T(), "2022-01-01", request.Birthdate.ToDateString())
	}
}
