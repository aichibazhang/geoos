package filter

import (
	"reflect"
	"testing"

	"github.com/spatial-go/geoos/algorithm/matrix"
)

func TestUniqueArrayFilter_FilterSteric(t *testing.T) {

	type args struct {
		matr matrix.Steric
	}
	tests := []struct {
		name string
		want []interface{}
		args args
	}{
		{name: "unique array filter",
			args: args{matrix.LineMatrix{{0, 0}, {1, 1}, {2, 1}, {3, 1}, {3, 0}, {0, 0}}},
			want: []interface{}{matrix.Matrix{0, 0}, matrix.Matrix{1, 1}, matrix.Matrix{2, 1}, matrix.Matrix{3, 1}, matrix.Matrix{3, 0}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UniqueArrayFilter{}
			entities := matrix.TransMatrixes(tt.args.matr)
			for _, v := range entities {
				u.Filter(v)
			}
			if got := u.Entities(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unique array filter = %v, want %v", got, tt.want)
			}
		})
	}
}
