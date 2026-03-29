package models

import (
	"time"

	"github.com/corecollectives/mist/utils"
)

type NotificationType string

const (
	NotificationDeploymentSuccess NotificationType = "deployment_success"
	NotificationDeploymentFailed  NotificationType = "deployment_failed"
	NotificationDeploymentStarted NotificationType = "deployment_started"
	NotificationSSLExpiryWarning  NotificationType = "ssl_expiry_warning"
	NotificationSSLRenewalSuccess NotificationType = "ssl_renewal_success"
	NotificationSSLRenewalFailed  NotificationType = "ssl_renewal_failed"
	NotificationResourceAlert     NotificationType = "resource_alert"
	NotificationAppError          NotificationType = "app_error"
	NotificationAppStopped        NotificationType = "app_stopped"
	NotificationBackupSuccess     NotificationType = "backup_success"
	NotificationBackupFailed      NotificationType = "backup_failed"
	NotificationUserInvited       NotificationType = "user_invited"
	NotificationMemberAdded       NotificationType = "member_added"
	NotificationSystemUpdate      NotificationType = "system_update"
	NotificationCustom            NotificationType = "custom"
)

type NotificationPriority string

const (
	PriorityLow    NotificationPriority = "low"
	PriorityNormal NotificationPriority = "normal"
	PriorityHigh   NotificationPriority = "high"
	PriorityUrgent NotificationPriority = "urgent"
)

type Notification struct {
	ID int64 `gorm:"primaryKey;autoIncrement:false" json:"id"`

	UserID *int64 `gorm:"index;constraint:OnDelete:CASCADE" json:"userId,omitempty"`

	Type    NotificationType `gorm:"index;not null" json:"type"`
	Title   string           `gorm:"not null" json:"title"`
	Message string           `gorm:"not null" json:"message"`

	Link         *string `json:"link,omitempty"`
	ResourceType *string `json:"resourceType,omitempty"`
	ResourceID   *int64  `json:"resourceId,omitempty"`

	EmailSent   bool       `gorm:"default:false" json:"emailSent"`
	EmailSentAt *time.Time `json:"emailSentAt,omitempty"`

	SlackSent   bool       `gorm:"default:false" json:"slackSent"`
	SlackSentAt *time.Time `json:"slackSentAt,omitempty"`

	DiscordSent   bool       `gorm:"default:false" json:"discordSent"`
	DiscordSentAt *time.Time `json:"discordSentAt,omitempty"`

	WebhookSent   bool       `gorm:"default:false" json:"webhookSent"`
	WebhookSentAt *time.Time `json:"webhookSentAt,omitempty"`

	IsRead bool       `gorm:"default:false;index" json:"isRead"`
	ReadAt *time.Time `json:"readAt,omitempty"`

	Priority NotificationPriority `gorm:"default:'normal';index" json:"priority"`

	Metadata *string `json:"metadata,omitempty"`

	CreatedAt time.Time  `gorm:"autoCreateTime;index:,sort:desc" json:"createdAt"`
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
}

func (n *Notification) ToJson() map[string]interface{} {
	return map[string]interface{}{
		"id":            n.ID,
		"userId":        n.UserID,
		"type":          n.Type,
		"title":         n.Title,
		"message":       n.Message,
		"link":          n.Link,
		"resourceType":  n.ResourceType,
		"resourceId":    n.ResourceID,
		"emailSent":     n.EmailSent,
		"emailSentAt":   n.EmailSentAt,
		"slackSent":     n.SlackSent,
		"slackSentAt":   n.SlackSentAt,
		"discordSent":   n.DiscordSent,
		"discordSentAt": n.DiscordSentAt,
		"webhookSent":   n.WebhookSent,
		"webhookSentAt": n.WebhookSentAt,
		"isRead":        n.IsRead,
		"readAt":        n.ReadAt,
		"priority":      n.Priority,
		"metadata":      n.Metadata,
		"createdAt":     n.CreatedAt,
		"expiresAt":     n.ExpiresAt,
	}
}

func (n *Notification) InsertInDB() error {
	n.ID = utils.GenerateRandomId()

	if n.Priority == "" {
		n.Priority = PriorityNormal
	}

	return db.Create(n).Error
}

func GetNotificationsByUserID(userID int64, unreadOnly bool) ([]Notification, error) {
	var notifications []Notification

	query := db.Where(db.Where("user_id = ?", userID).Or("user_id IS NULL"))

	if unreadOnly {
		query = query.Where("is_read = ?", false)
	}

	result := query.Order("created_at DESC").Limit(100).Find(&notifications)

	return notifications, result.Error
}

func (n *Notification) MarkAsRead() error {
	return db.Model(n).Updates(map[string]interface{}{
		"is_read": true,
		"read_at": time.Now(),
	}).Error
}

func MarkAllAsRead(userID int64) error {
	return db.Model(&Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": time.Now(),
		}).Error
}

func DeleteNotification(notificationID int64) error {
	return db.Delete(&Notification{}, notificationID).Error
}

func DeleteExpiredNotifications() error {
	return db.Where("expires_at IS NOT NULL AND expires_at < ?", time.Now()).
		Delete(&Notification{}).Error
}

//######################################################################################################
//ARCHIVED CODE BELOW

// package models

// import (
// 	"time"

// 	"github.com/corecollectives/mist/utils"
// )

// type NotificationType string

// const (
// 	NotificationDeploymentSuccess NotificationType = "deployment_success"
// 	NotificationDeploymentFailed  NotificationType = "deployment_failed"
// 	NotificationDeploymentStarted NotificationType = "deployment_started"
// 	NotificationSSLExpiryWarning  NotificationType = "ssl_expiry_warning"
// 	NotificationSSLRenewalSuccess NotificationType = "ssl_renewal_success"
// 	NotificationSSLRenewalFailed  NotificationType = "ssl_renewal_failed"
// 	NotificationResourceAlert     NotificationType = "resource_alert"
// 	NotificationAppError          NotificationType = "app_error"
// 	NotificationAppStopped        NotificationType = "app_stopped"
// 	NotificationBackupSuccess     NotificationType = "backup_success"
// 	NotificationBackupFailed      NotificationType = "backup_failed"
// 	NotificationUserInvited       NotificationType = "user_invited"
// 	NotificationMemberAdded       NotificationType = "member_added"
// 	NotificationSystemUpdate      NotificationType = "system_update"
// 	NotificationCustom            NotificationType = "custom"
// )

// type NotificationPriority string

// const (
// 	PriorityLow    NotificationPriority = "low"
// 	PriorityNormal NotificationPriority = "normal"
// 	PriorityHigh   NotificationPriority = "high"
// 	PriorityUrgent NotificationPriority = "urgent"
// )

// type Notification struct {
// 	ID            int64                `db:"id" json:"id"`
// 	UserID        *int64               `db:"user_id" json:"userId,omitempty"`
// 	Type          NotificationType     `db:"type" json:"type"`
// 	Title         string               `db:"title" json:"title"`
// 	Message       string               `db:"message" json:"message"`
// 	Link          *string              `db:"link" json:"link,omitempty"`
// 	ResourceType  *string              `db:"resource_type" json:"resourceType,omitempty"`
// 	ResourceID    *int64               `db:"resource_id" json:"resourceId,omitempty"`
// 	EmailSent     bool                 `db:"email_sent" json:"emailSent"`
// 	EmailSentAt   *time.Time           `db:"email_sent_at" json:"emailSentAt,omitempty"`
// 	SlackSent     bool                 `db:"slack_sent" json:"slackSent"`
// 	SlackSentAt   *time.Time           `db:"slack_sent_at" json:"slackSentAt,omitempty"`
// 	DiscordSent   bool                 `db:"discord_sent" json:"discordSent"`
// 	DiscordSentAt *time.Time           `db:"discord_sent_at" json:"discordSentAt,omitempty"`
// 	WebhookSent   bool                 `db:"webhook_sent" json:"webhookSent"`
// 	WebhookSentAt *time.Time           `db:"webhook_sent_at" json:"webhookSentAt,omitempty"`
// 	IsRead        bool                 `db:"is_read" json:"isRead"`
// 	ReadAt        *time.Time           `db:"read_at" json:"readAt,omitempty"`
// 	Priority      NotificationPriority `db:"priority" json:"priority"`
// 	Metadata      *string              `db:"metadata" json:"metadata,omitempty"` // JSON
// 	CreatedAt     time.Time            `db:"created_at" json:"createdAt"`
// 	ExpiresAt     *time.Time           `db:"expires_at" json:"expiresAt,omitempty"`
// }

// func (n *Notification) ToJson() map[string]interface{} {
// 	return map[string]interface{}{
// 		"id":            n.ID,
// 		"userId":        n.UserID,
// 		"type":          n.Type,
// 		"title":         n.Title,
// 		"message":       n.Message,
// 		"link":          n.Link,
// 		"resourceType":  n.ResourceType,
// 		"resourceId":    n.ResourceID,
// 		"emailSent":     n.EmailSent,
// 		"emailSentAt":   n.EmailSentAt,
// 		"slackSent":     n.SlackSent,
// 		"slackSentAt":   n.SlackSentAt,
// 		"discordSent":   n.DiscordSent,
// 		"discordSentAt": n.DiscordSentAt,
// 		"webhookSent":   n.WebhookSent,
// 		"webhookSentAt": n.WebhookSentAt,
// 		"isRead":        n.IsRead,
// 		"readAt":        n.ReadAt,
// 		"priority":      n.Priority,
// 		"metadata":      n.Metadata,
// 		"createdAt":     n.CreatedAt,
// 		"expiresAt":     n.ExpiresAt,
// 	}
// }

// func (n *Notification) InsertInDB() error {
// 	id := utils.GenerateRandomId()
// 	n.ID = id
// 	query := `
// 	INSERT INTO notifications (
// 		id, user_id, type, title, message, link,
// 		resource_type, resource_id, priority, metadata, expires_at
// 	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
// 	RETURNING created_at
// 	`
// 	err := db.QueryRow(query, n.ID, n.UserID, n.Type, n.Title, n.Message, n.Link,
// 		n.ResourceType, n.ResourceID, n.Priority, n.Metadata, n.ExpiresAt).Scan(&n.CreatedAt)
// 	return err
// }

// func GetNotificationsByUserID(userID int64, unreadOnly bool) ([]Notification, error) {
// 	var notifications []Notification
// 	query := `
// 	SELECT id, user_id, type, title, message, link, resource_type, resource_id,
// 	       email_sent, email_sent_at, slack_sent, slack_sent_at,
// 	       discord_sent, discord_sent_at, webhook_sent, webhook_sent_at,
// 	       is_read, read_at, priority, metadata, created_at, expires_at
// 	FROM notifications
// 	WHERE user_id = ? OR user_id IS NULL
// 	`
// 	if unreadOnly {
// 		query += " AND is_read = 0"
// 	}
// 	query += " ORDER BY created_at DESC LIMIT 100"

// 	rows, err := db.Query(query, userID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var notif Notification
// 		err := rows.Scan(
// 			&notif.ID, &notif.UserID, &notif.Type, &notif.Title, &notif.Message,
// 			&notif.Link, &notif.ResourceType, &notif.ResourceID,
// 			&notif.EmailSent, &notif.EmailSentAt, &notif.SlackSent, &notif.SlackSentAt,
// 			&notif.DiscordSent, &notif.DiscordSentAt, &notif.WebhookSent, &notif.WebhookSentAt,
// 			&notif.IsRead, &notif.ReadAt, &notif.Priority, &notif.Metadata,
// 			&notif.CreatedAt, &notif.ExpiresAt,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}
// 		notifications = append(notifications, notif)
// 	}

// 	return notifications, rows.Err()
// }

// func (n *Notification) MarkAsRead() error {
// 	query := `
// 	UPDATE notifications
// 	SET is_read = 1, read_at = CURRENT_TIMESTAMP
// 	WHERE id = ?
// 	`
// 	_, err := db.Exec(query, n.ID)
// 	return err
// }

// func MarkAllAsRead(userID int64) error {
// 	query := `
// 	UPDATE notifications
// 	SET is_read = 1, read_at = CURRENT_TIMESTAMP
// 	WHERE user_id = ? AND is_read = 0
// 	`
// 	_, err := db.Exec(query, userID)
// 	return err
// }

// func DeleteNotification(notificationID int64) error {
// 	query := `DELETE FROM notifications WHERE id = ?`
// 	_, err := db.Exec(query, notificationID)
// 	return err
// }

// func DeleteExpiredNotifications() error {
// 	query := `
// 	DELETE FROM notifications
// 	WHERE expires_at IS NOT NULL AND expires_at < CURRENT_TIMESTAMP
// 	`
// 	_, err := db.Exec(query)
// 	return err
// }
