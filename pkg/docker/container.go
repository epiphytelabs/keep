package docker

type Container struct {
	Env      map[string]string
	Image    string
	Labels   map[string]string
	Links    map[string]string
	Name     string
	Networks []string
}

func (c Container) Create() error {
	return Create(c)
}

func (c Container) Pull() error {
	return Pull(c.Image)
}

func (c Container) Rm() error {
	return Rm(c.Name)
}

func (c Container) Start() error {
	return Start(c)
}

func (c Container) Stop() error {
	return Stop(c.Name)
}
