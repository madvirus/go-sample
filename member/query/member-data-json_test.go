package query

import (
	"encoding/json"
	"github.com/golang-module/carbon/v2"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshalMemberData(t *testing.T) {
	date := carbon.Date{
		Carbon: carbon.CreateFromDate(2022, 10, 19),
	}
	m := MemberData{
		BirthDate: &date,
	}
	b, err := json.Marshal(m)
	if err != nil {
		log.Errorf("error occured: %v", err)
		assert.FailNow(t, "error occured", err)
	}
	log.Infof("json: %s", string(b))
	assert.JSONEq(t, `{"Id":0, "Name": "", "Email":"", "BirthDate":"2022-10-19", "RegisterDate": null}`, string(b))
}
