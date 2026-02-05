package fs

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/corecollectives/mist/constants"
	"github.com/corecollectives/mist/models"
)

var logPath string = constants.Constants["LogPath"].(string)

// helper function to create build log file
// build logs files are located in `/var/lib/mist/logs/{dep-commit-hash + appID + _build_logs}`
func CreateDockerBuildLogFile(depID int64) (*os.File, string, error) {
	commitHash, err := models.GetCommitHashByDeploymentID(depID)
	if err != nil {
		return nil, "", err
	}
	err = CreateDirIfNotExists(logPath, os.ModePerm)
	if err != nil {
		return nil, "", err
	}

	logFileName := commitHash + strconv.FormatInt(depID, 10) + "_build_logs"
	logPath := filepath.Join(logPath, logFileName)

	logfile, err := os.Create(logPath)
	if err != nil {
		return nil, "", err
	}
	return logfile, logPath, nil

}
