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

func BuildImage(ctx context.Context, imageTag, contextPath string, envVars map[string]string, logfile *os.File) error {
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
	env := make(map[string]*string)
	for k, v := range envVars {
		val := v
		env[k] = &val
	}
	buildOptions := client.ImageBuildOptions{
		Tags:      tags,
		Remove:    true,
		BuildArgs: env,
	}

	log.Info().Str("image_tag", imageTag).Msg("Building Docker image")

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

	// legacy exec method
	//
	//
	// buildArgs := []string{"build", "-t", imageTag}
	//
	// for key, value := range envVars {
	// 	buildArgs = append(buildArgs, "--build-arg", fmt.Sprintf("%s=%s", key, value))
	// }
	//
	// buildArgs = append(buildArgs, contextPath)
	//
	// log.Debug().Strs("build_args", buildArgs).Str("image_tag", imageTag).Msg("Building Docker image")
	//
	// ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	// defer cancel()
	//
	// cmd := exec.CommandContext(ctx, "docker", buildArgs...)
	// cmd.Stdout = logfile
	// cmd.Stderr = logfile
	//
	// if err := cmd.Run(); err != nil {
	// 	if ctx.Err() == context.DeadlineExceeded {
	// 		return fmt.Errorf("docker build timed out after 15 minutes")
	// 	}
	// 	exitCode := -1
	// 	if exitErr, ok := err.(*exec.ExitError); ok {
	// 		exitCode = exitErr.ExitCode()
	// 	}
	// 	return fmt.Errorf("docker build failed with exit code %d: %w", exitCode, err)
	// }
	// return nil
}

func PullDockerImage(ctx context.Context, imageName string, logfile *os.File) error {
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

	// legacy exec method
	//
	//
	// pullCmd := exec.CommandContext(ctx, "docker", "pull", imageName)
	// pullCmd.Stdout = logfile
	// pullCmd.Stderr = logfile
	//
	// if err := pullCmd.Run(); err != nil {
	// 	if ctx.Err() == context.DeadlineExceeded {
	// 		return fmt.Errorf("docker pull timed out after 15 minutes for image %s", imageName)
	// 	}
	// 	exitCode := -1
	// 	if exitErr, ok := err.(*exec.ExitError); ok {
	// 		exitCode = exitErr.ExitCode()
	// 	}
	// 	return fmt.Errorf("docker pull failed with exit code %d: %w", exitCode, err)
	// }
	// return nil
}
