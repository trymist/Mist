package git

import (
	"context"
	"fmt"
	"os"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/rs/zerolog/log"
)

func CloneRepo(ctx context.Context, url string, branch string, logFile *os.File, path string) error {
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
		return err
	}

	return nil
}
