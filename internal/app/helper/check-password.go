package helper

import (
	"golang.org/x/crypto/bcrypt"
)

func CheckPassword(existingPassword string, providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(existingPassword), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}
