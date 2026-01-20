package models

import "time"

type Registry struct {
	ID int64 `gorm:"primaryKey;autoIncrement:true" json:"id"`

	ProjectID int64 `gorm:"uniqueIndex:idx_project_registry;not null;constraint:OnDelete:CASCADE" json:"projectId"`

	RegistryURL string `gorm:"uniqueIndex:idx_project_registry;not null" json:"registryUrl"`

	Username string `json:"username"`
	Password string `json:"password"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

func (Registry) TableName() string {
	return "registries"
}
