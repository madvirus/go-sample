package transactor

import (
	"context"
	"gorm.io/gorm"
)

func CreateTransactor(db *gorm.DB) Transactor {
	return &gormTransactor{db: db}
}

const ctx_tx_name = "CTX_GORM_TX"

type gormTransactor struct {
	db *gorm.DB
}

func (t *gormTransactor) GetTx(ctx context.Context) *gorm.DB {
	ongoingTx := ctx.Value(ctx_tx_name)
	if ongoingTx == nil {
		return t.db
	}
	return ongoingTx.(*gorm.DB)
}

func (t *gormTransactor) Execute(run func(ctx context.Context) error) error {
	ctx := context.Background()
	return t.ExecuteContext(ctx, run)
}

func (t *gormTransactor) ExecuteContext(ctx context.Context, run func(ctx context.Context) error) error {
	value := ctx.Value(ctx_tx_name)
	_, ok := value.(*gorm.DB)
	if value != nil && ok {
		return run(ctx)
	}
	return t.db.Transaction(func(tx *gorm.DB) error {
		newCtx := context.WithValue(ctx, ctx_tx_name, tx)
		return run(newCtx)
	})
}
