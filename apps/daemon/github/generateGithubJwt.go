package github

import (
	"encoding/pem"
	"time"

	"github.com/corecollectives/mist/models"
	"github.com/golang-jwt/jwt"
)

func GenerateGithubJwt(appID int) (string, error) {
	appNumericId, appPrivateKeyPEM, err := models.GetGithubAppCredentials(appID)
	if err != nil {
		return "", err
	}

	block, _ := pem.Decode([]byte(appPrivateKeyPEM))
	if block == nil {
		return "", err
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(appPrivateKeyPEM))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(10 * time.Minute).Unix(),
		"iss": appNumericId,
	})

	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil

}
