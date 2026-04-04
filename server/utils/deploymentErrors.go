package utils

import (
	"fmt"
	"time"
)

type DeploymentStage string

const (
	StagePending     DeploymentStage = "pending"
	StageCloning     DeploymentStage = "cloning"
	StageBuilding    DeploymentStage = "building"
	StageDeploying   DeploymentStage = "deploying"
	StageSuccess     DeploymentStage = "success"
	StageFailed      DeploymentStage = "failed"
	StageRollingBack DeploymentStage = "rolling_back"
)

type DeploymentError struct {
	Stage     DeploymentStage
	DeployID  int64
	AppID     int64
	Err       error
	Message   string
	Retryable bool
	Timestamp time.Time
}

func (de *DeploymentError) Error() string {
	return fmt.Sprintf("[%s] deployment %d failed at %s: %s (error: %v)",
		de.Timestamp.Format(time.RFC3339), de.DeployID, de.Stage, de.Message, de.Err)
}

func NewDeploymentError(stage DeploymentStage, depID, appID int64, err error, message string, retryable bool) *DeploymentError {
	return &DeploymentError{
		Stage:     stage,
		DeployID:  depID,
		AppID:     appID,
		Err:       err,
		Message:   message,
		Retryable: retryable,
		Timestamp: time.Now(),
	}
}

func GetProgressFromStage(stage string) int {
	switch DeploymentStage(stage) {
	case StagePending:
		return 0
	case StageCloning:
		return 20
	case StageBuilding:
		return 50
	case StageDeploying:
		return 80
	case StageSuccess:
		return 100
	case StageFailed:
		return 0
	default:
		return 0
	}
}

func GetStageMessage(stage string) string {
	switch DeploymentStage(stage) {
	case StagePending:
		return "Deployment queued and waiting to start"
	case StageCloning:
		return "Cloning repository from Git"
	case StageBuilding:
		return "Building Docker image"
	case StageDeploying:
		return "Deploying container"
	case StageSuccess:
		return "Deployment completed successfully"
	case StageFailed:
		return "Deployment failed"
	case StageRollingBack:
		return "Rolling back to previous version"
	default:
		return "Unknown stage"
	}
}
