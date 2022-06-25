package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidLuhnNumber(t *testing.T) {
	type args struct {
		number int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "positive case",
			args: args{
				number: 79927398713,
			},
			want: true,
		},
		{
			name: "negative case",
			args: args{
				number: 123,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := ValidLuhnNumber(tt.args.number)

			assert.Equal(t, tt.want, res)
		})
	}
}
