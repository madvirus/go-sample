//go:build integration

package test_integration

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go-sample/appctx"
	"go-sample/transactor"
	"gorm.io/gorm"
	"testing"
)

func TestTransactorIntTestSuite(t *testing.T) {
	suite.Run(t, new(TransactorIntTestSuite))
}

type TransactorIntTestSuite struct {
	suite.Suite
	transactor   transactor.Transactor
	testDbHelper *TestDbHelper
	db           *gorm.DB
}

func (suite *TransactorIntTestSuite) SetupSuite() {
	appCtx, err := CreateAppCtx4Test()
	if err != nil {
		panic(err.Error())
	}
	suite.transactor, _ = appctx.GetByType[transactor.Transactor](appCtx)
	suite.testDbHelper = CreateTestDbHelper(appCtx.DB())
	suite.db = appCtx.DB()
}

func (suite *TransactorIntTestSuite) SetupTest() {
	suite.testDbHelper.Truncate()
}

func (suite *TransactorIntTestSuite) TestRollback() {
	suite.transactor.Execute(func(ctx context.Context) error {
		tx := suite.transactor.GetTx(ctx)
		result := tx.Exec("insert into member (id, name, email) values (1, 'name1', 'email1')")
		if result.Error != nil {
			return result.Error
		}
		result = tx.Exec("insert into member (id, name, email) values (2, 'name2', 'email2')")
		if result.Error != nil {
			return result.Error
		}
		result = tx.Exec("insert into member (id, name, email) values (1, 'name3', 'email3')")
		if result.Error != nil {
			// 중복 PK, 에러 리턴, 롤백 처리
			return result.Error
		}
		return nil
	})
	var count int64
	suite.db.Table("member").Count(&count)
	assert.Equal(suite.T(), int64(0), count)
}
