package cli

import (
	"github.com/ddollar/stdcli"
)

var (
	Version = "dev"
)

type Engine struct {
	*stdcli.Engine
}

func New() *Engine {
	e := &Engine{Engine: stdcli.New("keep", Version)}

	e.Register()

	return e
}

func (e *Engine) Register() {
	e.Command("install", "install an application", e.Install, stdcli.CommandOptions{
		Usage:    "<app>",
		Validate: stdcli.Args(1),
	})

	e.Command("uninstall", "uninstall an application", e.Uninstall, stdcli.CommandOptions{
		Usage:    "<app>",
		Validate: stdcli.Args(1),
	})
}
