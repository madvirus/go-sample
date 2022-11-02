package app

import (
	"context"
	"fmt"
	"github.com/golang-module/carbon/v2"
	log "github.com/sirupsen/logrus"
	"go-sample/common"
	"go-sample/member/domain"
	"go-sample/transactor"
)

type MemberRegistRequest struct {
	Name      string       `json:"name"`
	Email     string       `json:"email"`
	Birthdate *carbon.Date `json:"birthdate"`
}

type RegisterService interface {
	Register(req MemberRegistRequest) (int64, error)
}

func CreateRegisterService(transactor transactor.Transactor, repository domain.MemberRepository) RegisterService {
	return &registerService{
		transactor:       transactor,
		memberRepository: repository,
	}
}

type registerService struct {
	transactor       transactor.Transactor
	memberRepository domain.MemberRepository
}

type MemberErrorCode int

const (
	MemberAlreadyExists MemberErrorCode = 1
)

type MemberError struct {
	Code MemberErrorCode
}

func (m MemberError) Error() string {
	return fmt.Sprintf("ErrorCode: %d", 1)
}

func (r *registerService) Register(req MemberRegistRequest) (int64, error) {
	errors := validateMemberRegisterRequest(req)
	if len(errors) > 0 {
		return 0, common.CreateValidationError(errors)
	}

	var id int64
	err := r.transactor.Execute(func(ctx context.Context) error {
		member, err := r.memberRepository.FindByEmail(ctx, req.Email)
		if err != nil {
			return err
		}
		if member != nil {
			log.Info("이미 같은 이메일 회원 존재")
			return MemberError{Code: MemberAlreadyExists}
		}
		var newMember = &domain.Member{
			Name:         req.Name,
			Email:        req.Email,
			BirthDate:    req.Birthdate,
			RegisterDate: carbon.DateTime{carbon.Now()},
		}
		err = r.memberRepository.Save(ctx, newMember)
		if err != nil {
			return err
		}
		id = newMember.Id
		return nil
	})

	if err != nil {
		return 0, err
	}
	return id, nil
}

func validateMemberRegisterRequest(req MemberRegistRequest) []common.ErrorField {
	var errors []common.ErrorField
	if len(req.Name) == 0 {
		errors = append(errors, common.ErrorField{Name: "name", Message: "이름은 필수입니다."})
	}
	if len(req.Email) == 0 {
		errors = append(errors, common.ErrorField{Name: "email", Message: "이메일은 필수입니다."})
	}
	return errors
}
