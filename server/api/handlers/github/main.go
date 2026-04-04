package github

import (
	"encoding/base64"
	"encoding/json"
	"os"
)

type StateData struct {
	AppId  int `json:"appId"`
	UserId int `json:"userId"`
}

func GetFrontendBaseUrl() string {
	if os.Getenv("ENV") == "dev" {
		return "http://localhost:5173"
	}
	return ""
}

func GenerateState(appId int, userId int) string {
	payload := map[string]int{
		"appId":  appId,
		"userId": userId,
	}
	jsonBytes, _ := json.Marshal(payload)
	encoded := base64.StdEncoding.EncodeToString(jsonBytes)
	return encoded
}
