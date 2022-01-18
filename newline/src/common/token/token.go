package internal

import (
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

// TODO change this to newline.com own signature
var webTokenSignature = []byte("")

// GenerateJWTToken Generate Json Web Token
func GenerateJWTToken(d map[string]interface{}) (string, error) {
	mc := jwt.MapClaims{}
	for k, v := range d {
		mc[k] = v
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mc)
	return token.SignedString(webTokenSignature)
}

// ParseJWTToken Parse Json Web Token
func ParseJWTToken(encryptedToken string) (map[string]interface{}, error) {
	r := map[string]interface{}{}
	if strings.Count(encryptedToken, ".") != 2 {
		return r, fmt.Errorf("Unexpected encryptedToken: %s", encryptedToken)
	}
	token, err := jwt.Parse(encryptedToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return webTokenSignature, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		for k, v := range claims {
			r[k] = v
		}
	}
	return r, err
}
