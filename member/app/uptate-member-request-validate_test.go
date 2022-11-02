package app

import (
	"github.com/stretchr/testify/assert"
	"go-sample/common"
	"testing"
)

func TestValidateUpdateMemberRequest(t *testing.T) {
	type args struct {
		req UpdateMemberRequest
	}
	tests := []struct {
		name string
		args args
		want []common.ErrorField
	}{
		{
			name: "다 없음",
			args: args{req: UpdateMemberRequest{}},
			want: []common.ErrorField{
				{Name: "id", Message: "아이디는 필수입니다."},
				{Name: "name", Message: "이름은 필수입니다."},
			},
		},
		{
			name: "ID만 없음",
			args: args{req: UpdateMemberRequest{Name: "이름"}},
			want: []common.ErrorField{
				{Name: "id", Message: "아이디는 필수입니다."},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, validateUpdateMemberRequest(tt.args.req), "request = %v", tt.args.req)
		})
	}
}
