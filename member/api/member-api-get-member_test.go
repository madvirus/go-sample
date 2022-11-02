package api

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go-sample/member/query"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockMemberQueryService struct {
	mock.Mock
}

func (m *mockMemberQueryService) GetMembers() ([]query.MemberData, error) {
	args := m.Called()
	return args.Get(0).([]query.MemberData), args.Error(1)
}

func (m *mockMemberQueryService) GetMember(id int64) (*query.MemberData, error) {
	args := m.Called(id)
	return args.Get(0).(*query.MemberData), args.Error(1)
}

func TestMemberApiQueryTestSuite(t *testing.T) {
	suite.Run(t, new(MemberApiQueryTestSuite))
}

type MemberApiQueryTestSuite struct {
	suite.Suite
	api              *MemberApi
	router           *gin.Engine
	mockQueryService *mockMemberQueryService
}

func (suite *MemberApiQueryTestSuite) SetupTest() {
	router := gin.Default()
	mockQueryService := new(mockMemberQueryService)
	suite.mockQueryService = mockQueryService
	suite.api = CreateMemberApi(mockQueryService, nil, nil)
	suite.api.RegisterHandler(router)
	suite.router = router
}

func (suite *MemberApiQueryTestSuite) TestGetMember() {
	date := carbon.Date{
		carbon.CreateFromDate(2022, 11, 1),
	}
	datetime := carbon.DateTime{
		carbon.CreateFromDateTime(2022, 11, 1, 9, 10, 11),
	}
	md := query.MemberData{
		Id:           1,
		Name:         "name",
		Email:        "email@email.com",
		BirthDate:    &date,
		RegisterDate: &datetime,
	}
	suite.mockQueryService.On("GetMember", mock.Anything).Return(&md, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/members/1", nil)
	suite.router.ServeHTTP(w, req)

	arg := suite.mockQueryService.Calls[0].Arguments.Get(0)
	requestId, ok := arg.(int64)
	assert.True(suite.T(), ok)
	if ok {
		assert.Equal(suite.T(), int64(1), requestId)
	}
}
