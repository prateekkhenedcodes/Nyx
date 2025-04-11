package auth

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
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

func TestValidateJWT(t *testing.T) {
	userId := uuid.New().String()
	token, _ := MakeJWT(userId, "secret", time.Hour)

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserId  string
		wantErr     bool
	}{
		{
			"Correct tokenString", token, "secret", userId, false,
		},
		{
			"Incorrect tokenString", "xyz", "secret", "", true,
		},
		{
			"Incorrect tokenSecret", token, "xyz", "", true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserId, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUserId != tt.wantUserId {
				t.Errorf("ValidateJWT() gotUserId: %v, wantUserId: %v", gotUserId, tt.wantUserId)
			}
		})
	}
}

func TestStripAuth(t *testing.T) {
	authHeader1 := "Bearer TOKEN_STRING"
	authHeader2 := "Bearer "
	token1, _ := stripBearer(authHeader1)
	token2, _ := stripBearer(authHeader2)

	tests := []struct {
		name    string
		auth    string
		token   string
		wantErr bool
	}{
		{
			"Correct token", authHeader1, token1, false,
		},
		{
			"Correct token", authHeader2, token2, false,
		},
		{
			"Incorrect auth", "Bearer", token1, true,
		},
		{
			"Incorrect token", authHeader1, "TOKEN", true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := stripBearer(tt.auth)
			if !strings.Contains(tt.auth, token) {
				if (err != nil) != tt.wantErr {
					t.Errorf("stripBearer() error: %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}
