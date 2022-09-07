package cli

import (
	"fmt"

	"github.com/ddollar/stdcli"
	"github.com/epiphytelabs/keep/pkg/docker"
	"github.com/epiphytelabs/keep/pkg/repository"
)

func (e *Engine) List(ctx *stdcli.Context) error {
	if err := docker.Info(); err != nil {
		return fmt.Errorf("could not connect to docker")
	}

	is, err := repository.Installed()
	if err != nil {
		return err
	}

	for _, i := range is {
		fmt.Println(i)
	}

	return nil
}
