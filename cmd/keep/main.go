package main

import (
	"os"

	"github.com/epiphytelabs/keep/pkg/cli"
)

var (
	version = "dev"
)

func main() {
	os.Exit(cli.New(version).Execute(os.Args[1:]))
}
