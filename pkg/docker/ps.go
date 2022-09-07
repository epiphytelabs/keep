package docker

import (
	"bytes"
	"fmt"
)

func Ps(labels map[string]string) ([]Container, error) {
	args := []string{"ps", "-q"}

	for k, v := range labels {
		args = append(args, "--filter", fmt.Sprintf("label=%s=%s", k, v))
	}

	data, err := dockerInfo(args...)
	if err != nil {
		return nil, err
	}

	tdata := bytes.TrimSpace(data)

	if len(tdata) == 0 {
		return []Container{}, nil
	}

	ids := bytes.Split(tdata, []byte{'\n'})

	cs := []Container{}

	for _, id := range ids {
		c, err := ContainerInspect(string(id))
		if err != nil {
			return nil, err
		}

		cs = append(cs, *c)
	}

	return cs, nil
}
