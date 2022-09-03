package repository

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/epiphytelabs/keep"
	"github.com/epiphytelabs/keep/pkg/docker"
	"gopkg.in/yaml.v3"
)

type App struct {
	Environment Environment
	Image       string
	Name        string
	Port        int
	Resources   Resources
}

func (a *App) Container() (*docker.Container, error) {
	env, err := a.Environment.Map()
	if err != nil {
		return nil, err
	}

	labels := map[string]string{
		"system": "keep",
		"type":   "app",
		"app":    a.Name,
	}

	if a.Port > 0 {
		labels["port"] = fmt.Sprintf("%d", a.Port)
	}

	c := &docker.Container{
		Env:      env,
		Image:    a.Image,
		Labels:   labels,
		Links:    map[string]string{},
		Name:     fmt.Sprintf("keep-%s", a.Name),
		Networks: []string{"keep"},
	}

	for _, r := range a.Resources.List() {
		rc, err := r.Container()
		if err != nil {
			return nil, err
		}

		c.Links[rc.Name] = r.Name()
	}

	return c, nil
}

func apps() ([]App, error) {
	fs, err := keep.Apps.ReadDir("apps")
	if err != nil {
		return nil, err
	}

	as := []App{}

	for _, f := range fs {
		data, err := keep.Apps.ReadFile(filepath.Join("apps", f.Name()))
		if err != nil {
			return nil, err
		}

		var a App

		a.Environment = Environment{app: &a}
		a.Name = strings.TrimSuffix(filepath.Base(f.Name()), filepath.Ext(f.Name()))
		a.Resources = Resources{app: &a}

		if err := yaml.Unmarshal(data, &a); err != nil {
			return nil, err
		}

		as = append(as, a)
	}

	sort.Slice(as, func(i, j int) bool { return as[i].Name < as[j].Name })

	return as, nil
}
