package github

import "fmt"

func CreateCloneUrl(accessToken string, repoURL string) string {
	if len(repoURL) > 8 && repoURL[:8] == "https://" {
		repoURL = fmt.Sprintf("https://x-access-token:%s@%s", accessToken, repoURL[8:])
	}
	return repoURL
}
