package repository

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/epiphytelabs/keep/pkg/random"
	"gopkg.in/yaml.v3"
)

type Environment struct {
	app   *App
	names []string
	vars  map[string]EnvironmentVariable
}

type EnvironmentVariable interface {
	String() (string, error)
}

type EnvironmentVariableResource struct {
	app      *App
	name     string
	property string
}

type EnvironmentVariableSecret struct {
	app    *App
	secret string
}

type EnvironmentVariableStatic struct {
	string
}

func (e *Environment) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var w map[string]yaml.Node

	if err := unmarshal(&w); err != nil {
		return err
	}

	e.names = []string{}
	e.vars = map[string]EnvironmentVariable{}

	var err error

	for k, v := range w {
		switch v.Tag {
		case "!!int", "!!str":
			e.vars[k], err = e.static(v.Value)
		case "!resource":
			e.vars[k], err = e.resource(v.Value)
		case "!secret":
			e.vars[k], err = e.secret(v.Value)
		default:
			err = fmt.Errorf("unknown environment variable tag: %s", v.Tag)
		}

		if err != nil {
			return err
		}
	}

	for k := range e.vars {
		e.names = append(e.names, k)
	}

	sort.Strings(e.names)

	return nil
}

func (e *Environment) Get(name string) (string, error) {
	return e.vars[name].String()
}

func (e *Environment) Map() (map[string]string, error) {
	m := map[string]string{}

	for _, k := range e.Names() {
		v, err := e.Get(k)
		if err != nil {
			return nil, err
		}

		m[k] = v
	}

	return m, nil
}

func (e *Environment) Names() []string {
	return e.names
}

func (e *Environment) resource(args string) (*EnvironmentVariableResource, error) {
	parts := strings.SplitN(args, " ", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid args for !resource: %s", args)
	}

	v := &EnvironmentVariableResource{
		app:      e.app,
		name:     parts[0],
		property: parts[1],
	}

	return v, nil
}

func (e *Environment) static(args string) (*EnvironmentVariableStatic, error) {
	return &EnvironmentVariableStatic{args}, nil
}

func (e *Environment) secret(args string) (EnvironmentVariable, error) {
	size, err := strconv.Atoi(args)
	if err != nil {
		return nil, err
	}

	secret, err := random.String(size)
	if err != nil {
		return nil, err
	}

	v := &EnvironmentVariableSecret{
		app:    e.app,
		secret: secret,
	}

	return v, nil
}

func (v *EnvironmentVariableResource) String() (string, error) {
	r, err := v.app.Resources.Get(v.name)
	if err != nil {
		return "", err
	}

	u, err := r.URL()
	if err != nil {
		return "", err
	}

	switch v.property {
	case "host":
		return u.Hostname(), nil
	case "password":
		pw, _ := u.User.Password()
		return pw, nil
	case "path":
		return strings.TrimPrefix(u.Path, "/"), nil
	case "port":
		return u.Port(), nil
	case "username":
		return u.User.Username(), nil
	default:
		return "", fmt.Errorf("unknown resource part: %s", v.property)
	}
}

func (v *EnvironmentVariableSecret) String() (string, error) {
	return v.secret, nil
}

func (v *EnvironmentVariableStatic) String() (string, error) {
	return v.string, nil
}
