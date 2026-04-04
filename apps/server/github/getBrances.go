package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Branch struct {
	Name string `json:"name"`
}

func GetGitHubBranches(token string, repo string) ([]Branch, error) {
	if repo == "" {
		return nil, fmt.Errorf("repo name cannot be empty")
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/branches", repo)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "StakBio-App")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("github request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("github error: %s", string(body))
	}

	var branches []Branch
	if err := json.NewDecoder(resp.Body).Decode(&branches); err != nil {
		return nil, fmt.Errorf("failed to decode github response: %w", err)
	}

	return branches, nil
}
