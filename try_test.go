package grill

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/lovlin-thakkar/swiggy-grill/mock"
)

func TestTry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAssertion := mock.NewMockAssertion(ctrl)

	tests := []struct {
		name   string
		output [][]interface{}
	}{
		{"NoErrors", [][]interface{}{{nil, 3}}},
		{"ErrorInBetween", [][]interface{}{{nil, 1}, {fmt.Errorf("error"), 1}, {nil, 3}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var calls []*gomock.Call
			for _, o := range tt.output {
				calls = append(calls, mockAssertion.EXPECT().Assert().Return(o[0]).Times(o[1].(int)))
			}
			gomock.InOrder(calls...)
			Try(time.Second, 3, mockAssertion).Assert()
		})
	}
}
