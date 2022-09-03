package docker

func Stop(id string) error {
	return docker("stop", id)
}
