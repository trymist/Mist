package docker

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/moby/go-archive"
	"github.com/moby/moby/client"
	"github.com/rs/zerolog/log"
)

func BuildDockerImageWithBuildArgs(ctx context.Context, imageTag, contextPath string, buildArgs map[string]string, logfile *os.File) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, 15*time.Minute)
	defer cancel()
	cli, err := client.New(client.FromEnv)
	if err != nil {
		return fmt.Errorf("error opening moby client: %s", err.Error())
	}
	buildCtx, err := archive.TarWithOptions(contextPath, &archive.TarOptions{
		ExcludePatterns: []string{},
	})

	if err != nil {
		return fmt.Errorf("error building build Context")
	}
	var tags []string
	tags = append(tags, imageTag)
	buildArgsMap := make(map[string]*string)
	for k, v := range buildArgs {
		val := v
		buildArgsMap[k] = &val
	}
	buildOptions := client.ImageBuildOptions{
		Tags:      tags,
		Remove:    true,
		BuildArgs: buildArgsMap,
	}

	log.Info().Str("image_tag", imageTag).Int("build_args_count", len(buildArgs)).Msg("Building Docker image with build-time arguments")

	resp, err := cli.ImageBuild(timeoutCtx, buildCtx, buildOptions)
	if err != nil {
		if timeoutCtx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("image build timed out after 15 minutes")
		}
		if timeoutCtx.Err() == context.Canceled {
			return context.Canceled
		}
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(logfile, resp.Body)
	if err != nil {
		return err
	}
	return nil

}

func PullPrebuiltDockerImage(ctx context.Context, imageName string, logfile *os.File) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, 15*time.Minute)
	defer cancel()

	cli, err := client.New(client.FromEnv)
	if err != nil {
		return fmt.Errorf("error opening moby client: %s", err.Error())
	}

	log.Debug().Str("image_name", imageName).Msg("pulling image")
	resp, err := cli.ImagePull(timeoutCtx, imageName, client.ImagePullOptions{})
	if err != nil {
		if timeoutCtx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("image pull timed out after 15 minutes")
		}
		if timeoutCtx.Err() == context.Canceled {
			return context.Canceled
		}
		return err
	}
	defer resp.Close()
	_, err = io.Copy(logfile, resp)
	if err != nil {
		return err
	}
	return nil

}
