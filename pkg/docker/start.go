package docker

func Start(c Container) error {
	return docker("start", c.Name)
}
