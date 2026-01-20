package lib

import (
	"time"

	"github.com/corecollectives/mist/models"
	"github.com/corecollectives/mist/queue"
	"github.com/rs/zerolog/log"
)

func CleanupOnStartup() error {
	err := cleanupUpdates()
	if err != nil {
		return err
	}
	err = cleanupDeployments()
	if err != nil {
		return err
	}
	return nil
}

func cleanupDeployments() error {
	deployments, err := models.GetIncompleteDeployments()
	if err != nil {
		return err
	}
	if len(deployments) == 0 {
		return nil
	}

	log.Info().
		Int("count", len(deployments)).
		Msg("Found incomplete deployments on startup, cleaning up")

	errorMsg := "system died before deployment could complete"
	for _, dep := range deployments {
		if dep.Status == "deploying" || dep.Status == "building" {
			log.Warn().
				Int64("deployment_id", dep.ID).
				Str("status", string(dep.Status)).
				Msg("Marking interrupted deployment as failed")

			err = models.UpdateDeploymentStatus(dep.ID, "failed", "failed", 0, &errorMsg)
			if err != nil {
				log.Error().Err(err).Int64("deployment_id", dep.ID).Msg("Failed to mark deployment as failed")
				return err
			}
		} else if dep.Status == "pending" {
			log.Info().
				Int64("deployment_id", dep.ID).
				Msg("Re-queuing pending deployment")

			err = queue.GetQueue().AddJob(dep.ID)
			if err != nil {
				log.Error().Err(err).Int64("deployment_id", dep.ID).Msg("Failed to re-queue pending deployment")
				return err
			}
		} else {
			continue
		}
	}

	log.Info().Msg("Deployment cleanup completed successfully")
	return nil
}

func cleanupUpdates() error {
	logs, err := models.GetUpdateLogs(1)
	if err != nil {
		return err
	}

	if len(logs) == 0 {
		return nil
	}

	latestLog := logs[0]

	if latestLog.Status != "in_progress" {
		return nil
	}

	log.Info().
		Int64("update_log_id", latestLog.ID).
		Str("from_version", latestLog.VersionFrom).
		Str("to_version", latestLog.VersionTo).
		Str("age", time.Since(latestLog.StartedAt).String()).
		Msg("Found in-progress update on startup, checking status")

	currentVersion, err := models.GetSystemSetting("version")
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

		completionLog := *latestLog.Logs + "\n✅ Update completed successfully (verified on restart)\n"
		err = models.UpdateUpdateLogStatus(latestLog.ID, "success", completionLog, nil)
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
	failureLog := *latestLog.Logs + "\n❌ " + errMsg + "\n"
	err = models.UpdateUpdateLogStatus(latestLog.ID, "failed", failureLog, &errMsg)
	if err != nil {
		log.Error().Err(err).Int64("update_log_id", latestLog.ID).Msg("Failed to mark failed update")
		return err
	}

	log.Info().
		Int64("update_log_id", latestLog.ID).
		Msg("Marked failed update as failed")

	return nil
}
