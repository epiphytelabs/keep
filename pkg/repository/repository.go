package repository

import (
	"fmt"
)

func Get(name string) (*App, error) {
	as, err := apps()
	if err != nil {
		return nil, err
	}

	for _, a := range as {
		if a.Name == name {
			return &a, nil
		}
	}

	return nil, fmt.Errorf("app not found: %s", name)
}

func Installed() ([]string, error) {
	return installed()
}

func List() ([]App, error) {
	return []App{}, nil
}
