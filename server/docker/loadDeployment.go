package docker

import (
	"github.com/corecollectives/mist/models"
	"gorm.io/gorm"
)

func LoadDeployment(depId int64, db *gorm.DB) (*models.Deployment, error) {
	// row := db.QueryRow("SELECT id, app_id, commit_hash, commit_message, triggered_by, logs, status, created_at, finished_at FROM deployments WHERE id = ?", depId)
	dep := &models.Deployment{}
	err := db.Table("deployments").Where("id=?", depId).First(&dep).Error

	if err != nil {
		return nil, err
	}
	return dep, nil

}
