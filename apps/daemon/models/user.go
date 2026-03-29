package models

import (
	"time"

	"github.com/corecollectives/mist/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID           int64  `gorm:"primaryKey;autoIncrement:false" json:"id"`
	Username     string `gorm:"uniqueIndex;not null" json:"username"`
	Email        string `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string `gorm:"not null" json:"-"`
	Role         string `gorm:"default:'user';index" json:"role"`

	FullName  *string `json:"fullName,omitempty"`
	AvatarURL *string `json:"avatarUrl,omitempty"`
	Bio       *string `json:"bio,omitempty"`

	EmailVerified           bool       `gorm:"default:false" json:"emailVerified"`
	EmailVerificationToken  *string    `json:"-"`
	EmailVerificationSentAt *time.Time `json:"-"`

	PasswordResetToken     *string    `json:"-"`
	PasswordResetExpiresAt *time.Time `json:"-"`
	PasswordChangedAt      *time.Time `json:"-"`

	TwoFactorEnabled     bool    `gorm:"default:false" json:"twoFactorEnabled"`
	TwoFactorSecret      *string `json:"-"`
	TwoFactorBackupCodes *string `json:"-"`

	LastLoginAt         *time.Time `json:"lastLoginAt,omitempty"`
	LastLoginIP         *string    `json:"lastLoginIp,omitempty"`
	FailedLoginAttempts int        `gorm:"default:0" json:"failedLoginAttempts"`
	AccountLockedUntil  *time.Time `json:"accountLockedUntil,omitempty"`

	Timezone                string  `gorm:"default:'UTC'" json:"timezone"`
	Language                string  `gorm:"default:'en'" json:"language"`
	NotificationPreferences *string `json:"notificationPreferences,omitempty"`

	IsActive  bool           `gorm:"default:true;index" json:"isActive"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

func (u *User) Create() error {
	u.ID = utils.GenerateRandomId()
	if u.Role == "" {
		u.Role = "user"
	}
	return db.Create(u).Error
}

func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hashedPassword)
	return nil
}

func GetUserByID(userID int64) (*User, error) {
	var user User
	err := db.First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func DeleteUserByID(userID int64) error {
	return db.Delete(&User{}, userID).Error
}

func UpdateUser(u *User) error {
	return db.Model(u).Updates(map[string]interface{}{
		"username":   u.Username,
		"email":      u.Email,
		"role":       u.Role,
		"avatar_url": u.AvatarURL,
		"updated_at": time.Now(),
	}).Error
}

func (u *User) MatchPassword(password string) bool {
	var storedHash string
	err := db.Model(&User{}).
		Select("password_hash").
		Where("id = ?", u.ID).
		Scan(&storedHash).Error

	if err != nil {
		return false
	}

	return bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password)) == nil
}

func (u *User) UpdatePassword() error {
	return db.Model(u).Updates(map[string]interface{}{
		"password_hash": u.PasswordHash,
		"updated_at":    time.Now(),
	}).Error
}

func GetUserByEmail(email string) (*User, error) {
	var user User
	err := db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByUsername(username string) (*User, error) {
	var user User
	err := db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateUserPassword(userID int64, passwordHash string) error {
	return db.Model(&User{ID: userID}).Updates(map[string]interface{}{
		"password_hash": passwordHash,
		"updated_at":    time.Now(),
	}).Error
}

func GetAllUsers() ([]User, error) {
	var users []User
	err := db.Find(&users).Error
	return users, err
}

func GetUserRole(userID int64) (string, error) {
	var role string
	err := db.Model(&User{}).
		Select("role").
		Where("id = ?", userID).
		Scan(&role).Error

	if err != nil {
		return "", err
	}
	return role, nil
}

func GetUserCount() (int64, error) {
	var count int64
	err := db.Model(&User{}).Count(&count).Error
	return count, err
}

//##############################################################################################################
//ARCHIVED CODE BELOW

// package models

// import (
// 	"time"

// 	"github.com/corecollectives/mist/utils"
// 	"golang.org/x/crypto/bcrypt"
// )

// type User struct {
// 	ID           int64     `json:"id"`
// 	Username     string    `json:"username"`
// 	Email        string    `json:"email"`
// 	PasswordHash string    `json:"-"`
// 	Role         string    `json:"role"`
// 	AvatarURL    *string   `json:"avatarUrl"`
// 	CreatedAt    time.Time `json:"createdAt"`
// 	UpdatedAt    time.Time `json:"updatedAt"`
// }

// func (u *User) Create() error {
// 	query := `
// 		INSERT INTO users (id, username, email, password_hash, role, avatar_url)
// 		VALUES ($1, $2, $3, $4, $5, $6)
// 		RETURNING id, username, email, role, avatar_url, created_at, updated_at
// 	`
// 	u.ID = utils.GenerateRandomId()
// 	err := db.QueryRow(query, u.ID, u.Username, u.Email, u.PasswordHash, u.Role, u.AvatarURL).Scan(
// 		&u.ID, &u.Username, &u.Email, &u.Role, &u.AvatarURL, &u.CreatedAt, &u.UpdatedAt,
// 	)
// 	return err
// }

// func (u *User) SetPassword(password string) error {
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return err
// 	}
// 	u.PasswordHash = string(hashedPassword)
// 	return nil
// }

// func GetUserByID(userID int64) (*User, error) {
// 	query := `
// 	  SELECT id, username, email, role, avatar_url, created_at, updated_at
// 	  FROM users
// 	  WHERE id = $1
// 	`
// 	user := &User{}
// 	err := db.QueryRow(query, userID).Scan(
// 		&user.ID, &user.Username, &user.Email, &user.Role, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return user, nil
// }

// func DeleteUserByID(userID int64) error {
// 	query := `DELETE FROM users WHERE id = $1`
// 	_, err := db.Exec(query, userID)
// 	return err
// }

// func UpdateUser(u *User) error {
// 	query := `
// 		UPDATE users
// 		SET username = $1, email = $2,  role = $3, avatar_url = $4, updated_at = CURRENT_TIMESTAMP
// 		WHERE id = $5
// 		RETURNING updated_at
// 	`
// 	return db.QueryRow(query, u.Username, u.Email, u.Role, u.AvatarURL, u.ID).Scan(&u.UpdatedAt)
// }

// func (u *User) MatchPassword(password string) bool {
// 	query := `
// 		SELECT password_hash
// 		FROM users
// 		WHERE id = $1
// 	`
// 	var storedHash string
// 	err := db.QueryRow(query, u.ID).Scan(&storedHash)
// 	if err != nil {
// 		return false
// 	}

// 	return bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password)) == nil
// }

// func (u *User) UpdatePassword() error {
// 	query := `
// 		UPDATE users
// 		SET password_hash = $1, updated_at = CURRENT_TIMESTAMP
// 		WHERE id = $2
// 		RETURNING updated_at
//   `
// 	return db.QueryRow(query, u.PasswordHash, u.ID).Scan(&u.UpdatedAt)

// }

// func GetUserByEmail(email string) (*User, error) {
// 	query := `
// 	  SELECT id, username, email, role, avatar_url, created_at, updated_at
// 	  FROM users
// 	  WHERE email = $1
// 	`
// 	user := &User{}
// 	err := db.QueryRow(query, email).Scan(
// 		&user.ID, &user.Username, &user.Email, &user.Role, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return user, nil
// }

// func GetUserByUsername(username string) (*User, error) {
// 	query := `
// 	  SELECT id, username, email, role, avatar_url, created_at, updated_at
// 	  FROM users
// 	  WHERE username = $1
// 	`
// 	user := &User{}
// 	err := db.QueryRow(query, username).Scan(
// 		&user.ID, &user.Username, &user.Email, &user.Role, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt,
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return user, nil
// }

// func UpdateUserPassword(userID int64, passwordHash string) error {
// 	query := `
// 		UPDATE users
// 		SET password_hash = $1, updated_at = CURRENT_TIMESTAMP
// 		WHERE id = $2
// 	`
// 	_, err := db.Exec(query, passwordHash, userID)
// 	return err
// }

// func GetAllUsers() ([]User, error) {
// 	query := `
// 		SELECT id, username, email, role, avatar_url, created_at, updated_at
// 		FROM users
// 	`
// 	rows, err := db.Query(query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var users []User
// 	for rows.Next() {
// 		user := User{}
// 		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.AvatarURL, &user.CreatedAt, &user.UpdatedAt)
// 		if err != nil {
// 			return nil, err
// 		}
// 		users = append(users, user)
// 	}
// 	return users, nil
// }

// func GetUserRole(userID int64) (string, error) {
// 	query := `
// 		SELECT role
// 		FROM users
// 		WHERE id = $1
// 	`
// 	var role string
// 	err := db.QueryRow(query, userID).Scan(&role)
// 	if err != nil {
// 		return "", err
// 	}
// 	return role, nil
// }

// func GetUserCount() (int, error) {
// 	query := `SELECT COUNT(*) FROM users`
// 	var count int
// 	err := db.QueryRow(query).Scan(&count)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return count, nil
// }
