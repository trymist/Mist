package main

import (
	"mist/internal/api"
	"mist/internal/db"
)

func main() {
	err := db.Init()
	if err != nil {
		panic(err)
	}
	api.StartServer()
}
