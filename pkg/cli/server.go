package cli

import (
	"fmt"

	"github.com/ddollar/stdcli"
	"github.com/epiphytelabs/keep/pkg/docker"
	"github.com/mitchellh/go-homedir"
)

func (e *Engine) ServerCertificate(ctx *stdcli.Context) error {

	return nil
}

func (e *Engine) ServerInstall(ctx *stdcli.Context) error {
	ctx.Startf("installing server")

	if _, err := docker.ContainerInspect("keep"); err == nil {
		return fmt.Errorf("already installed")
	}

	config, err := homedir.Expand("~/.keep")
	if err != nil {
		return err
	}

	c := docker.Container{
		Name:  "keep",
		Image: fmt.Sprintf("epiphytelabs/keep:%s", e.Version),
		Ports: []string{
			"443:443/tcp",
			"53944:53/udp",
		},
		Networks: []string{"keep"},
		Volumes: map[string]string{
			config:                 "/etc/keep",
			"/var/run/docker.sock": "/var/run/docker.sock",
		},
	}

	// if err := docker.Pull(c.Image); err != nil {
	// 	return err
	// }

	if err := c.Create(); err != nil {
		return err
	}

	if err := c.Start(); err != nil {
		return err
	}

	return ctx.OK()
}

func (e *Engine) ServerUninstall(ctx *stdcli.Context) error {
	_ = e.serverUninstallStop(ctx)
	_ = e.serverUninstallRemove(ctx)

	return nil
}

func (e *Engine) serverUninstallRemove(ctx *stdcli.Context) error {
	ctx.Startf("removing server")

	if err := docker.ContainerRm("keep"); err != nil {
		return ctx.Error(err)
	}

	return ctx.OK()
}

func (e *Engine) serverUninstallStop(ctx *stdcli.Context) error {
	ctx.Startf("stopping server")

	if err := docker.ContainerStop("keep"); err != nil {
		return ctx.Error(err)
	}

	return ctx.OK()
}
