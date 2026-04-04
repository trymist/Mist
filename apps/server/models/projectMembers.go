package models

import "time"

type ProjectMember struct {
	UserID int64 `gorm:"primaryKey;autoIncrement:false" json:"user_id"`

	ProjectID int64 `gorm:"primaryKey;autoIncrement:false" json:"project_id"`

	AddedAt time.Time `gorm:"autoCreateTime" json:"added_at"`
}

func (ProjectMember) TableName() string {
	return "project_members"
}
