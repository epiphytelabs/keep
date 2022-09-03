package docker

import (
	"errors"
	"os/exec"
	"strings"
)

func docker(args ...string) error {
	cmd := exec.Command("docker", args...)

	if data, err := cmd.CombinedOutput(); err != nil {
		return errors.New(strings.TrimSpace(string(data)))
	}

	return nil
}

func dockerInfo(args ...string) ([]byte, error) {
	data, err := exec.Command("docker", args...).CombinedOutput()
	if err != nil {
		return nil, errors.New(strings.TrimSpace(string(data)))
	}

	return data, nil
}
