package main

import (
	"mist/internal/api"
	"mist/internal/db"
	"mist/internal/utils"
)

func main() {
	err := db.Init()
	utils.Init()
	if err != nil {
		panic(err)
	}
	api.StartServer()
}
