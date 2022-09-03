package docker

import "fmt"

func Create(c Container) error {
	args := []string{
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

	// if c.Ports != nil {
	// 	for k, v := range c.Ports {
	// 		args = append(args, "--publish", fmt.Sprintf("%d:%d", k, v))
	// 	}
	// }

	args = append(args, c.Image)

	if err := docker(args...); err != nil {
		return err
	}

	return nil
}
