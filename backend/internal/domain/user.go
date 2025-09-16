package domain

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int32  `json:"id"`
	Username       string `json:"username"`
	Name           string `json:"name"`
	MeetLink       string `json:"-"`
	HashedPassword []byte `json:"-"`
}

func HashPassword(plainPassword string) ([]byte, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %v", err)
	}
	return h, nil
}

func ValidatePassword(password string, hashedPassword []byte) bool {
	return bcrypt.CompareHashAndPassword(hashedPassword, []byte(password)) == nil
}
