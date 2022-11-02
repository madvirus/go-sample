//go:build integration

package test_integration

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func CreateTestDbHelper(db *gorm.DB) *TestDbHelper {
	return &TestDbHelper{db: db}
}

type TestDbHelper struct {
	db *gorm.DB
}

func (h *TestDbHelper) Truncate() {
	tables := []string{
		"member",
	}
	for _, name := range tables {
		h.db.Exec("truncate table " + name)
		log.Infof("truncated table %s", name)
	}
}

func (h *TestDbHelper) Save(value any) {
	h.db.Save(value)
}
