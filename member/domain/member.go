package domain

import (
	"github.com/golang-module/carbon/v2"
)

type Member struct {
	Id           int64 `gorm:"primarykey"`
	Name         string
	Email        string
	BirthDate    *carbon.Date    `gorm:"column:birthDate"`
	RegisterDate carbon.DateTime `gorm:"column:registerDate"`
}

func (m *Member) TableName() string {
	return "member"
}

// Update 고민이다! Member 필드는 어차피 public 인데....
func (m *Member) Update(name string, birthDate *carbon.Date) {
	m.Name = name
	m.BirthDate = birthDate
}
