//go:build integration

package member

import (
	"context"
	"github.com/golang-module/carbon/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go-sample/appctx"
	"go-sample/member/app"
	"go-sample/member/domain"
	"go-sample/test-integration"
	"testing"
)

func TestRegisterServiceIntTestSuite(t *testing.T) {
	suite.Run(t, new(RegisterServiceIntTestSuite))
}

type RegisterServiceIntTestSuite struct {
	suite.Suite
	registerService app.RegisterService
	repository      domain.MemberRepository
	testDbHelper    *test_integration.TestDbHelper
}

func (suite *RegisterServiceIntTestSuite) SetupSuite() {
	appCtx, err := test_integration.CreateAppCtx4Test()
	if err != nil {
		panic(err.Error())
	}
	registerService, err := appctx.GetByType[app.RegisterService](appCtx)
	if err != nil {
		panic(err.Error())
	}
	repository, err := appctx.GetByType[domain.MemberRepository](appCtx)
	if err != nil {
		panic(err.Error())
	}
	suite.registerService = registerService
	suite.repository = repository
	suite.testDbHelper = test_integration.CreateTestDbHelper(appCtx.DB())
}

func (suite *RegisterServiceIntTestSuite) SetupTest() {
	suite.testDbHelper.Truncate()
}

func (suite *RegisterServiceIntTestSuite) TestRegisterSuccess() {
	date := carbon.Date{Carbon: carbon.CreateFromDate(2022, 1, 1)}
	req := app.MemberRegistRequest{
		Name:      "name",
		Email:     "b@b.com",
		Birthdate: &date,
	}
	newId, err := suite.registerService.Register(req)
	if err != nil {
		assert.Fail(suite.T(), "error occured", "error: %v", err)
		return
	}
	mem, err := suite.repository.FindById(context.Background(), newId)
	if err != nil {
		assert.Fail(suite.T(), "error occured", "error: %v", err)
		return
	}
	assert.Equal(suite.T(), "name", mem.Name)
	assert.Equal(suite.T(), "b@b.com", mem.Email)
	assert.Equal(suite.T(), "2022-01-01", mem.BirthDate.ToDateString())
}

func (suite *RegisterServiceIntTestSuite) TestDupTest() {
	suite.testDbHelper.Save(&domain.Member{
		Name:         "n",
		Email:        "dup@dup.com",
		RegisterDate: carbon.DateTime{Carbon: carbon.Now()},
	})

	date := carbon.Date{Carbon: carbon.CreateFromDate(2022, 1, 1)}
	req := app.MemberRegistRequest{
		Name:      "name",
		Email:     "dup@dup.com",
		Birthdate: &date,
	}
	_, err := suite.registerService.Register(req)
	assert.NotNil(suite.T(), err)
	memberError, ok := err.(app.MemberError)
	assert.True(suite.T(), ok, "err should be MemberError")
	assert.Equal(suite.T(), memberError.Code, app.MemberAlreadyExists)
}
