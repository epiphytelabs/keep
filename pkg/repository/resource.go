package repository

import (
	"fmt"
	"net/url"

	"github.com/epiphytelabs/keep/pkg/docker"
	"github.com/epiphytelabs/keep/pkg/random"
)

type Resource interface {
	Container() (*docker.Container, error)
	Name() string
	URL() (*url.URL, error)
}

type Resources struct {
	app       *App
	resources []Resource
}

func (a *App) resource(name string) (Resource, error) {
	switch name {
	case "postgres":
		return a.resourcePostgres()
	default:
		return nil, fmt.Errorf("unknown resource type: %s", name)
	}
}

func (rs Resources) Get(name string) (Resource, error) {
	for _, r := range rs.resources {
		if r.Name() == name {
			return r, nil
		}
	}

	return nil, fmt.Errorf("resource not found: %s", name)
}

func (rs Resources) List() []Resource {
	return rs.resources
}

func (rs *Resources) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var w []string

	if err := unmarshal(&w); err != nil {
		return err
	}

	rs.resources = []Resource{}

	for _, name := range w {
		r, err := rs.app.resource(name)
		if err != nil {
			return err
		}

		rs.resources = append(rs.resources, r)
	}

	return nil
}

type ResourcePostgres struct {
	app      *App
	password string
}

func (a *App) resourcePostgres() (*ResourcePostgres, error) {
	pw, err := random.String(32)
	if err != nil {
		return nil, err
	}

	r := &ResourcePostgres{app: a, password: pw}

	return r, nil
}

func (r *ResourcePostgres) Container() (*docker.Container, error) {
	c := &docker.Container{
		Name:  fmt.Sprintf("keep-%s-postgres", r.app.Name),
		Image: "postgres:14",
		Env: map[string]string{
			"POSTGRES_DB":       r.app.Name,
			"POSTGRES_PASSWORD": r.password,
			"POSTGRES_USER":     r.app.Name,
		},
	}

	return c, nil
}

func (r *ResourcePostgres) Name() string {
	return "postgres"
}

func (r *ResourcePostgres) URL() (*url.URL, error) {
	u := &url.URL{
		Path:   fmt.Sprintf("/%s", r.app.Name),
		User:   url.UserPassword(r.app.Name, r.password),
		Host:   "postgres:5432",
		Scheme: "postgres",
	}

	return u, nil
}
