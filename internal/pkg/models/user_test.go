package models

import "testing"

func TestUser_GeneratePasswordHashAndComparePassword(t *testing.T) {
	tests := []struct {
		name                                string
		passwordForGenerateHash             string
		passwordToCompareWithHashedPassword string
		wantInCompare                       bool
		wantErrInGenerate                   bool
	}{
		{name: "OK", passwordForGenerateHash: "qwerty", passwordToCompareWithHashedPassword: "qwerty", wantInCompare: true},
		{name: "Unmatched password", passwordForGenerateHash: "qwerty", passwordToCompareWithHashedPassword: "12345", wantInCompare: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := User{}

			if err := u.GeneratePasswordHash(tt.passwordForGenerateHash); err != nil {
				t.Errorf("GeneratePasswordHash() = %v, want %v", err, tt.wantErrInGenerate)
			}
			if got := u.ComparePassword(tt.passwordToCompareWithHashedPassword); got != tt.wantInCompare {
				t.Errorf("ComparePassword() = %v, want %v", got, tt.wantInCompare)
			}
		})
	}
}
