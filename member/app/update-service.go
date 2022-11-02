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

type UpdateMemberRequest struct {
	Id        int64        `json:"id"`
	Name      string       `json:"name"`
	Birthdate *carbon.Date `json:"birthdate"`
}

type UpdateService interface {
	UpdateMember(req UpdateMemberRequest) error
}

func CreateUpdateService(
	transactor transactor.Transactor,
	repository domain.MemberRepository,
) UpdateService {
	return &updateService{
		transactor:       transactor,
		memberRepository: repository,
	}
}

type updateService struct {
	transactor       transactor.Transactor
	memberRepository domain.MemberRepository
}

func (u *updateService) UpdateMember(req UpdateMemberRequest) error {
	errors := validateUpdateMemberRequest(req)
	if len(errors) > 0 {
		return common.CreateValidationError(errors)
	}

	err := u.transactor.Execute(func(ctx context.Context) error {
		member, err := u.memberRepository.FindById(ctx, req.Id)
		if err != nil {
			return err
		}
		if member == nil {
			log.Infof("UpdateMember: Member[id=%d]가 존재하지 않음", req.Id)
			return common.CreateNoDataFoundError(fmt.Sprintf("UpdateMember: Member[id=%d]가 존재하지 않음", req.Id))
		}
		member.Update(req.Name, req.Birthdate)

		err = u.memberRepository.Save(ctx, member)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func validateUpdateMemberRequest(req UpdateMemberRequest) []common.ErrorField {
	var errors []common.ErrorField
	if req.Id == 0 {
		errors = append(errors, common.ErrorField{Name: "id", Message: "아이디는 필수입니다."})
	}
	if len(req.Name) == 0 {
		errors = append(errors, common.ErrorField{Name: "name", Message: "이름은 필수입니다."})
	}
	return errors
}
