package main

import (
	"os"

	"github.com/epiphytelabs/keep/pkg/cli"
)

func main() {
	os.Exit(cli.New().Execute(os.Args[1:]))
}
