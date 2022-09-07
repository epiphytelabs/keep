package cli

import (
	"fmt"

	"github.com/ddollar/stdcli"
	"github.com/epiphytelabs/keep/pkg/docker"
	"github.com/epiphytelabs/keep/pkg/repository"
)

func (e *Engine) Uninstall(ctx *stdcli.Context) error {
	name := ctx.Arg(0)

	a, err := repository.Get(name)
	if err != nil {
		return err
	}

	if !a.Installed {
		return fmt.Errorf("not installed")
	}

	if err := e.uninstallRemoveApp(ctx, a); err != nil {
		_ = ctx.Error(err)
	}

	for _, r := range a.Resources.List() {
		if err := e.uninstallRemoveResource(ctx, r); err != nil {
			_ = ctx.Error(err)
		}
	}

	if err := e.uninstallRemoveNetwork(ctx, name); err != nil {
		_ = ctx.Error(err)
	}

	return nil
}

func (e *Engine) uninstallRemoveApp(ctx *stdcli.Context, a *repository.App) error {
	ctx.Startf("removing app <id>%s</id>", a.Name)

	c, err := a.Container()
	if err != nil {
		return err
	}

	if err := c.Stop(); err != nil {
		return err
	}

	if err := c.Rm(); err != nil {
		return err
	}

	return ctx.OK()
}

func (e *Engine) uninstallRemoveNetwork(ctx *stdcli.Context, name string) error {
	ctx.Startf("removing network <id>%s</id>", name)

	if err := docker.NetworkRemove(fmt.Sprintf("keep-%s", name)); err != nil {
		return err
	}

	return ctx.OK()
}

func (e *Engine) uninstallRemoveResource(ctx *stdcli.Context, r repository.Resource) error {
	ctx.Startf("removing resource <id>%s</id>", r.Name())

	c, err := r.Container()
	if err != nil {
		return err
	}

	if err := c.Stop(); err != nil {
		return err
	}

	if err := c.Rm(); err != nil {
		return err
	}

	return ctx.OK()
}
