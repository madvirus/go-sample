package infra

import (
	"context"
	"go-sample/member/domain"
)

func CreateMemoryMemberRepository() domain.MemberRepository {
	repository := memoryMemberRepository{
		values: make(map[int64]*domain.Member),
	}
	return &repository
}

type memoryMemberRepository struct {
	values map[int64]*domain.Member
}

func (m *memoryMemberRepository) FindById(ctx context.Context, id int64) (*domain.Member, error) {
	member, exists := m.values[id]
	if !exists {
		return nil, nil
	}
	return member, nil
}

func (m *memoryMemberRepository) FindByEmail(ctx context.Context, email string) (*domain.Member, error) {
	for _, member := range m.values {
		if member.Email == email {
			return member, nil
		}
	}
	return nil, nil
}

func (m *memoryMemberRepository) Save(ctx context.Context, member *domain.Member) error {
	m.values[member.Id] = member
	return nil
}
