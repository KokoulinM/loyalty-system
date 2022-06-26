package auth

import (
	"reflect"
	"testing"

	"github.com/KokoulinM/go-musthave-diploma-tpl/internal/config"
)

func TestCreateToken(t *testing.T) {
	type args struct {
		userID string
		cfg    config.ConfigToken
	}
	tests := []struct {
		name    string
		args    args
		want    *TokenDetails
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateToken(tt.args.userID, tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}
