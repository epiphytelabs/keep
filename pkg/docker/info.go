package docker

func Info() error {
	return docker("info")
}
