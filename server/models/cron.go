package models

import "time"

type Cron struct {
	ID        int64     `gorm:"primaryKey;autoIncrement:true" json:"id"`
	AppID     int64     `gorm:"index;constraint:OnDelete:CASCADE;not null" json:"app_id"`
	Name      string    `gorm:"index;not null" json:"name"`
	Schedule  string    `gorm:"not null" json:"schedule"`
	Command   string    `gorm:"not null" json:"command"`
	LastRun   *time.Time `gorm:"type:timestamp" json:"last_run"`
	NextRun   *time.Time `gorm:"type:timestamp" json:"next_run"`
	Enable    bool      `gorm:"default:true" json:"enable"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
