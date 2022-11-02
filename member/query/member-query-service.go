package query

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type MemberQueryService interface {
	GetMembers() ([]MemberData, error)
	GetMember(id int64) (*MemberData, error)
}

func CreateMemberQueryService(db *gorm.DB) MemberQueryService {
	return &memberQueryService{
		db: db,
	}
}

type memberQueryService struct {
	db *gorm.DB
}

func (m *memberQueryService) GetMembers() ([]MemberData, error) {
	var members []MemberData
	result := m.db.Find(&members)
	if result.Error != nil {
		return nil, result.Error
	}
	return members, nil
}

func (m *memberQueryService) GetMember(id int64) (*MemberData, error) {
	var member MemberData
	result := m.db.Take(&member, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Infof("MemberData[id=%d] not found", id)
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &member, nil
}
