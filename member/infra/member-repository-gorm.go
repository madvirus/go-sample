package infra

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"go-sample/member/domain"
	"go-sample/transactor"
	"gorm.io/gorm"
)

func CreateMemberRepository(transactor transactor.Transactor) domain.MemberRepository {
	return &gormMemberRepository{
		transactor: transactor,
	}
}

type gormMemberRepository struct {
	transactor transactor.Transactor
}

func (repo *gormMemberRepository) Save(ctx context.Context, member *domain.Member) error {
	tx := repo.transactor.GetTx(ctx)
	result := tx.Save(member)
	if result.Error != nil {
		log.Errorf("fail to save %v", member)
		return result.Error
	}
	return nil
}

func (repo *gormMemberRepository) FindByEmail(ctx context.Context, email string) (*domain.Member, error) {
	tx := repo.transactor.GetTx(ctx)
	var m domain.Member
	result := tx.Take(&m, "email=?", email)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Infof("Member[email=%s] not found", email)
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &m, nil
}

func (repo *gormMemberRepository) FindById(ctx context.Context, id int64) (*domain.Member, error) {
	tx := repo.transactor.GetTx(ctx)
	var m domain.Member
	result := tx.Take(&m, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Infof("Member[id=%d] not found", id)
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &m, nil
}
