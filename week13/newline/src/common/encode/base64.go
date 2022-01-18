package encode

import (
	"encoding/base64"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// EncodePassword EncodePassword
func EncodePassword(password string) string {
	bcryptHashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logrus.Error(err)
	}
	base64BcryptHashedPassword := base64.StdEncoding.EncodeToString(bcryptHashedPassword)
	return base64BcryptHashedPassword
}

// CompareEncodedPassword CompareEncodedPassword
func CompareEncodedPassword(encodedPassword, password string) bool {
	bcryptHashedPassword, _ := base64.StdEncoding.DecodeString(encodedPassword)
	err := bcrypt.CompareHashAndPassword(bcryptHashedPassword, []byte(password))
	return err == nil
}
