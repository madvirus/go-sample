package domain

import (
	"context"
)

type MemberRepository interface {
	FindById(ctx context.Context, id int64) (*Member, error)
	FindByEmail(ctx context.Context, email string) (*Member, error)
	Save(ctx context.Context, member *Member) error
}
