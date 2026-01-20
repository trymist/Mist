package models

import (
	"fmt"
	"time"

	"github.com/corecollectives/mist/utils"
	"gorm.io/gorm"
)

type DeploymentStatus string

const (
	DeploymentStatusPending    DeploymentStatus = "pending"
	DeploymentStatusBuilding   DeploymentStatus = "building"
	DeploymentStatusDeploying  DeploymentStatus = "deploying"
	DeploymentStatusSuccess    DeploymentStatus = "success"
	DeploymentStatusFailed     DeploymentStatus = "failed"
	DeploymentStatusStopped    DeploymentStatus = "stopped"
	DeploymentStatusRolledBack DeploymentStatus = "rolled_back"
)

type Deployment struct {
	ID int64 `gorm:"primaryKey;autoIncrement:false" json:"id"`

	AppID int64 `gorm:"index:idx_deployments_app_id;not null;constraint:OnDelete:CASCADE" json:"app_id"`

	CommitHash string `gorm:"index:idx_deployments_commit_hash;not null" json:"commit_hash"`

	CommitMessage *string `json:"commit_message,omitempty"`
	CommitAuthor  *string `json:"commit_author,omitempty"`

	TriggeredBy *int64 `gorm:"constraint:OnDelete:SET NULL" json:"triggered_by,omitempty"`

	DeploymentNumber *int `json:"deployment_number,omitempty"`

	ContainerID   *string `json:"container_id,omitempty"`
	ContainerName *string `json:"container_name,omitempty"`
	ImageTag      *string `json:"image_tag,omitempty"`

	Logs          *string `json:"logs,omitempty"`
	BuildLogsPath *string `json:"build_logs_path,omitempty"`

	Status DeploymentStatus `gorm:"default:'pending';index:idx_deployments_status" json:"status"`

	Stage    string `gorm:"default:'pending'" json:"stage"`
	Progress int    `gorm:"default:0" json:"progress"`

	ErrorMessage *string `json:"error_message,omitempty"`

	CreatedAt time.Time `gorm:"autoCreateTime;index:idx_deployments_created_at,sort:desc" json:"created_at"`

	StartedAt  *time.Time `json:"started_at,omitempty"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
	Duration   *int       `json:"duration,omitempty"`

	IsActive bool `gorm:"default:false;index:idx_deployments_is_active" json:"is_active"`

	RolledBackFrom *int64 `gorm:"constraint:OnDelete:SET NULL" json:"rolled_back_from,omitempty"`
}

func (d *Deployment) ToJson() map[string]interface{} {
	return map[string]interface{}{
		"id":               d.ID,
		"appId":            d.AppID,
		"commitHash":       d.CommitHash,
		"commitMessage":    d.CommitMessage,
		"commitAuthor":     d.CommitAuthor,
		"triggeredBy":      d.TriggeredBy,
		"deploymentNumber": d.DeploymentNumber,
		"containerId":      d.ContainerID,
		"containerName":    d.ContainerName,
		"imageTag":         d.ImageTag,
		"logs":             d.Logs,
		"buildLogsPath":    d.BuildLogsPath,
		"status":           d.Status,
		"stage":            d.Stage,
		"progress":         d.Progress,
		"errorMessage":     d.ErrorMessage,
		"createdAt":        d.CreatedAt,
		"startedAt":        d.StartedAt,
		"finishedAt":       d.FinishedAt,
		"duration":         d.Duration,
		"isActive":         d.IsActive,
		"rolledBackFrom":   d.RolledBackFrom,
	}
}

func GetDeploymentsByAppID(appID int64) ([]Deployment, error) {
	var deployments []Deployment
	result := db.Where("app_id = ?", appID).Order("created_at DESC").Find(&deployments)
	return deployments, result.Error
}

func (d *Deployment) CreateDeployment() error {
	d.ID = utils.GenerateRandomId()
	var maxDepNum *int
	result := db.Model(&Deployment{}).
		Where("app_id = ?", d.AppID).
		Pluck("MAX(deployment_number)", &maxDepNum)

	currentNum := 0
	if result.Error == nil && maxDepNum != nil {
		currentNum = *maxDepNum
	}
	newNum := currentNum + 1
	d.DeploymentNumber = &newNum

	d.Status = DeploymentStatusPending
	d.Stage = "pending"
	d.Progress = 0
	d.IsActive = false

	return db.Create(d).Error
}

func GetDeploymentByID(depID int64) (*Deployment, error) {
	var deployment Deployment
	result := db.First(&deployment, "id = ?", depID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &deployment, nil
}

func GetDeploymentByCommitHash(commitHash string) (*Deployment, error) {
	var deployment Deployment
	result := db.First(&deployment, "commit_hash = ?", commitHash)
	if result.Error != nil {
		return nil, result.Error
	}
	return &deployment, nil
}

func GetDeploymentByAppIDAndCommitHash(appID int64, commitHash string) (*Deployment, error) {
	var deployment Deployment
	result := db.First(&deployment, "app_id = ? AND commit_hash = ?", appID, commitHash)
	if result.Error != nil {
		return nil, result.Error
	}
	return &deployment, nil
}
func GetActiveDeploymentByAppID(appID int64) (*Deployment, error) {
	var deployment Deployment
	err := db.Where("app_id = ? AND is_active = ?", appID, true).First(&deployment)
	if err.Error != nil {
		return nil, err.Error
	}
	return &deployment, nil
}

func GetCommitHashByDeploymentID(depID int64) (string, error) {
	var d Deployment
	result := db.Select("commit_hash").First(&d, "id = ?", depID)
	if result.Error != nil {
		return "", result.Error
	}
	return d.CommitHash, nil
}

func UpdateDeploymentStatus(depID int64, status, stage string, progress int, errorMsg *string) error {
	var d Deployment
	result := db.First(&d, "id = ?", depID)
	if result.Error != nil {
		return result.Error
	}

	updates := map[string]interface{}{
		"status":        status,
		"stage":         stage,
		"progress":      progress,
		"error_message": errorMsg,
	}
	if status == string(DeploymentStatusFailed) || status == string(DeploymentStatusSuccess) || status == string(DeploymentStatusStopped) {
		now := time.Now()
		updates["finished_at"] = &now
		if d.StartedAt != nil {
			duration := int(now.Sub(*d.StartedAt).Seconds())
			updates["duration"] = &duration
		}
	}
	if errorMsg != nil {
		fmt.Println("updated dep status: ", *errorMsg)
	}
	return db.Model(d).Updates(updates).Error
}
func MarkDeploymentStarted(depID int64) error {
	now := time.Now()
	return db.Model(&Deployment{}).Where("id = ?", depID).Update("started_at", now).Error
}

func MarkDeploymentActive(depID int64, appID int64) error {
	return db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&Deployment{}).Where("app_id=?", appID).Update("is_active", false).Error
		if err != nil {
			return err
		}
		err = tx.Model(&Deployment{}).Where("id = ?", depID).Update("is_active", true).Error
		if err != nil {
			return err
		}
		return nil
	})
}
func UpdateContainerInfo(depID int64, containerID, containerName, imageTag string) error {
	updates := map[string]interface{}{
		"container_id":   containerID,
		"container_name": containerName,
		"image_tag":      imageTag,
	}
	return db.Model(&Deployment{}).Where("id = ?", depID).Updates(updates).Error
}

func GetDeploymentStatus(depID int64) (string, error) {
	var status string
	result := db.Model(&Deployment{}).Select("status").Where("id = ?", depID).Scan(&status)
	if result.Error != nil {
		return "", result.Error
	}

	if result.RowsAffected == 0 {
		return "", gorm.ErrRecordNotFound
	}

	return status, nil
}

//#############################################################################################################
//ARCHIVED CODE BELOW------>

// package models

// import (
// 	"time"

// 	"github.com/corecollectives/mist/utils"
// )

// type DeploymentStatus string

// const (
// 	DeploymentStatusPending    DeploymentStatus = "pending"
// 	DeploymentStatusBuilding   DeploymentStatus = "building"
// 	DeploymentStatusDeploying  DeploymentStatus = "deploying"
// 	DeploymentStatusSuccess    DeploymentStatus = "success"
// 	DeploymentStatusFailed     DeploymentStatus = "failed"
// 	DeploymentStatusStopped    DeploymentStatus = "stopped"
// 	DeploymentStatusRolledBack DeploymentStatus = "rolled_back"
// )

// type Deployment struct {
// 	ID               int64            `db:"id" json:"id"`
// 	AppID            int64            `db:"app_id" json:"app_id"`
// 	CommitHash       string           `db:"commit_hash" json:"commit_hash"`
// 	CommitMessage    *string          `db:"commit_message" json:"commit_message,omitempty"`
// 	CommitAuthor     *string          `db:"commit_author" json:"commit_author,omitempty"`
// 	TriggeredBy      *int64           `db:"triggered_by" json:"triggered_by,omitempty"`
// 	DeploymentNumber *int             `db:"deployment_number" json:"deployment_number,omitempty"`
// 	ContainerID      *string          `db:"container_id" json:"container_id,omitempty"`
// 	ContainerName    *string          `db:"container_name" json:"container_name,omitempty"`
// 	ImageTag         *string          `db:"image_tag" json:"image_tag,omitempty"`
// 	Logs             *string          `db:"logs" json:"logs,omitempty"`
// 	BuildLogsPath    *string          `db:"build_logs_path" json:"build_logs_path,omitempty"`
// 	Status           DeploymentStatus `db:"status" json:"status"`
// 	Stage            string           `db:"stage" json:"stage"`
// 	Progress         int              `db:"progress" json:"progress"`
// 	ErrorMessage     *string          `db:"error_message" json:"error_message,omitempty"`
// 	CreatedAt        time.Time        `db:"created_at" json:"created_at"`
// 	StartedAt        *time.Time       `db:"started_at" json:"started_at,omitempty"`
// 	FinishedAt       *time.Time       `db:"finished_at" json:"finished_at,omitempty"`
// 	Duration         *int             `db:"duration" json:"duration,omitempty"`
// 	IsActive         bool             `db:"is_active" json:"is_active"`
// 	RolledBackFrom   *int64           `db:"rolled_back_from" json:"rolled_back_from,omitempty"`
// }

// func (d *Deployment) ToJson() map[string]interface{} {
// 	return map[string]interface{}{
// 		"id":               d.ID,
// 		"appId":            d.AppID,
// 		"commitHash":       d.CommitHash,
// 		"commitMessage":    d.CommitMessage,
// 		"commitAuthor":     d.CommitAuthor,
// 		"triggeredBy":      d.TriggeredBy,
// 		"deploymentNumber": d.DeploymentNumber,
// 		"containerId":      d.ContainerID,
// 		"containerName":    d.ContainerName,
// 		"imageTag":         d.ImageTag,
// 		"logs":             d.Logs,
// 		"buildLogsPath":    d.BuildLogsPath,
// 		"status":           d.Status,
// 		"stage":            d.Stage,
// 		"progress":         d.Progress,
// 		"errorMessage":     d.ErrorMessage,
// 		"createdAt":        d.CreatedAt,
// 		"startedAt":        d.StartedAt,
// 		"finishedAt":       d.FinishedAt,
// 		"duration":         d.Duration,
// 		"isActive":         d.IsActive,
// 		"rolledBackFrom":   d.RolledBackFrom,
// 	}
// }

// func GetDeploymentsByAppID(appID int64) ([]Deployment, error) {
// 	query := `
// 	SELECT id, app_id, commit_hash, commit_message, commit_author, triggered_by,
// 	       deployment_number, container_id, container_name, image_tag,
// 	       logs, build_logs_path, status, stage, progress, error_message,
// 	       created_at, started_at, finished_at, duration, is_active, rolled_back_from
// 	FROM deployments
// 	WHERE app_id = ?
// 	ORDER BY created_at DESC
// 	`

// 	rows, err := db.Query(query, appID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var deployments []Deployment
// 	for rows.Next() {
// 		var d Deployment
// 		if err := rows.Scan(
// 			&d.ID, &d.AppID, &d.CommitHash, &d.CommitMessage, &d.CommitAuthor,
// 			&d.TriggeredBy, &d.DeploymentNumber, &d.ContainerID, &d.ContainerName,
// 			&d.ImageTag, &d.Logs, &d.BuildLogsPath, &d.Status, &d.Stage, &d.Progress,
// 			&d.ErrorMessage, &d.CreatedAt, &d.StartedAt, &d.FinishedAt, &d.Duration,
// 			&d.IsActive, &d.RolledBackFrom,
// 		); err != nil {
// 			return nil, err
// 		}
// 		deployments = append(deployments, d)
// 	}

// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}

// 	return deployments, nil
// }

// func (d *Deployment) CreateDeployment() error {
// 	id := utils.GenerateRandomId()
// 	d.ID = id

// 	var maxDeploymentNum int
// 	err := db.QueryRow(`SELECT COALESCE(MAX(deployment_number), 0) FROM deployments WHERE app_id = ?`, d.AppID).Scan(&maxDeploymentNum)
// 	if err == nil {
// 		deploymentNum := maxDeploymentNum + 1
// 		d.DeploymentNumber = &deploymentNum
// 	}

// 	query := `
// 	INSERT INTO deployments (
// 		id, app_id, commit_hash, commit_message, commit_author, triggered_by,
// 		deployment_number, status, stage, progress
// 	) VALUES (?, ?, ?, ?, ?, ?, ?, 'pending', 'pending', 0)
// 	RETURNING created_at
// 	`
// 	err = db.QueryRow(query, d.ID, d.AppID, d.CommitHash, d.CommitMessage,
// 		d.CommitAuthor, d.TriggeredBy, d.DeploymentNumber).Scan(&d.CreatedAt)
// 	if err != nil {
// 		return err
// 	}

// 	d.Status = DeploymentStatusPending
// 	d.Stage = "pending"
// 	d.Progress = 0
// 	d.IsActive = false
// 	return nil
// }

// func GetDeploymentByID(depID int64) (*Deployment, error) {
// 	query := `
// 	SELECT id, app_id, commit_hash, commit_message, commit_author, triggered_by,
// 	       deployment_number, container_id, container_name, image_tag,
// 	       logs, build_logs_path, status, stage, progress, error_message,
// 	       created_at, started_at, finished_at, duration, is_active, rolled_back_from
// 	FROM deployments
// 	WHERE id = ?
// 	`

// 	var d Deployment
// 	if err := db.QueryRow(query, depID).Scan(
// 		&d.ID, &d.AppID, &d.CommitHash, &d.CommitMessage, &d.CommitAuthor,
// 		&d.TriggeredBy, &d.DeploymentNumber, &d.ContainerID, &d.ContainerName,
// 		&d.ImageTag, &d.Logs, &d.BuildLogsPath, &d.Status, &d.Stage, &d.Progress,
// 		&d.ErrorMessage, &d.CreatedAt, &d.StartedAt, &d.FinishedAt, &d.Duration,
// 		&d.IsActive, &d.RolledBackFrom,
// 	); err != nil {
// 		return nil, err
// 	}

// 	return &d, nil
// }

// func GetDeploymentByCommitHash(commitHash string) (*Deployment, error) {
// 	query := `
// 	SELECT id, app_id, commit_hash, commit_message, commit_author, triggered_by,
// 	       deployment_number, container_id, container_name, image_tag,
// 	       logs, build_logs_path, status, stage, progress, error_message,
// 	       created_at, started_at, finished_at, duration, is_active, rolled_back_from
// 	FROM deployments
// 	WHERE commit_hash = ?
// 	LIMIT 1
// 	`
// 	var d Deployment
// 	if err := db.QueryRow(query, commitHash).Scan(
// 		&d.ID, &d.AppID, &d.CommitHash, &d.CommitMessage, &d.CommitAuthor,
// 		&d.TriggeredBy, &d.DeploymentNumber, &d.ContainerID, &d.ContainerName,
// 		&d.ImageTag, &d.Logs, &d.BuildLogsPath, &d.Status, &d.Stage, &d.Progress,
// 		&d.ErrorMessage, &d.CreatedAt, &d.StartedAt, &d.FinishedAt, &d.Duration,
// 		&d.IsActive, &d.RolledBackFrom,
// 	); err != nil {
// 		return nil, err
// 	}

// 	return &d, nil

// }

// func GetDeploymentByAppIDAndCommitHash(appID int64, commitHash string) (*Deployment, error) {
// 	query := `
// 	SELECT id, app_id, commit_hash, commit_message, commit_author, triggered_by,
// 	       deployment_number, container_id, container_name, image_tag,
// 	       logs, build_logs_path, status, stage, progress, error_message,
// 	       created_at, started_at, finished_at, duration, is_active, rolled_back_from
// 	FROM deployments
// 	WHERE app_id = ? AND commit_hash = ?
// 	LIMIT 1
// 	`
// 	var d Deployment
// 	if err := db.QueryRow(query, appID, commitHash).Scan(
// 		&d.ID, &d.AppID, &d.CommitHash, &d.CommitMessage, &d.CommitAuthor,
// 		&d.TriggeredBy, &d.DeploymentNumber, &d.ContainerID, &d.ContainerName,
// 		&d.ImageTag, &d.Logs, &d.BuildLogsPath, &d.Status, &d.Stage, &d.Progress,
// 		&d.ErrorMessage, &d.CreatedAt, &d.StartedAt, &d.FinishedAt, &d.Duration,
// 		&d.IsActive, &d.RolledBackFrom,
// 	); err != nil {
// 		return nil, err
// 	}

// 	return &d, nil

// }

// func GetActiveDeploymentByAppID(appID int64) (*Deployment, error) {
// 	query := `
// 	SELECT id, app_id, commit_hash, commit_message, commit_author, triggered_by,
// 	       deployment_number, container_id, container_name, image_tag,
// 	       logs, build_logs_path, status, stage, progress, error_message,
// 	       created_at, started_at, finished_at, duration, is_active, rolled_back_from
// 	FROM deployments
// 	WHERE app_id = ? AND is_active = 1
// 	LIMIT 1
// 	`

// 	var d Deployment
// 	if err := db.QueryRow(query, appID).Scan(
// 		&d.ID, &d.AppID, &d.CommitHash, &d.CommitMessage, &d.CommitAuthor,
// 		&d.TriggeredBy, &d.DeploymentNumber, &d.ContainerID, &d.ContainerName,
// 		&d.ImageTag, &d.Logs, &d.BuildLogsPath, &d.Status, &d.Stage, &d.Progress,
// 		&d.ErrorMessage, &d.CreatedAt, &d.StartedAt, &d.FinishedAt, &d.Duration,
// 		&d.IsActive, &d.RolledBackFrom,
// 	); err != nil {
// 		return nil, err
// 	}

// 	return &d, nil
// }

// func GetCommitHashByDeploymentID(depID int64) (string, error) {
// 	var commitHash string
// 	err := db.QueryRow(`SELECT commit_hash FROM deployments WHERE id = ?`, depID).Scan(&commitHash)
// 	if err != nil {
// 		return "", err
// 	}
// 	return commitHash, nil
// }

// func UpdateDeploymentStatus(depID int64, status, stage string, progress int, errorMsg *string) error {
// 	query := `
// 	UPDATE deployments
// 	SET status = ?, stage = ?, progress = ?, error_message = ?,
// 	    finished_at = CASE WHEN ? IN ('success', 'failed', 'stopped') THEN CURRENT_TIMESTAMP ELSE finished_at END,
// 	    duration = CASE WHEN ? IN ('success', 'failed', 'stopped') THEN
// 	                  CAST((julianday(CURRENT_TIMESTAMP) - julianday(started_at)) * 86400 AS INTEGER)
// 	               ELSE duration END
// 	WHERE id = ?
// 	`
// 	_, err := db.Exec(query, status, stage, progress, errorMsg, status, status, depID)
// 	return err
// }

// func MarkDeploymentStarted(depID int64) error {
// 	query := `UPDATE deployments SET started_at = CURRENT_TIMESTAMP WHERE id = ?`
// 	_, err := db.Exec(query, depID)
// 	return err
// }

// func MarkDeploymentActive(depID int64, appID int64) error {
// 	_, err := db.Exec(`UPDATE deployments SET is_active = 0 WHERE app_id = ?`, appID)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = db.Exec(`UPDATE deployments SET is_active = 1 WHERE id = ?`, depID)
// 	return err
// }

//	func UpdateContainerInfo(depID int64, containerID, containerName, imageTag string) error {
//		query := `
//		UPDATE deployments
//		SET container_id = ?, container_name = ?, image_tag = ?
//		WHERE id = ?
//		`
//		_, err := db.Exec(query, containerID, containerName, imageTag, depID)
//		return err
//	}
func GetIncompleteDeployments() ([]Deployment, error) {
	var deployments []Deployment

	err := db.
		Where("status IN ?", []string{"building", "deploying"}).
		Order("created_at DESC").
		Find(&deployments).Error

	if err != nil {
		return nil, err
	}

	return deployments, nil
	// query := `
	// SELECT id, app_id, commit_hash, commit_message, commit_author, triggered_by,
	//        deployment_number, container_id, container_name, image_tag,
	//        logs, build_logs_path, status, stage, progress, error_message,
	//        created_at, started_at, finished_at, duration, is_active, rolled_back_from
	// FROM deployments
	// WHERE status = 'building' OR status = 'deploying'
	// ORDER BY created_at DESC
	// `
	//
	// rows, err := db.Query(query)
	// if err != nil {
	// 	return nil, err
	// }
	// defer rows.Close()
	//
	// var deployments []Deployment
	// for rows.Next() {
	// 	var d Deployment
	// 	if err := rows.Scan(
	// 		&d.ID, &d.AppID, &d.CommitHash, &d.CommitMessage, &d.CommitAuthor,
	// 		&d.TriggeredBy, &d.DeploymentNumber, &d.ContainerID, &d.ContainerName,
	// 		&d.ImageTag, &d.Logs, &d.BuildLogsPath, &d.Status, &d.Stage, &d.Progress,
	// 		&d.ErrorMessage, &d.CreatedAt, &d.StartedAt, &d.FinishedAt, &d.Duration,
	// 		&d.IsActive, &d.RolledBackFrom,
	// 	); err != nil {
	// 		return nil, err
	// 	}
	// 	deployments = append(deployments, d)
	// }
	//
	// if err := rows.Err(); err != nil {
	// 	return nil, err
	// }
	//
	// return deployments, nil
}
