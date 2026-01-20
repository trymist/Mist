package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type AuditLog struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       *int64    `gorm:"index" json:"user_id"`
	Username     *string   `gorm:"-" json:"username"`
	Email        *string   `gorm:"-" json:"email"`
	Action       string    `gorm:"not null" json:"action"`
	ResourceType string    `gorm:"not null" json:"resourceType"`
	ResourceID   *int64    `json:"resourceId"`
	ResourceName *string   `gorm:"-" json:"resourceName"`
	Details      *string   `json:"details"`
	IPAddress    *string   `gorm:"-" json:"ipAddress"`
	UserAgent    *string   `gorm:"-" json:"userAgent"`
	TriggerType  string    `gorm:"-" json:"triggerType"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

type AuditLogDetails struct {
	Before interface{} `json:"before,omitempty"`
	After  interface{} `json:"after,omitempty"`
	Reason string      `json:"reason,omitempty"`
	Extra  interface{} `json:"extra,omitempty"`
}

func (a *AuditLog) Create() error {
	return db.Create(a).Error
}

//hlper function for join logic

func getAuditLogsQuery() *gorm.DB {
	return db.Table("audit_logs").Select(`audit_logs.id, 
			audit_logs.user_id, 
			users.username, 
			users.email, 
			audit_logs.action, 
			audit_logs.resource_type, 
			audit_logs.resource_id, 
			audit_logs.details, 
			audit_logs.created_at`).Joins("LEFT JOIN users ON audit_logs.user_id=users.id").Order("audit_logs.created_at DESC")

}

// helper function to reduce repetition
func logHelper(logs []AuditLog) []AuditLog {
	for i := range logs {
		log := &logs[i]
		if log.UserID == nil {
			log.TriggerType = "system"
			if log.Details != nil && (*log.Details != "") {
				var detailsMap map[string]interface{}
				if err := json.Unmarshal([]byte(*log.Details), &detailsMap); err == nil {
					if triggerType, ok := detailsMap["trigger_type"].(string); ok {
						log.TriggerType = triggerType
					}
				}
			}
		} else {
			log.TriggerType = "user"
		}
	}
	return logs
}

func GetAllAuditLogs(limit, offset int) ([]AuditLog, error) {
	var logs []AuditLog
	err := getAuditLogsQuery().Limit(limit).Offset(offset).Scan(&logs).Error
	if err != nil {
		return nil, err
	}
	return logHelper(logs), nil
}

func GetAuditLogsByUserID(userID int64, limit, offset int) ([]AuditLog, error) {
	var logs []AuditLog
	err := getAuditLogsQuery().Where("audit_logs.user_id = ?", userID).Limit(limit).Offset(offset).Scan(&logs).Error
	if err != nil {
		return nil, err
	}
	for i := range logs {
		logs[i].TriggerType = "user"
	}
	return logs, nil
}

func GetAuditLogsByResource(resourceType string, resourceID int64, limit, offset int) ([]AuditLog, error) {
	var logs []AuditLog
	err := getAuditLogsQuery().Where("audit_logs.resource_type = ? AND audit_logs.resource_id = ?", resourceType, resourceID).
		Limit(limit).Offset(offset).Scan(&logs).Error
	if err != nil {
		return nil, err
	}
	return logHelper(logs), nil
}

func GetAuditLogsCount() (int64, error) {
	var count int64
	err := db.Model(&AuditLog{}).Count(&count).Error
	return count, err
}

func LogAudit(userID *int64, action, resourceType string, resourceID *int64, details interface{}) error {
	var detailsJSON *string
	if details != nil {
		jsonBytes, err := json.Marshal(details)
		if err != nil {
			return err
		}
		jsonStr := string(jsonBytes)
		detailsJSON = &jsonStr
	}

	log := &AuditLog{
		UserID:       userID,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Details:      detailsJSON,
	}
	return log.Create()
}

func LogUserAudit(userID int64, action, resourceType string, resourceID *int64, details interface{}) error {
	return LogAudit(&userID, action, resourceType, resourceID, details)
}

func LogWebhookAudit(action, resourceType string, resourceID *int64, details map[string]interface{}) error {
	if details == nil {
		details = make(map[string]interface{})
	}
	details["trigger_type"] = "webhook"
	return LogAudit(nil, action, resourceType, resourceID, details)
}

func LogSystemAudit(action, resourceType string, resourceID *int64, details interface{}) error {
	detailsMap := make(map[string]interface{})
	detailsMap["trigger_type"] = "system"
	if details != nil {
		detailsMap["data"] = details
	}
	return LogAudit(nil, action, resourceType, resourceID, detailsMap)
}

func GetAuditLogsByResourceType(resourceType string, limit, offset int) ([]AuditLog, error) {
	var logs []AuditLog
	err := getAuditLogsQuery().Where("audit_logs.resource_type = ?", resourceType).
		Limit(limit).Offset(offset).Scan(&logs).Error
	if err != nil {
		return nil, err
	}
	return logHelper(logs), nil
}

func GetAuditLogByID(id int64) (*AuditLog, error) {
	var logs []AuditLog
	err := getAuditLogsQuery().
		Where("audit_logs.id = ?", id).
		Limit(1).
		Scan(&logs).Error

	if err != nil {
		return nil, err
	}
	if len(logs) == 0 {
		return nil, nil
	}

	enrichedLogs := logHelper(logs)
	return &enrichedLogs[0], nil
}
