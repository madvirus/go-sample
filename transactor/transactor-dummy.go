package transactor

import (
	"context"
	"gorm.io/gorm"
)

func CreateDummyTransactor() Transactor {
	return &dummyTransactor{}
}

type dummyTransactor struct {
}

func (d *dummyTransactor) Execute(run func(ctx context.Context) error) error {
	ctx := context.Background()
	return d.ExecuteContext(ctx, run)
}

func (d *dummyTransactor) ExecuteContext(ctx context.Context, run func(ctx context.Context) error) error {
	return run(ctx)
}

func (d *dummyTransactor) GetTx(ctx context.Context) *gorm.DB {
	return nil
}
