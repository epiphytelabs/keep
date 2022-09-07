package cli

import (
	"fmt"

	"github.com/ddollar/stdcli"
	"github.com/epiphytelabs/keep/pkg/docker"
	"github.com/epiphytelabs/keep/pkg/repository"
)

func (e *Engine) Install(ctx *stdcli.Context) error {
	name := ctx.Arg(0)

	a, err := repository.Get(name)
	if err != nil {
		return err
	}

	if a.Installed {
		return fmt.Errorf("already installed")
	}

	net, err := e.installCreateNetwork(ctx, name)
	if err != nil {
		return err
	}

	for _, r := range a.Resources.List() {
		if err := e.installCreateResource(ctx, net, r); err != nil {
			return err
		}
	}

	if err := e.installCreateApp(ctx, net, a); err != nil {
		return err
	}

	ctx.Writef("url: <url>https://%s.app.keep</url>", name)

	return nil
}

func (e *Engine) installCreateApp(ctx *stdcli.Context, net *docker.Network, a *repository.App) error {
	ctx.Startf("creating app <id>%s</id>", a.Name)

	c, err := a.Container()
	if err != nil {
		return err
	}

	if err := docker.Pull(c.Image); err != nil {
		return err
	}

	if err := c.Create(); err != nil {
		return err
	}

	if err := net.Connect(*c); err != nil {
		return err
	}

	if err := c.Start(); err != nil {
		return err
	}

	return ctx.OK()
}

func (e *Engine) installCreateNetwork(ctx *stdcli.Context, name string) (*docker.Network, error) {
	ctx.Startf("creating network <id>%s</id>", name)

	net, err := docker.NetworkCreate(fmt.Sprintf("keep-%s", name))
	if err != nil {
		return nil, err
	}

	return net, ctx.OK()
}

func (e *Engine) installCreateResource(ctx *stdcli.Context, net *docker.Network, r repository.Resource) error {
	ctx.Startf("creating resource <id>%s</id>", r.Name())

	c, err := r.Container()
	if err != nil {
		return err
	}

	if err := docker.Pull(c.Image); err != nil {
		return err
	}

	if err := c.Create(); err != nil {
		return err
	}

	if err := net.Connect(*c); err != nil {
		return err
	}

	if err := c.Start(); err != nil {
		return err
	}

	return ctx.OK()
}
