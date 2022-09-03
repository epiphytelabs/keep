package docker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
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

	args = []string{"inspect"}

	for _, id := range ids {
		args = append(args, string(id))
	}

	data, err = dockerInfo(args...)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &rcs); err != nil {
		return nil, err
	}

	cs := []Container{}

	for _, rc := range rcs {
		env := map[string]string{}

		for _, e := range rc.Config.Env {
			if parts := strings.SplitN(e, "=", 2); len(parts) == 2 {
				env[parts[0]] = parts[1]
			}
		}

		ns := []string{}

		for k := range rc.NetworkSettings.Networks {
			ns = append(ns, k)
		}

		cs = append(cs, Container{
			Env:      env,
			Image:    rc.Config.Image,
			Labels:   rc.Config.Labels,
			Name:     rc.Name[1:],
			Networks: ns,
		})
	}

	return cs, nil
}
