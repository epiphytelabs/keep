package docker

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
)

func Events(filters map[string]string) (<-chan string, error) {
	ch := make(chan string)

	args := []string{"events"}

	for k, v := range filters {
		args = append(args, "--filter", fmt.Sprintf("%s=%s", k, v))
	}

	cmd := exec.Command("docker", args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	go eventReader(stdout, ch)

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return ch, nil
}

func eventReader(r io.Reader, ch chan<- string) {
	s := bufio.NewScanner(r)

	for s.Scan() {
		ch <- s.Text()
	}
}
