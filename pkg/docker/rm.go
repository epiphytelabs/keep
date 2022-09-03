package docker

func Rm(name string) error {
	return docker("rm", name)
}
