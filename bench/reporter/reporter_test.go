package reporter_test

import (
	"reflect"
	"testing"

	"github.com/isucon10-qualify/isucon10-qualify/bench/reporter"
)

func Test_uniqMsgs(t *testing.T) {
	type args struct {
		allMsgs []string
	}
	tests := []struct {
		name string
		args args
		want []reporter.Message
	}{
		{
			args: args{
				allMsgs: []string{"A", "B", "B"},
			},
			want: []reporter.Message{
				{Text: "A", Count: 1},
				{Text: "B", Count: 2},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := reporter.UniqMsgs(tt.args.allMsgs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("uniqMsgs() = %v, want %v", got, tt.want)
			}
		})
	}
}
