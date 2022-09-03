package main

import (
	"fmt"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/epiphytelabs/keep/pkg/docker"
	"github.com/epiphytelabs/keep/pkg/server"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	rs, err := routes()
	if err != nil {
		return err
	}

	s := server.New(rs)

	go update(s, "start")
	go update(s, "stop")

	if err := s.Listen(":80"); err != nil {
		return err
	}

	return nil
}

func routes() (server.Routes, error) {
	cs, err := docker.Ps(map[string]string{"system": "keep", "type": "app"})
	if err != nil {
		return nil, err
	}

	routes := server.Routes{}

	for _, c := range cs {
		if port := c.Labels["port"]; port != "" {
			remote := &url.URL{Scheme: "http", Host: fmt.Sprintf("%s:%s", c.Name, port)}
			routes[fmt.Sprintf("%s.keep", c.Labels["app"])] = httputil.NewSingleHostReverseProxy(remote)
		}
	}

	return routes, nil
}

func update(s *server.Server, event string) {
	ch, err := docker.Events(map[string]string{"type": "container", "event": event})
	if err != nil {
		fmt.Printf("err: %+v\n", err)
	}

	for range ch {
		if rs, err := routes(); err != nil {
			fmt.Printf("err: %+v\n", err)
		} else {
			s.UpdateRoutes(rs)
		}
	}
}
