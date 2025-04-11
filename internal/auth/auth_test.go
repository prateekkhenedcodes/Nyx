package auth

import (
	"testing"
)

func TestCheckPassHash(t *testing.T) {
	pass1 := "examplepass1"
	pass2 := "examplepass2"
	hashPass1, _ := HashNyxCode(pass1)
	hashPass2, _ := HashNyxCode(pass2)

	tests := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{
			"Correct password", pass1, hashPass1, false,
		},
		{
			"Correct password", pass2, hashPass2, false,
		},
		{
			"password doesn't match hash", pass1, hashPass2, true,
		},
		{
			"password doesn't match hash", pass2, hashPass1, true,
		},
		{
			"Empty password", "", hashPass1, true,
		},
		{
			"incorrect password", "wrongpass", hashPass1, true,
		},
		{
			"incorrect hash", pass1, "wronghash", true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPassHash(tt.hash, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPassHash() error: %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
