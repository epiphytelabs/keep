package cli

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/ddollar/stdcli"
	"github.com/epiphytelabs/keep/pkg/docker"
	"github.com/epiphytelabs/keep/pkg/repository"
	"github.com/pkg/browser"
)

func (e *Engine) Install(ctx *stdcli.Context) error {
	if err := docker.Info(); err != nil {
		return fmt.Errorf("could not connect to docker")
	}

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

	url := fmt.Sprintf("https://%s.app.keep", name)

	if err := e.installWait(ctx, url); err != nil {
		return err
	}

	_ = ctx.Writef("%s\n", url)

	_ = browser.OpenURL(url)

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

func (e *Engine) installWait(ctx *stdcli.Context, url string) error {
	ctx.Startf("waiting for app to start")

	tick := time.NewTicker(5 * time.Second)
	timeout := time.After(5 * time.Minute)

	c := &http.Client{
		Timeout: 2 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	for {
		select {
		case <-tick.C:
			res, err := c.Get(url)
			if err == nil && res.StatusCode == 200 {
				return ctx.OK()
			}
		case <-timeout:
			return fmt.Errorf("timeout")
		}
	}
}
