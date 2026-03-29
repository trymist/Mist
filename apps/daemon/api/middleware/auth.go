package middleware

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/corecollectives/mist/api/handlers"
	"github.com/corecollectives/mist/models"
	"github.com/golang-jwt/jwt"
	"github.com/rs/zerolog/log"
)

var jwtSecret []byte

func InitJWTSecret() error {
	settings, err := models.GetSystemSettings()
	if err != nil {
		return err
	}
	jwtSecret = []byte(settings.JwtSecret)
	log.Info().Msg("JWT secret loaded from database")
	return nil
}

func GenerateJWT(userID int64, email, role string) (string, error) {
	if len(jwtSecret) == 0 {
		if err := InitJWTSecret(); err != nil {
			return "", err
		}
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(31 * 24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

type JWTClaims struct {
	UserID int64
	Email  string
	Role   string
}

func VerifyJWT(tokenStr string) (*JWTClaims, error) {
	if len(jwtSecret) == 0 {
		if err := InitJWTSecret(); err != nil {
			return nil, err
		}
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			if int64(exp) < time.Now().Unix() {
				return nil, errors.New("token expired")
			}
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			return nil, errors.New("invalid user_id in token")
		}

		email, ok := claims["email"].(string)
		if !ok {
			return nil, errors.New("invalid email in token")
		}

		role, ok := claims["role"].(string)
		if !ok {
			return nil, errors.New("invalid role in token")
		}

		return &JWTClaims{
			UserID: int64(userIDFloat),
			Email:  email,
			Role:   role,
		}, nil
	}

	return nil, errors.New("invalid token")
}

type contextKey string

const userContextKey = contextKey("user-data")

func AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("mist_token")
			if err != nil {
				handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Authentication required", "No token provided")
				return
			}

			tokenString := cookie.Value
			claims, err := VerifyJWT(tokenString)

			if err != nil {
				handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "Invalid or expired token", err.Error())
				return
			}

			user, err := models.GetUserByID(claims.UserID)
			if err != nil {
				handlers.SendResponse(w, http.StatusInternalServerError, false, nil, "Failed to retrieve user", err.Error())
				return
			}

			if user == nil {
				handlers.SendResponse(w, http.StatusUnauthorized, false, nil, "User not found", "Invalid token user")
				return
			}

			ctx := context.WithValue(r.Context(), userContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUser(r *http.Request) (*models.User, bool) {
	user, ok := r.Context().Value(userContextKey).(*models.User)
	return user, ok
}
