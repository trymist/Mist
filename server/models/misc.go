package models

type LatestCommit struct {
	SHA     string `json:"sha"`
	Message string `json:"message"`
	URL     string `json:"html_url"`
	Author  string `json:"author"`
}
