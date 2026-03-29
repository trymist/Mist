package queue

import (
	"gorm.io/gorm"
)

func InitQueue(db *gorm.DB) *Queue {
	q := NewQueue(5, db)
	return q
}
