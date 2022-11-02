package appctx

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	memapi "go-sample/member/api"
	memapp "go-sample/member/app"
	meminfra "go-sample/member/infra"
	memquery "go-sample/member/query"
	"go-sample/transactor"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type AppCtx interface {
	Get(name string) (any, error)
	GetAll() []any
	DB() *gorm.DB
}

type appCtxData struct {
	db         *gorm.DB
	components map[string]any
}

func (ac *appCtxData) Get(name string) (any, error) {
	comp, found := ac.components[name]
	if !found {
		return nil, fmt.Errorf("component %s not found", name)
	}
	return comp, nil
}

func (ac *appCtxData) GetAll() []any {
	result := []any{}
	for _, v := range ac.components {
		result = append(result, v)
	}
	return result
}

func (ac *appCtxData) DB() *gorm.DB {
	return ac.db
}

// Initialize config를 이용해서 AppCtx를 생성
func Initialize(config Config) (AppCtx, error) {
	db, err := createDb(config)
	if err != nil {
		return nil, err
	}
	return createAppCtx(db)
}

func createDb(config Config) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(config.DbDns), &gorm.Config{})
	if err != nil {
		log.Fatalf("createDb: Failed to connect database: %v", err)
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("createDb: Failed to get db.DB(): %v", err)
		return nil, err
	}
	sqlDB.SetMaxIdleConns(config.DbMaxIdleConns)
	sqlDB.SetMaxOpenConns(config.DbMaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(config.DbConnMaxLifetime) * time.Second)

	log.Infof("db configuration: MaxIdleConns=%d", config.DbMaxIdleConns)
	log.Infof("db configuration: MaxOpenConns=%d", config.DbMaxOpenConns)
	log.Infof("db configuration: ConnMaxLifetime=%d", config.DbConnMaxLifetime)

	return db, err
}

func createAppCtx(db *gorm.DB) (AppCtx, error) {
	transactor := transactor.CreateTransactor(db)
	memberRepository := meminfra.CreateMemberRepository(transactor)
	registerService := memapp.CreateRegisterService(transactor, memberRepository)
	updateService := memapp.CreateUpdateService(transactor, memberRepository)
	memberQueryService := memquery.CreateMemberQueryService(db)
	memberApi := memapi.CreateMemberApi(memberQueryService, registerService, updateService)

	components := make(map[string]any)
	components["transactor"] = transactor
	components["memberRepository"] = memberRepository
	components["registerService"] = registerService
	components["updateService"] = updateService
	components["memberQueryService"] = memberQueryService
	components["memberApi"] = memberApi
	return &appCtxData{db: db, components: components}, nil
}

func Get[T any](ac AppCtx, name string) (T, error) {
	comp, err := ac.Get(name)
	if err != nil {
		var t T
		return t, err
	}
	typedComp, ok := comp.(T)
	if !ok {
		var t T
		return t, fmt.Errorf("%s component type: %T", name, comp)
	}
	return typedComp, nil
}

func GetByType[T any](ac AppCtx) (T, error) {
	comps := ac.GetAll()
	for _, c := range comps {
		if v, ok := c.(T); ok {
			return v, nil
		}
	}
	var t T
	return t, fmt.Errorf("component not found")
}

func GetAllByType[T any](ac AppCtx) []T {
	comps := ac.GetAll()
	result := []T{}
	for _, c := range comps {
		if v, ok := c.(T); ok {
			result = append(result, v)
		}
	}
	return result
}
