package main

import (
	"os"

	"github.com/trymist/mist/cli"
)

func main() {
	cli.Run(os.Args[1:])
}
