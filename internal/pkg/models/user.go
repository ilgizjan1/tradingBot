package models

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID            int    `json:"-" db:"id"`
	Name          string `json:"name" binding:"required"`
	Username      string `json:"username" binding:"required"`
	Password      string `json:"password" binding:"required" db:"password_hash"`
	PublicAPIKey  string `json:"public_api_key" binding:"required" db:"public_api_key"`
	PrivateAPIKey string `json:"private_api_key" binding:"required" db:"private_api_key"`
}

func (u *User) GeneratePasswordHash(password string) error {
	byteHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return err
	}
	u.Password = string(byteHash)
	return nil
}

func (u *User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
