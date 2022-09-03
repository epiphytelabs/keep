package docker

import "fmt"

type Network struct {
	Name string
}

func NetworkCreate(name string) (*Network, error) {
	if err := docker("network", "create", name); err != nil {
		return nil, err
	}

	n := &Network{
		Name: name,
	}

	return n, nil
}

func NetworkRemove(name string) error {
	return docker("network", "rm", name)
}

func (n Network) Connect(c Container) error {
	args := []string{
		"network",
		"connect",
	}

	if c.Links != nil {
		for k, v := range c.Links {
			args = append(args, "--link", fmt.Sprintf("%s:%s", k, v))
		}
	}

	args = append(args, n.Name, c.Name)

	return docker(args...)
}
