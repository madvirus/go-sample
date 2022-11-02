//go:build integration

package test_integration

import (
	"go-sample/appctx"
)

func CreateAppCtx4Test() (appctx.AppCtx, error) {
	config := appctx.Config{
		DbDns:             "testuser:testuser@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local",
		DbMaxIdleConns:    2,
		DbMaxOpenConns:    2,
		DbConnMaxLifetime: 60,
	}
	return appctx.Initialize(config)
}
