package utils

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type DeploymentLogger struct {
	logger     zerolog.Logger
	deployID   int64
	appID      int64
	commitHash string
}

func NewDeploymentLogger(depID, appID int64, commit string) *DeploymentLogger {
	return &DeploymentLogger{
		logger: log.With().
			Int64("deployment_id", depID).
			Int64("app_id", appID).
			Str("commit_hash", commit).
			Logger(),
		deployID:   depID,
		appID:      appID,
		commitHash: commit,
	}
}

func (dl *DeploymentLogger) Info(message string) {
	dl.logger.Info().Msg(message)
}

func (dl *DeploymentLogger) InfoWithFields(message string, fields map[string]interface{}) {
	event := dl.logger.Info()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(message)
}

func (dl *DeploymentLogger) Error(err error, message string) {
	dl.logger.Error().Err(err).Msg(message)
}

func (dl *DeploymentLogger) ErrorWithFields(err error, message string, fields map[string]interface{}) {
	event := dl.logger.Error().Err(err)
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(message)
}

func (dl *DeploymentLogger) Warn(message string) {
	dl.logger.Warn().Msg(message)
}

func (dl *DeploymentLogger) WarnWithFields(message string, fields map[string]interface{}) {
	event := dl.logger.Warn()
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(message)
}

func (dl *DeploymentLogger) Debug(message string) {
	dl.logger.Debug().Msg(message)
}

func InitLogger() {
	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
}
