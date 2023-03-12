package auth

import (
	// "log"

	// "remood/models"

	"golang.org/x/crypto/bcrypt"
)

// func ValidateUsername(user *models.User, username string) error {
// 	if err := user.GetOne("username", username); err != nil {
// 		return err
// 	}
// 	return nil
// }

func ValidatePassword(hashedPassword string, password string) error {
	// log.Println(user, password)
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return err
	}
	return nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(hashedPassword), err
}
