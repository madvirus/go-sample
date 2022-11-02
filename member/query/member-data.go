package query

import "github.com/golang-module/carbon/v2"

type MemberData struct {
	Id           int64 `gorm:"primarykey"`
	Name         string
	Email        string
	BirthDate    *carbon.Date     `gorm:"column:birthDate"`
	RegisterDate *carbon.DateTime `gorm:"column:registerDate"`
}

func (d *MemberData) TableName() string {
	return "member"
}
