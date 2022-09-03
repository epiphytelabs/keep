package docker

func Pull(image string) error {
	return docker("pull", image)
}
