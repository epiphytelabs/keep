package cli

import (
	"github.com/ddollar/stdcli"
)

type Engine struct {
	*stdcli.Engine
}

func New(version string) *Engine {
	e := &Engine{Engine: stdcli.New("keep", version)}

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
