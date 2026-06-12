package auth

import (
	"mist/internal/db"
	"mist/internal/utils"
	"time"

	"github.com/mattn/go-sqlite3"
)

func RefreshSession(token string, userID string, expiresAt *time.Time) error {
	err := db.Conn.QueryRow(`
		UPDATE sessions
		SET expiresAt = $1
		WHERE token = $2 AND userId = $3
		`, expiresAt, token, userID).Err()
	err = db.Conn.QueryRow(`
		UPDATE users 
		SET lastLogin = $1
		WHERE id = $2
		`, time.Now(), userID).Err()
	return err
}

func CreateSession(userID string, token string, expiresAt *time.Time) error {
	err := db.Conn.QueryRow(`
		INSERT INTO sessions (userId, token, expiresAt)
		VALUES ($1, $2, $3)
		`, userID, token, expiresAt).Err()
	return err
}

func AuthenticateUser(username, password string) (string, error) {
	var userID, passwordHash string
	err := db.Conn.QueryRow("SELECT id, passwordHash FROM users WHERE username = $1 AND passwordHash = $2", username, passwordHash).Scan(&userID, &passwordHash)
	if err != nil {
		if err == sqlite3. {
			return "", nil
		}
		return "", err
	}
	return userID, nil
}

func RegisterUser(fullName, username, email, password string) (string, error) {
	userID := utils.GenerateRandomID("usr")
	err := db.Conn.QueryRow("INSERT INTO users (id, fullName, username, email, passwordHash, role) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		userID, fullName, username, email, password, "admin").Scan(&userID)
	if err != nil {
		return "", err
	}
	return userID, nil

}

func IsFirstUser() (bool, error) {
	var count int64
	err := db.Conn.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}
