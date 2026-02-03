package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/corecollectives/mist/constants"
	"github.com/corecollectives/mist/github"
	"github.com/corecollectives/mist/models"
)

func GetLatestCommit(appID int64, userID int64) (*models.LatestCommit, error) {
	gitProvider, err := models.GetGitProviderNameByAppID(appID)
	if err != nil {
		fmt.Printf("error getting provider name by app id", err.Error())
		return nil, err
	}

	gitCloneUrl, err := models.GetCloneUrlfromAppID(appID)
	if err != nil {
		return nil, err
	}
	if gitProvider == nil && gitCloneUrl != nil {
		fmt.Printf("no provider found")
		return latestRemoteCommit(*gitCloneUrl)
	} else if gitProvider == nil && gitCloneUrl == nil {
		return nil, fmt.Errorf("git url or provider not given")
	}

	if *gitProvider == models.GitProviderGitHub {
		fmt.Printf("github provider found")
		return github.GetLatestCommit(appID, userID)
	}
	return nil, fmt.Errorf("failed to get latest commit")

}

// for repositories which are linked via git clone url and not any github provider
// we temporarily clone it (not fully) to get the latest commit information
// the path for this temp repo is `/var/lib/mist/git-meta/repo-`
func latestRemoteCommit(repoURL string) (*models.LatestCommit, error) {

	tmpPath := filepath.Join(constants.Constants["RootPath"].(string), "git-meta")
	if err := os.MkdirAll(tmpPath, 0o755); err != nil {
		return nil, err
	}
	tmpDir, err := os.MkdirTemp(tmpPath, "repo-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("git init failed: %w: %s", err, out)
	}

	cmd = exec.Command("git", "remote", "add", "origin", repoURL)
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("git remote add failed: %w: %s", err, out)
	}

	// use `--filter=blob:none` to not to fetch any actual file at this point, we only fetching the
	// git metadata not the actual repo
	cmd = exec.Command(
		"git",
		"fetch",
		"--depth=1",
		"--filter=blob:none",
		"origin",
		"HEAD",
	)
	cmd.Dir = tmpDir
	if out, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("git fetch failed: %w: %s", err, out)
	}

	cmd = exec.Command(
		"git",
		"show",
		"-s",
		"--format=%H%n%an%n%s",
		"FETCH_HEAD",
	)
	cmd.Dir = tmpDir

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git show failed: %w", err)
	}

	parts := strings.SplitN(strings.TrimSpace(string(out)), "\n", 3)
	if len(parts) < 3 {
		return nil, fmt.Errorf("unexpected git show output")
	}

	return &models.LatestCommit{
		SHA:     parts[0],
		Author:  parts[1],
		Message: parts[2],
		URL:     "",
	}, nil
}
