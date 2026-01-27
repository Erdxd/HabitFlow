package encrypto

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func Encrypto(password string) (string, error) {
	HasgedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Failed to encrypt your password")
	}
	return string(HasgedPassword), nil

}
func CheckPassword(HashedPassword, Password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(HashedPassword), []byte(Password))
	return err == nil

}
