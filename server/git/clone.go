package git

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/corecollectives/mist/github"
	"github.com/corecollectives/mist/models"
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/rs/zerolog/log"
)

func CloneGitRepo(ctx context.Context, url string, branch string, logFile *os.File, path string) error {
	_, err := fmt.Fprintf(logFile, "[GIT]: Cloning into %s\n", path)
	if err != nil {
		log.Warn().Msg("error logging into log file")
	}
	_, err = git.PlainCloneContext(ctx, path, &git.CloneOptions{
		URL: url,
		// Progress:      logFile,
		ReferenceName: plumbing.NewBranchReferenceName(branch),
		SingleBranch:  true,
	})
	if err != nil {
		if ctx.Err() == context.Canceled {
			return fmt.Errorf("deployment stopped by user")
		}
		return err
	}

	return nil
}

func CloneRepo(ctx context.Context, appId int64, logFile *os.File) error {
	log.Info().Int64("app_id", appId).Msg("Starting repository clone")

	userId, err := models.GetUserIDByAppID(appId)
	if err != nil {
		return fmt.Errorf("failed to get user id by app id: %w", err)
	}

	cloneURL, accessToken, shouldMigrate, err := models.GetAppCloneURL(appId, *userId)
	if err != nil {
		return fmt.Errorf("failed to get clone URL: %w", err)
	}

	_, _, branch, _, projectId, name, err := models.GetAppGitInfo(appId)
	if err != nil {
		return fmt.Errorf("failed to fetch app: %w", err)
	}

	if shouldMigrate {
		log.Info().Int64("app_id", appId).Msg("Migrating legacy app to new git format")
		// for legacy GitHub apps, we don't have a git_provider_id
		// we just update the git_clone_url
		err = models.UpdateAppGitCloneURL(appId, cloneURL, nil)
		if err != nil {
			log.Warn().Err(err).Int64("app_id", appId).Msg("Failed to migrate app git info, continuing anyway")
		}
	}

	// construct authenticated clone URL if we have an access token
	repoURL := cloneURL
	if accessToken != "" {
		// insert token into the URL
		// for GitHub: https://x-access-token:TOKEN@github.com/user/repo.git
		// for GitLab: https://oauth2:TOKEN@gitlab.com/user/repo.git
		// for Bitbucket: https://x-token-auth:TOKEN@bitbucket.org/user/repo.git
		// for Gitea: https://TOKEN@gitea.com/user/repo.git

		// simple approach: insert token after https://
		// if len(cloneURL) > 8 && cloneURL[:8] == "https://" {
		// 	repoURL = fmt.Sprintf("https://x-access-token:%s@%s", accessToken, cloneURL[8:])
		// }
		repoURL = github.CreateCloneUrl(accessToken, repoURL)
	}

	path := fmt.Sprintf("/var/lib/mist/projects/%d/apps/%s", projectId, name)

	if _, err := os.Stat(path + "/.git"); err == nil {
		log.Info().Str("path", path).Msg("Repository already exists, removing directory")

		if err := os.RemoveAll(path); err != nil {
			return fmt.Errorf("failed to remove existing repository: %w", err)

		}
	}

	if err := os.MkdirAll(path, 0o755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	log.Info().Str("clone_url", cloneURL).Str("branch", branch).Str("path", path).Msg("Cloning repository")

	ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	// old command implementation
	//
	//
	// cmd := exec.CommandContext(ctx, "git", "clone", "--branch", branch, repoURL, path)
	// output, err := cmd.CombinedOutput()
	// lines := strings.Split(string(output), "\n")
	// for _, line := range lines {
	// 	if len(line) > 0 {
	// 		fmt.Fprintf(logFile, "[GIT] %s\n", line)
	// 	}
	// }

	// new git sdk implementation
	err = CloneGitRepo(ctx, repoURL, branch, logFile, path)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("git clone timed out after 10 minutes")
		}
		return fmt.Errorf("error cloning repository: %v\n", err)
	}

	log.Info().Int64("app_id", appId).Str("path", path).Msg("Repository cloned successfully")
	return nil
}
