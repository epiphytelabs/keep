package docker

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Container struct {
	Env      map[string]string
	Image    string
	Labels   map[string]string
	Links    map[string]string
	Name     string
	Networks []string
	Ports    []string
	Volumes  map[string]string
}

func (c Container) Create() error {
	return ContainerCreate(c)
}

func (c Container) Rm() error {
	return ContainerRm(c.Name)
}

func (c Container) Start() error {
	return ContainerStart(c)
}

func (c Container) Stop() error {
	return ContainerStop(c.Name)
}

func ContainerCreate(c Container) error {
	args := []string{
		"container",
		"create",
		"--name", c.Name,
		"--restart", "unless-stopped",
	}

	for k, v := range c.Env {
		args = append(args, "--env", fmt.Sprintf("%s=%s", k, v))
	}

	for k, v := range c.Labels {
		args = append(args, "--label", fmt.Sprintf("%s=%s", k, v))
	}

	for _, n := range c.Networks {
		args = append(args, "--network", n)
	}

	for _, p := range c.Ports {
		args = append(args, "--publish", p)
	}

	for local, remote := range c.Volumes {
		args = append(args, "--volume", fmt.Sprintf("%s:%s", local, remote))
	}

	args = append(args, c.Image)

	if err := docker(args...); err != nil {
		return err
	}

	return nil
}

func ContainerInspect(id string) (*Container, error) {
	data, err := dockerInfo("container", "inspect", id)
	if err != nil {
		return nil, err
	}

	var rcs []struct {
		Config struct {
			Env      []string
			Hostname string
			Image    string
			Labels   map[string]string
		}
		Name            string
		NetworkSettings struct {
			Networks map[string]struct {
				Aliases   []string
				IPAddress string
			}
		}
	}

	if err := json.Unmarshal(data, &rcs); err != nil {
		return nil, err
	}

	switch len(rcs) {
	case 0:
		return nil, fmt.Errorf("container not found: %s", id)
	case 1:
	default:
		return nil, fmt.Errorf("multiple containers found for: %s", id)
	}

	env := map[string]string{}

	for _, e := range rcs[0].Config.Env {
		if parts := strings.SplitN(e, "=", 2); len(parts) == 2 {
			env[parts[0]] = parts[1]
		}
	}

	ns := []string{}

	for k := range rcs[0].NetworkSettings.Networks {
		ns = append(ns, k)
	}

	c := &Container{
		Env:      env,
		Image:    rcs[0].Config.Image,
		Labels:   rcs[0].Config.Labels,
		Name:     rcs[0].Name[1:],
		Networks: ns,
	}

	return c, nil
}

func ContainerRm(name string) error {
	return docker("container", "rm", name)
}

func ContainerStart(c Container) error {
	return docker("container", "start", c.Name)
}

func ContainerStop(id string) error {
	return docker("container", "stop", id)
}
