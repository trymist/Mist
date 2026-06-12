package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"mist/internal/db"
	"time"
)

var jwtsecret string

func Init() {
	err := db.Conn.QueryRow("SELECT value FROM system_settings WHERE key = 'jwtSecret'").Scan(&jwtsecret)
	if err != nil {
		jwtsecret = GenerateRandomID("jwt")
	}
}

func SignJWT(userID string) (string, error) {
	claims := jwt.MapClaims{
		"userID": userID,
		"exp":    jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		"iat":    jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtsecret))
}

func VerifyJWT(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtsecret, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["userID"].(string)
		if !ok {
			return "", jwt.ErrInvalidKey
		}
		return userID, nil
	} else {
		return "", jwt.ErrInvalidKey
	}

}
