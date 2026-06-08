package main

import (
	"mist/internal/db"
)

func main() {
	err := db.Init()
	if err != nil {
		panic(err)
	}
}
