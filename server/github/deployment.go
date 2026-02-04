package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func CreateDeployment(repo string, branch string, userId int) (int64, error) {

	token, _, err := GetGitHubAccessToken(userId)
	if err != nil {
		return 0, fmt.Errorf("error getting GH token %w", err)
	}
	url := fmt.Sprintf("https://api.github.com/repos/%s/deployments", repo)
	body := fmt.Sprintf(`{
		"ref":"%s",
		"environment":"production",
		"required_contexts":[]
	}`, branch)

	req, _ := http.NewRequest("POST", url, strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{Timeout: 30 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		b, _ := io.ReadAll(res.Body)
		return 0, fmt.Errorf("github api error: %s", b)
	}

	var out struct {
		ID int64 `json:"id"`
	}

	json.NewDecoder(res.Body).Decode(&out)

	return out.ID, nil
}

func UpdateDeployment(repo string, depID int64, state string, message string, userID int) error {
	token, _, err := GetGitHubAccessToken(userID)
	if err != nil {
		return fmt.Errorf("error getting GH token %w", err)
	}
	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/deployments/%d/statuses",
		repo,
		depID,
	)
	body := fmt.Sprintf(`{
		"state":"%s",
		"description":"%s"
	}`, state, message)

	req, _ := http.NewRequest("POST", url, strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 300 {
		b, _ := io.ReadAll(res.Body)
		return fmt.Errorf("github api error: %s", b)
	}
	return nil

}
