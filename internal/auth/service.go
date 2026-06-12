package auth

import (
	"database/sql"
	"fmt"
	"mist/internal/db"
	"mist/internal/utils"
	"time"
)

func RefreshSession(token string, userID string, expiresAt *time.Time) error {
	_, err := db.Conn.Exec(`
		UPDATE sessions
		SET expiresAt = ?, token = ?
		WHERE userId = ?
	`, expiresAt, token, userID)
	if err != nil {
		return err
	}

	_, err = db.Conn.Exec(`
		UPDATE users
		SET lastLogin = ?
		WHERE id = ?
	`, time.Now(), userID)

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
	err := db.Conn.QueryRow("SELECT id, passwordHash FROM users WHERE username = $1", username).Scan(&userID, &passwordHash)
	if err != nil {
		fmt.Println("Error querying user:", err)
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("invalid username or password")
		}
		return "", err
	}
	fmt.Println("User found, checking password")
	fmt.Println("Password hash from DB:", passwordHash)
	passwordMatch := utils.CheckPasswordHash(password, passwordHash)
	if !passwordMatch {
		return "", fmt.Errorf("invalid username or password")
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
