package models

import (
	"time"

	"github.com/corecollectives/mist/utils"
	"gorm.io/gorm"
)

type ApiToken struct {
	ID          int64      `gorm:"primaryKey;autoIncrement:false" json:"id"`
	UserID      int64      `gorm:"index;not null" json:"userId"`
	Name        string     `gorm:"type:varchar(255);not null" json:"name"`
	TokenHash   string     `gorm:"uniqueIndex;type:varchar(255);not null" json:"-"`
	TokenPrefix string     `gorm:"index;type:varchar(50);not null" json:"tokenPrefix"`
	Scopes      *string    `json:"scopes,omitempty"`
	LastUsedAt  *time.Time `json:"lastUsedAt,omitempty"`
	LastUsedIP  *string    `json:"lastUsedIp,omitempty"`
	UsageCount  int        `gorm:"default:0" json:"usageCount"`
	ExpiresAt   *time.Time `gorm:"index" json:"expiresAt,omitempty"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	RevokedAt   *time.Time `json:"revokedAt,omitempty"`
}

func (t *ApiToken) ToJson() map[string]interface{} {
	return map[string]interface{}{
		"id":          t.ID,
		"userId":      t.UserID,
		"name":        t.Name,
		"tokenPrefix": t.TokenPrefix,
		"scopes":      t.Scopes,
		"lastUsedAt":  t.LastUsedAt,
		"lastUsedIp":  t.LastUsedIP,
		"usageCount":  t.UsageCount,
		"expiresAt":   t.ExpiresAt,
		"createdAt":   t.CreatedAt,
		"revokedAt":   t.RevokedAt,
	}
}

func (t *ApiToken) InsertInDB() error {
	t.ID = utils.GenerateRandomId()
	result := db.Create(t)
	return result.Error
}

func GetApiTokensByUserID(userID int64) ([]ApiToken, error) {
	var tokens []ApiToken
	result := db.Where("user_id=? AND revoked_at IS NULL", userID).Order("created_at DESC").Find(&tokens)
	return tokens, result.Error
}

func GetApiTokenByHash(tokenHash string) (*ApiToken, error) {
	var token ApiToken
	result := db.Where("token_hash=? AND revoked_at IS NULL", tokenHash).First(&token)
	if result.Error != nil {
		return nil, result.Error
	}
	return &token, nil
}

func (t *ApiToken) UpdateUsage(lastUsedIP string) error {
	return db.Model(t).Updates(map[string]interface{}{
		"last_used_at": time.Now(),
		"last_used_ip": lastUsedIP,
		"usage_count":  gorm.Expr("usage_count + ?", 1),
	}).Error
}

func (t *ApiToken) Revoke() error {
	return db.Model(t).Update("revoked_at", time.Now()).Error
}

func DeleteApiToken(tokenID int64) error {
	return db.Delete(&ApiToken{}, tokenID).Error
}
