package docker

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/moby/moby/client"
	"github.com/rs/zerolog/log"
)

// remove more that keepCount images to save storage
func CleanupOldImages(appID int64, keepCount int) error {
	if keepCount < 1 {
		keepCount = 5
	}

	// TODO: this part is wrong and isn't working, fix it and include it
	imagePattern := fmt.Sprintf("mist-app-%d-", appID)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	cli, err := client.New(client.FromEnv)
	if err != nil {
		return fmt.Errorf("error creating moby client: %s", err.Error())
	}

	filterArgs := make(client.Filters)
	filterArgs.Add("reference", fmt.Sprintf("%s*", imagePattern))

	imageListResult, err := cli.ImageList(ctx, client.ImageListOptions{
		Filters: filterArgs,
	})
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("listing images timed out")
		}
		return fmt.Errorf("failed to list images: %w", err)
	}

	imageList := imageListResult.Items

	if len(imageList) == 0 {
		return nil
	}

	if len(imageList) <= keepCount {
		return nil
	}

	type imageInfo struct {
		id      string
		created int64
	}

	var images []imageInfo
	for _, img := range imageList {
		images = append(images, imageInfo{
			id:      img.ID,
			created: img.Created,
		})
	}

	sort.Slice(images, func(i, j int) bool {
		return images[i].created > images[j].created
	})

	if len(images) > keepCount {
		imagesToRemove := images[keepCount:]

		for _, img := range imagesToRemove {
			rmiCtx, rmiCancel := context.WithTimeout(context.Background(), 1*time.Minute)
			_, err := cli.ImageRemove(rmiCtx, img.id, client.ImageRemoveOptions{
				Force: true,
			})
			if err != nil {
				log.Warn().Err(err).Str("image_id", img.id).Msg("Failed to remove old image")
			}
			rmiCancel()
		}
	}

	return nil

	// legacy exec method
	//
	//
	// imagePattern := fmt.Sprintf("mist-app-%d-", appID)
	//
	// ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	// defer cancel()
	//
	// listCmd := exec.CommandContext(ctx, "docker", "images",
	// 	"--filter", fmt.Sprintf("reference=%s*", imagePattern),
	// 	"--format", "{{.Repository}}:{{.Tag}} {{.CreatedAt}}",
	// )
	//
	// output, err := listCmd.Output()
	// if err != nil {
	// 	if ctx.Err() == context.DeadlineExceeded {
	// 		return fmt.Errorf("listing images timed out")
	// 	}
	// 	return fmt.Errorf("failed to list images: %w", err)
	// }
	//
	// if len(output) == 0 {
	// 	return nil
	// }
	//
	// lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	// if len(lines) <= keepCount {
	// 	return nil
	// }
	//
	// type imageInfo struct {
	// 	name      string
	// 	timestamp string
	// }
	//
	// var images []imageInfo
	// for _, line := range lines {
	// 	if line == "" {
	// 		continue
	// 	}
	// 	parts := strings.SplitN(line, " ", 2)
	// 	if len(parts) == 2 {
	// 		images = append(images, imageInfo{
	// 			name:      parts[0],
	// 			timestamp: parts[1],
	// 		})
	// 	}
	// }
	//
	// sort.Slice(images, func(i, j int) bool {
	// 	return images[i].timestamp > images[j].timestamp
	// })
	//
	// if len(images) > keepCount {
	// 	imagesToRemove := images[keepCount:]
	//
	// 	for _, img := range imagesToRemove {
	// 		rmiCtx, rmiCancel := context.WithTimeout(context.Background(), 1*time.Minute)
	// 		rmiCmd := exec.CommandContext(rmiCtx, "docker", "rmi", "-f", img.name)
	// 		if err := rmiCmd.Run(); err != nil {
	// 			log.Warn().Err(err).Str("image", img.name).Msg("Failed to remove old image")
	// 		}
	// 		rmiCancel()
	// 	}
	// }
	//
	// return nil
}

// cleanup dangling images, (triggered from the dashboard)
func CleanupDanglingImages() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	cli, err := client.New(client.FromEnv)
	if err != nil {
		return fmt.Errorf("error creating moby client: %s", err.Error())
	}

	filterArgs := make(client.Filters)
	filterArgs.Add("dangling", "true")

	_, err = cli.ImagePrune(ctx, client.ImagePruneOptions{
		Filters: filterArgs,
	})
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("pruning images timed out")
		}
		return fmt.Errorf("failed to prune dangling images: %w", err)
	}

	return nil

	// legacy exec method
	//
	//
	// ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	// defer cancel()
	//
	// pruneCmd := exec.CommandContext(ctx, "docker", "image", "prune", "-f")
	// if err := pruneCmd.Run(); err != nil {
	// 	if ctx.Err() == context.DeadlineExceeded {
	// 		return fmt.Errorf("pruning images timed out")
	// 	}
	// 	return fmt.Errorf("failed to prune dangling images: %w", err)
	// }
	//
	// return nil
}

// cleanup stopped containers, triggered from dashboard
func CleanupStoppedContainers() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	cli, err := client.New(client.FromEnv)
	if err != nil {
		return fmt.Errorf("error creating moby client: %s", err.Error())
	}

	_, err = cli.ContainerPrune(ctx, client.ContainerPruneOptions{})
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("pruning containers timed out")
		}
		return fmt.Errorf("failed to prune stopped containers: %w", err)
	}

	return nil

	// legacy exec method
	//
	//
	// ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	// defer cancel()
	//
	// pruneCmd := exec.CommandContext(ctx, "docker", "container", "prune", "-f")
	// if err := pruneCmd.Run(); err != nil {
	// 	if ctx.Err() == context.DeadlineExceeded {
	// 		return fmt.Errorf("pruning containers timed out")
	// 	}
	// 	return fmt.Errorf("failed to prune stopped containers: %w", err)
	// }
	//
	// return nil
}

// system prune, triggered from dashboard
func SystemPrune() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	cli, err := client.New(client.FromEnv)
	if err != nil {
		return "", fmt.Errorf("error creating moby client: %s", err.Error())
	}

	var output string

	containerReport, err := cli.ContainerPrune(ctx, client.ContainerPruneOptions{})
	if err != nil {
		return output, fmt.Errorf("failed to prune containers: %w", err)
	}
	output += fmt.Sprintf("Deleted Containers: %v\n", containerReport.Report.ContainersDeleted)
	output += fmt.Sprintf("Space Reclaimed: %d bytes\n", containerReport.Report.SpaceReclaimed)

	filterArgs := make(client.Filters)
	filterArgs.Add("dangling", "true")
	imageReport, err := cli.ImagePrune(ctx, client.ImagePruneOptions{
		Filters: filterArgs,
	})
	if err != nil {
		return output, fmt.Errorf("failed to prune images: %w", err)
	}
	output += fmt.Sprintf("Deleted Images: %v\n", imageReport.Report.ImagesDeleted)
	output += fmt.Sprintf("Space Reclaimed: %d bytes\n", imageReport.Report.SpaceReclaimed)

	networkReport, err := cli.NetworkPrune(ctx, client.NetworkPruneOptions{})
	if err != nil {
		return output, fmt.Errorf("failed to prune networks: %w", err)
	}
	output += fmt.Sprintf("Deleted Networks: %v\n", networkReport.Report.NetworksDeleted)

	buildCacheReport, err := cli.BuildCachePrune(ctx, client.BuildCachePruneOptions{})
	if err != nil {
		return output, fmt.Errorf("failed to prune build cache: %w", err)
	}
	output += fmt.Sprintf("Space Reclaimed from build cache: %d bytes\n", buildCacheReport.Report.SpaceReclaimed)

	if ctx.Err() == context.DeadlineExceeded {
		return output, fmt.Errorf("system prune timed out")
	}

	return output, nil

	// legacy exec method
	//
	//
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	// defer cancel()
	//
	// pruneCmd := exec.CommandContext(ctx, "docker", "system", "prune", "-f")
	// output, err := pruneCmd.CombinedOutput()
	// if err != nil {
	// 	if ctx.Err() == context.DeadlineExceeded {
	// 		return "", fmt.Errorf("system prune timed out")
	// 	}
	// 	return string(output), fmt.Errorf("failed to run system prune: %w", err)
	// }
	//
	// return string(output), nil
}

// system prune all, triggered from dashboard
func SystemPruneAll() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	cli, err := client.New(client.FromEnv)
	if err != nil {
		return "", fmt.Errorf("error creating moby client: %s", err.Error())
	}

	var output string

	containerReport, err := cli.ContainerPrune(ctx, client.ContainerPruneOptions{})
	if err != nil {
		return output, fmt.Errorf("failed to prune containers: %w", err)
	}
	output += fmt.Sprintf("Deleted Containers: %v\n", containerReport.Report.ContainersDeleted)
	output += fmt.Sprintf("Space Reclaimed: %d bytes\n", containerReport.Report.SpaceReclaimed)

	imageReport, err := cli.ImagePrune(ctx, client.ImagePruneOptions{})
	if err != nil {
		return output, fmt.Errorf("failed to prune images: %w", err)
	}
	output += fmt.Sprintf("Deleted Images: %v\n", imageReport.Report.ImagesDeleted)
	output += fmt.Sprintf("Space Reclaimed: %d bytes\n", imageReport.Report.SpaceReclaimed)

	networkReport, err := cli.NetworkPrune(ctx, client.NetworkPruneOptions{})
	if err != nil {
		return output, fmt.Errorf("failed to prune networks: %w", err)
	}
	output += fmt.Sprintf("Deleted Networks: %v\n", networkReport.Report.NetworksDeleted)

	buildCacheReport, err := cli.BuildCachePrune(ctx, client.BuildCachePruneOptions{
		All: true,
	})
	if err != nil {
		return output, fmt.Errorf("failed to prune build cache: %w", err)
	}
	output += fmt.Sprintf("Space Reclaimed from build cache: %d bytes\n", buildCacheReport.Report.SpaceReclaimed)

	if ctx.Err() == context.DeadlineExceeded {
		return output, fmt.Errorf("aggressive system prune timed out")
	}

	return output, nil

	// legacy exec method
	//
	//
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	// defer cancel()
	//
	// pruneCmd := exec.CommandContext(ctx, "docker", "system", "prune", "-a", "-f")
	// output, err := pruneCmd.CombinedOutput()
	// if err != nil {
	// 	if ctx.Err() == context.DeadlineExceeded {
	// 		return "", fmt.Errorf("aggressive system prune timed out")
	// 	}
	// 	return string(output), fmt.Errorf("failed to run aggressive system prune: %w", err)
	// }
	//
	// return string(output), nil
}
