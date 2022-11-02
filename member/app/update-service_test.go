package app

import (
	"context"
	"github.com/golang-module/carbon/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go-sample/common"
	"go-sample/member/domain"
	"go-sample/member/infra"
	"go-sample/transactor"
	"testing"
)

func TestUpdateMemberServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UpdateMemberServiceTestSuite))
}

type UpdateMemberServiceTestSuite struct {
	suite.Suite
	repository domain.MemberRepository
	service    UpdateService
}

func (suite *UpdateMemberServiceTestSuite) SetupTest() {
	transactor := transactor.CreateDummyTransactor()
	suite.repository = infra.CreateMemoryMemberRepository()
	suite.service = CreateUpdateService(transactor, suite.repository)
}

func (suite *UpdateMemberServiceTestSuite) TestInvalidRequest() {
	err := suite.service.UpdateMember(UpdateMemberRequest{})
	ve, ok := err.(*common.ValidationError)
	assert.True(suite.T(), ok)
	if ok {
		assert.NotEmpty(suite.T(), ve.Errors)
	}
}

func (suite *UpdateMemberServiceTestSuite) TestNotFound() {
	err := suite.service.UpdateMember(UpdateMemberRequest{
		Id:        100,
		Name:      "newname",
		Birthdate: nil,
	})

	_, ok := err.(*common.NoDataFoundError)
	assert.True(suite.T(), ok)
}

func (suite *UpdateMemberServiceTestSuite) TestUpdated() {
	birth := carbon.Date{Carbon: carbon.CreateFromDate(2022, 10, 11)}
	suite.repository.Save(context.Background(), &domain.Member{
		Id:           101,
		Name:         "oldname",
		Email:        "a@a.com",
		BirthDate:    &birth,
		RegisterDate: carbon.DateTime{carbon.Now()},
	})

	newBirthDate := carbon.Date{Carbon: carbon.CreateFromDate(2022, 10, 29)}
	suite.service.UpdateMember(UpdateMemberRequest{
		Id:        101,
		Name:      "newname",
		Birthdate: &newBirthDate,
	})

	mem, _ := suite.repository.FindById(context.Background(), 101)
	assert.Equal(suite.T(), "newname", mem.Name)
	assert.Equal(suite.T(), "2022-10-29", mem.BirthDate.ToDateString())
}
