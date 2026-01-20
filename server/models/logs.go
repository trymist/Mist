package models

import "time"

type LogSource string

const (
	LogSourceApp    LogSource = "app"
	LogSourceSystem LogSource = "system"
)

type LogLevel string

const (
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
	LogLevelDebug LogLevel = "debug"
)

type Logs struct {
	ID        int64     `gorm:"primaryKey;autoIncrement:true"`
	Source    LogSource `gorm:"not null"`
	SourceID  *int64    `gorm:"index"`
	Message   string    `gorm:"not null"`
	Level     LogLevel  `gorm:"not null;default:'info'"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
