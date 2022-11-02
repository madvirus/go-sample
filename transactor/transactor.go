package transactor

import (
	"context"
	"gorm.io/gorm"
)

type Transactor interface {
	Execute(run func(ctx context.Context) error) error
	ExecuteContext(ctx context.Context, run func(ctx context.Context) error) error
	GetTx(ctx context.Context) *gorm.DB
}
