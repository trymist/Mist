package models

import (
	"time"
)

type AppRepositorySourceType string

const (
	SourceGitProvider AppRepositorySourceType = "git_provider"
	SourceGithubApp   AppRepositorySourceType = "github_app"
)

type AppRepositories struct {
	ID           int64                   `gorm:"primaryKey;autoIncrement:true" json:"id"`
	AppID        int64                   `gorm:"uniqueIndex:idx_app_repo_unique;not null;constraint:OnDelete:CASCADE" json:"app_id"`
	SourceType   AppRepositorySourceType `gorm:"not null" json:"source_type"`
	SourceID     int64                   `gorm:"not null" json:"source_id"`
	RepoFullName string                  `gorm:"uniqueIndex:idx_app_repo_unique;not null" json:"repo_full_name"`
	RepoURL      string                  `gorm:"not null" json:"repo_url"`
	Branch       string                  `gorm:"default:'main'" json:"branch"`
	WebhookID    int64                   `json:"webhook_id"`
	AutoDeploy   bool                    `gorm:"default:false" json:"auto_deploy"`
	LastSyncedAt *time.Time              `json:"last_synced_at,omitempty"`
}
