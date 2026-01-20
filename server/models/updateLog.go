package models

import (
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type UpdateStatus string

const (
	UpdateStatusInProgress UpdateStatus = "in_progress"
	UpdateStatusSuccess    UpdateStatus = "success"
	UpdateStatusFailed     UpdateStatus = "failed"
)

type UpdateLog struct {
	ID           int64        `gorm:"primaryKey;autoIncrement:true" json:"id"`
	VersionFrom  string       `gorm:"not null" json:"version_from"`
	VersionTo    string       `gorm:"not null" json:"version_to"`
	Status       UpdateStatus `gorm:"index;not null" json:"status"`
	Logs         *string      `json:"logs"`
	ErrorMessage *string      `json:"error_message"`
	StartedBy    int64        `gorm:"not null;constraint:OnDelete:CASCADE" json:"started_by"`
	StartedAt    time.Time    `gorm:"autoCreateTime;index:,sort:desc" json:"started_at"`
	CompletedAt  *time.Time   `json:"completed_at"`
	Username     *string      `gorm:"-" json:"username"`
}

func CreateUpdateLog(versionFrom, versionTo string, startedBy int64) (*UpdateLog, error) {
	emptyLogs := ""
	updateLog := &UpdateLog{
		VersionFrom: versionFrom,
		VersionTo:   versionTo,
		Status:      UpdateStatusInProgress,
		Logs:        &emptyLogs,
		StartedBy:   startedBy,
	}

	if err := db.Create(updateLog).Error; err != nil {
		log.Error().Err(err).Msg("Failed to create update log")
		return nil, err
	}

	log.Info().
		Int64("update_log_id", updateLog.ID).
		Str("from", versionFrom).
		Str("to", versionTo).
		Int64("started_by", startedBy).
		Msg("Update log created")

	return updateLog, nil
}

func UpdateUpdateLogStatus(id int64, status UpdateStatus, logs string, errorMessage *string) error {
	err := db.Model(&UpdateLog{ID: id}).Updates(map[string]interface{}{
		"status":        status,
		"logs":          logs,
		"error_message": errorMessage,
		"completed_at":  time.Now(),
	}).Error

	if err != nil {
		log.Error().Err(err).Int64("update_log_id", id).Msg("Failed to update log status")
		return err
	}

	log.Info().
		Int64("update_log_id", id).
		Str("status", string(status)).
		Msg("Update log status updated")

	return nil
}

func AppendUpdateLog(id int64, logLine string) error {
	err := db.Model(&UpdateLog{ID: id}).
		Update("logs", gorm.Expr("logs || ?", logLine+"\n")).Error

	if err != nil {
		log.Error().Err(err).Int64("update_log_id", id).Msg("Failed to append log line")
		return err
	}

	return nil
}

func GetUpdateLogs(limit int) ([]UpdateLog, error) {
	var logs []UpdateLog

	err := db.Table("update_logs").
		Select("update_logs.*, users.username").
		Joins("LEFT JOIN users ON update_logs.started_by = users.id").
		Order("update_logs.started_at DESC").
		Limit(limit).
		Scan(&logs).Error

	if err != nil {
		log.Error().Err(err).Msg("Failed to query update logs")
		return nil, err
	}

	return logs, nil
}

func GetUpdateLogByID(id int64) (*UpdateLog, error) {
	var updateLog UpdateLog

	err := db.Table("update_logs").
		Select("update_logs.*, users.username").
		Joins("LEFT JOIN users ON update_logs.started_by = users.id").
		Where("update_logs.id = ?", id).
		Scan(&updateLog).Error

	if err != nil {
		log.Error().Err(err).Int64("update_log_id", id).Msg("Failed to get update log by ID")
		return nil, err
	}

	if updateLog.ID == 0 {
		return nil, nil
	}

	return &updateLog, nil
}

func GetUpdateLogsAsString() (string, error) {
	logs, err := GetUpdateLogs(10)
	if err != nil {
		return "", err
	}

	if len(logs) == 0 {
		return "No update history available", nil
	}

	var builder strings.Builder
	builder.WriteString("Recent Update History:\n")
	builder.WriteString("======================\n\n")

	for _, log := range logs {
		builder.WriteString("Version: ")
		builder.WriteString(log.VersionFrom)
		builder.WriteString(" → ")
		builder.WriteString(log.VersionTo)
		builder.WriteString("\n")
		builder.WriteString("Status: ")
		builder.WriteString(string(log.Status))
		builder.WriteString("\n")
		builder.WriteString("Started: ")
		builder.WriteString(log.StartedAt.Format("2006-01-02 15:04:05"))
		builder.WriteString(" by ")
		if log.Username != nil {
			builder.WriteString(*log.Username)
		} else {
			builder.WriteString("unknown")
		}

		builder.WriteString("\n")
		if log.CompletedAt != nil {
			builder.WriteString("Completed: ")
			builder.WriteString(log.CompletedAt.Format("2006-01-02 15:04:05"))
			builder.WriteString("\n")
		}
		if log.ErrorMessage != nil && *log.ErrorMessage != "" {
			builder.WriteString("Error: ")
			builder.WriteString(*log.ErrorMessage)
			builder.WriteString("\n")
		}
		builder.WriteString("\n")
	}

	return builder.String(), nil
}

func CheckAndCompletePendingUpdates() error {
	logs, err := GetUpdateLogs(1)
	if err != nil {
		return err
	}

	if len(logs) == 0 {
		return nil
	}

	latestLog := logs[0]

	if latestLog.Status != UpdateStatusInProgress {
		return nil
	}

	log.Info().
		Int64("update_log_id", latestLog.ID).
		Str("from_version", latestLog.VersionFrom).
		Str("to_version", latestLog.VersionTo).
		Str("age", time.Since(latestLog.StartedAt).String()).
		Msg("Found in-progress update on startup, checking status")

	currentVersion, err := GetSystemSetting("version")
	if err != nil {
		log.Error().Err(err).Msg("Failed to get current version for update completion check")
		return err
	}

	if currentVersion == "" {
		currentVersion = "1.0.0"
	}

	if currentVersion == latestLog.VersionTo {
		log.Info().
			Int64("update_log_id", latestLog.ID).
			Str("version", currentVersion).
			Msg("Completing successful update that was interrupted by service restart")

		existing := ""
		if latestLog.Logs != nil {
			existing = *latestLog.Logs
		}

		completionLog := existing + "\n✅ Update completed successfully (verified on restart)\n"

		err = UpdateUpdateLogStatus(latestLog.ID, UpdateStatusSuccess, completionLog, nil)
		if err != nil {
			log.Error().Err(err).Int64("update_log_id", latestLog.ID).Msg("Failed to complete pending update")
			return err
		}

		log.Info().
			Int64("update_log_id", latestLog.ID).
			Str("from", latestLog.VersionFrom).
			Str("to", latestLog.VersionTo).
			Msg("Successfully completed pending update")
		return nil
	}

	log.Warn().
		Int64("update_log_id", latestLog.ID).
		Str("expected_version", latestLog.VersionTo).
		Str("current_version", currentVersion).
		Str("age", time.Since(latestLog.StartedAt).String()).
		Msg("Update appears to have failed (version mismatch detected on startup)")

	errMsg := "Update process was interrupted and version does not match target"
	existing := ""
	if latestLog.Logs != nil {
		existing = *latestLog.Logs
	}

	failureLog := existing + "\n❌ " + errMsg + "\n"
	err = UpdateUpdateLogStatus(latestLog.ID, UpdateStatusFailed, failureLog, &errMsg)
	if err != nil {
		log.Error().Err(err).Int64("update_log_id", latestLog.ID).Msg("Failed to mark failed update")
		return err
	}

	log.Info().
		Int64("update_log_id", latestLog.ID).
		Msg("Marked failed update as failed")

	return nil
}

//############################################################################################################################
//ARCHIVED CODE BELOW

// package models

// import (
// 	"database/sql"
// 	"strings"
// 	"time"

// 	"github.com/rs/zerolog/log"
// )

// type UpdateStatus string

// const (
// 	UpdateStatusInProgress UpdateStatus = "in_progress"
// 	UpdateStatusSuccess    UpdateStatus = "success"
// 	UpdateStatusFailed     UpdateStatus = "failed"
// )

// type UpdateLog struct {
// 	ID           int64
// 	VersionFrom  string
// 	VersionTo    string
// 	Status       UpdateStatus
// 	Logs         *string
// 	ErrorMessage *string
// 	StartedBy    int64
// 	StartedAt    time.Time
// 	CompletedAt  *time.Time
// 	Username     *string
// }

// func CreateUpdateLog(versionFrom, versionTo string, startedBy int64) (*UpdateLog, error) {
// 	query := `
// 		INSERT INTO update_logs (version_from, version_to, status, logs, started_by, started_at)
// 		VALUES (?, ?, 'in_progress', '', ?, CURRENT_TIMESTAMP)
// 		RETURNING id, version_from, version_to, status, logs, error_message, started_by, started_at, completed_at
// 	`

// 	updateLog := &UpdateLog{}
// 	err := db.QueryRow(query, versionFrom, versionTo, startedBy).Scan(
// 		&updateLog.ID,
// 		&updateLog.VersionFrom,
// 		&updateLog.VersionTo,
// 		&updateLog.Status,
// 		&updateLog.Logs,
// 		&updateLog.ErrorMessage,
// 		&updateLog.StartedBy,
// 		&updateLog.StartedAt,
// 		&updateLog.CompletedAt,
// 	)

// 	if err != nil {
// 		log.Error().Err(err).Msg("Failed to create update log")
// 		return nil, err
// 	}

// 	log.Info().
// 		Int64("update_log_id", updateLog.ID).
// 		Str("from", versionFrom).
// 		Str("to", versionTo).
// 		Int64("started_by", startedBy).
// 		Msg("Update log created")

// 	return updateLog, nil
// }

// func UpdateUpdateLogStatus(id int64, status UpdateStatus, logs string, errorMessage *string) error {
// 	query := `
// 		UPDATE update_logs
// 		SET status = ?, logs = ?, error_message = ?, completed_at = CURRENT_TIMESTAMP
// 		WHERE id = ?
// 	`

// 	_, err := db.Exec(query, status, logs, errorMessage, id)
// 	if err != nil {
// 		log.Error().Err(err).Int64("update_log_id", id).Msg("Failed to update log status")
// 		return err
// 	}

// 	log.Info().
// 		Int64("update_log_id", id).
// 		Str("status", string(status)).
// 		Msg("Update log status updated")

// 	return nil
// }

// func AppendUpdateLog(id int64, logLine string) error {
// 	query := `
// 		UPDATE update_logs
// 		SET logs = logs || ?
// 		WHERE id = ?
// 	`

// 	_, err := db.Exec(query, logLine+"\n", id)
// 	if err != nil {
// 		log.Error().Err(err).Int64("update_log_id", id).Msg("Failed to append log line")
// 		return err
// 	}

// 	return nil
// }

// func GetUpdateLogs(limit int) ([]UpdateLog, error) {
// 	query := `
// 		SELECT
// 			ul.id, ul.version_from, ul.version_to, ul.status,
// 			ul.logs, ul.error_message, ul.started_by, ul.started_at,
// 			ul.completed_at, u.username
// 		FROM update_logs ul
// 		LEFT JOIN users u ON ul.started_by = u.id
// 		ORDER BY ul.started_at DESC
// 		LIMIT ?
// 	`

// 	rows, err := db.Query(query, limit)
// 	if err != nil {
// 		log.Error().Err(err).Msg("Failed to query update logs")
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var logs []UpdateLog
// 	for rows.Next() {
// 		var updateLog UpdateLog
// 		err := rows.Scan(
// 			&updateLog.ID,
// 			&updateLog.VersionFrom,
// 			&updateLog.VersionTo,
// 			&updateLog.Status,
// 			&updateLog.Logs,
// 			&updateLog.ErrorMessage,
// 			&updateLog.StartedBy,
// 			&updateLog.StartedAt,
// 			&updateLog.CompletedAt,
// 			&updateLog.Username,
// 		)
// 		if err != nil {
// 			log.Error().Err(err).Msg("Failed to scan update log row")
// 			return nil, err
// 		}
// 		logs = append(logs, updateLog)
// 	}

// 	return logs, nil
// }

// func GetUpdateLogByID(id int64) (*UpdateLog, error) {
// 	query := `
// 		SELECT
// 			ul.id, ul.version_from, ul.version_to, ul.status,
// 			ul.logs, ul.error_message, ul.started_by, ul.started_at,
// 			ul.completed_at, u.username
// 		FROM update_logs ul
// 		LEFT JOIN users u ON ul.started_by = u.id
// 		WHERE ul.id = ?
// 	`

// 	updateLog := &UpdateLog{}
// 	err := db.QueryRow(query, id).Scan(
// 		&updateLog.ID,
// 		&updateLog.VersionFrom,
// 		&updateLog.VersionTo,
// 		&updateLog.Status,
// 		&updateLog.Logs,
// 		&updateLog.ErrorMessage,
// 		&updateLog.StartedBy,
// 		&updateLog.StartedAt,
// 		&updateLog.CompletedAt,
// 		&updateLog.Username,
// 	)

// 	if err == sql.ErrNoRows {
// 		return nil, nil
// 	}

// 	if err != nil {
// 		log.Error().Err(err).Int64("update_log_id", id).Msg("Failed to get update log by ID")
// 		return nil, err
// 	}

// 	return updateLog, nil
// }

// func GetUpdateLogsAsString() (string, error) {
// 	logs, err := GetUpdateLogs(10)
// 	if err != nil {
// 		return "", err
// 	}

// 	if len(logs) == 0 {
// 		return "No update history available", nil
// 	}

// 	var builder strings.Builder
// 	builder.WriteString("Recent Update History:\n")
// 	builder.WriteString("======================\n\n")

// 	for _, log := range logs {
// 		builder.WriteString("Version: ")
// 		builder.WriteString(log.VersionFrom)
// 		builder.WriteString(" → ")
// 		builder.WriteString(log.VersionTo)
// 		builder.WriteString("\n")
// 		builder.WriteString("Status: ")
// 		builder.WriteString(string(log.Status))
// 		builder.WriteString("\n")
// 		builder.WriteString("Started: ")
// 		builder.WriteString(log.StartedAt.Format("2006-01-02 15:04:05"))
// 		builder.WriteString(" by ")
// 		if log.Username != nil {
// 			builder.WriteString(*log.Username)
// 		} else {
// 			builder.WriteString("unknown")
// 		}

// 		builder.WriteString("\n")
// 		if log.CompletedAt != nil {
// 			builder.WriteString("Completed: ")
// 			builder.WriteString(log.CompletedAt.Format("2006-01-02 15:04:05"))
// 			builder.WriteString("\n")
// 		}
// 		if log.ErrorMessage != nil && *log.ErrorMessage != "" {
// 			builder.WriteString("Error: ")
// 			builder.WriteString(*log.ErrorMessage)
// 			builder.WriteString("\n")
// 		}
// 		builder.WriteString("\n")
// 	}

// 	return builder.String(), nil
// }

// func CheckAndCompletePendingUpdates() error {
// 	logs, err := GetUpdateLogs(1)
// 	if err != nil {
// 		return err
// 	}

// 	if len(logs) == 0 {
// 		return nil
// 	}

// 	latestLog := logs[0]

// 	if latestLog.Status != UpdateStatusInProgress {
// 		return nil
// 	}

// 	log.Info().
// 		Int64("update_log_id", latestLog.ID).
// 		Str("from_version", latestLog.VersionFrom).
// 		Str("to_version", latestLog.VersionTo).
// 		Str("age", time.Since(latestLog.StartedAt).String()).
// 		Msg("Found in-progress update on startup, checking status")

// 	currentVersion, err := GetSystemSetting("version")
// 	if err != nil {
// 		log.Error().Err(err).Msg("Failed to get current version for update completion check")
// 		return err
// 	}

// 	if currentVersion == "" {
// 		currentVersion = "1.0.0"
// 	}

// 	if currentVersion == latestLog.VersionTo {
// 		log.Info().
// 			Int64("update_log_id", latestLog.ID).
// 			Str("version", currentVersion).
// 			Msg("Completing successful update that was interrupted by service restart")

// 		existing := ""
// 		if latestLog.Logs != nil {
// 			existing = *latestLog.Logs
// 		}

// 		completionLog := existing + "\n✅ Update completed successfully (verified on restart)\n"

// 		err = UpdateUpdateLogStatus(latestLog.ID, UpdateStatusSuccess, completionLog, nil)
// 		if err != nil {
// 			log.Error().Err(err).Int64("update_log_id", latestLog.ID).Msg("Failed to complete pending update")
// 			return err
// 		}

// 		log.Info().
// 			Int64("update_log_id", latestLog.ID).
// 			Str("from", latestLog.VersionFrom).
// 			Str("to", latestLog.VersionTo).
// 			Msg("Successfully completed pending update")
// 		return nil
// 	}

// 	log.Warn().
// 		Int64("update_log_id", latestLog.ID).
// 		Str("expected_version", latestLog.VersionTo).
// 		Str("current_version", currentVersion).
// 		Str("age", time.Since(latestLog.StartedAt).String()).
// 		Msg("Update appears to have failed (version mismatch detected on startup)")

// 	errMsg := "Update process was interrupted and version does not match target"
// 	existing := ""
// 	if latestLog.Logs != nil {
// 		existing = *latestLog.Logs
// 	}

// 	failureLog := existing + "\n❌ " + errMsg + "\n"
// 	err = UpdateUpdateLogStatus(latestLog.ID, UpdateStatusFailed, failureLog, &errMsg)
// 	if err != nil {
// 		log.Error().Err(err).Int64("update_log_id", latestLog.ID).Msg("Failed to mark failed update")
// 		return err
// 	}

// 	log.Info().
// 		Int64("update_log_id", latestLog.ID).
// 		Msg("Marked failed update as failed")

// 	return nil
// }
