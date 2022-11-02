package app

import (
	"github.com/stretchr/testify/assert"
	"go-sample/common"
	"testing"
)

func Test_validateMemberRegisterRequest(t *testing.T) {
	type args struct {
		req MemberRegistRequest
	}
	tests := []struct {
		name string
		args args
		want []common.ErrorField
	}{
		{
			name: "다 없음",
			args: args{req: MemberRegistRequest{}},
			want: []common.ErrorField{
				{Name: "name", Message: "이름은 필수입니다."},
				{Name: "email", Message: "이메일은 필수입니다."},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, validateMemberRegisterRequest(tt.args.req), "validateMemberRegisterRequest(%v)", tt.args.req)
		})
	}
}
