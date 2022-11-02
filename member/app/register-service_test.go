package app

import (
	"github.com/golang-module/carbon/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go-sample/common"
	"go-sample/member/domain"
	"go-sample/member/infra"
	"go-sample/transactor"
	"testing"
)

func TestRegisterServiceTestSuite(t *testing.T) {
	suite.Run(t, new(RegisterServiceTestSuite))
}

type RegisterServiceTestSuite struct {
	suite.Suite
	repository domain.MemberRepository
	service    RegisterService
}

func (suite *RegisterServiceTestSuite) SetupTest() {
	transactor := transactor.CreateDummyTransactor()
	suite.repository = infra.CreateMemoryMemberRepository()
	suite.service = CreateRegisterService(transactor, suite.repository)
}

func (suite *RegisterServiceTestSuite) TestInvalidRequest() {
	_, err := suite.service.Register(MemberRegistRequest{})

	ve, ok := err.(*common.ValidationError)
	assert.True(suite.T(), ok)
	if ok {
		assert.NotEmpty(suite.T(), ve.Errors)
	}
}

func (suite *RegisterServiceTestSuite) Test_RegisterService_Success() {
	date := carbon.Date{Carbon: carbon.CreateFromDate(2022, 10, 19)}
	id, err := suite.service.Register(MemberRegistRequest{
		Name:      "이름",
		Email:     "a@a.com",
		Birthdate: &date,
	})
	assert.Nil(suite.T(), err)
	mem, err := suite.repository.FindById(nil, id)
	assert.NotNil(suite.T(), mem)
	assert.Equal(suite.T(), "이름", mem.Name)
	assert.Equal(suite.T(), "a@a.com", mem.Email)
	assert.Equal(suite.T(), "2022-10-19", mem.BirthDate.ToDateString())
}

func (suite *RegisterServiceTestSuite) Test_RegisterService_DupEmail() {
	// 이메일 a@a.com 회원 존재
	suite.repository.Save(nil, &domain.Member{
		Id:    1,
		Name:  "name",
		Email: "a@a.com",
	})

	date := carbon.Date{Carbon: carbon.CreateFromDate(2022, 10, 19)}

	_, err := suite.service.Register(MemberRegistRequest{
		Name:      "이름",
		Email:     "a@a.com",
		Birthdate: &date,
	})

	assert.NotNil(suite.T(), err)
	memberError, ok := err.(MemberError)
	assert.True(suite.T(), ok, "err should be MemberError")
	assert.Equal(suite.T(), memberError.Code, MemberAlreadyExists)
}
