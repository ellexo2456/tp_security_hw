package certs

import (
	"os"
	"os/exec"
	"strings"
)

func Load() (map[string][]byte, error) {
	files, err := os.ReadDir("./src/certs/list")
	if err != nil {
		return nil, err
	}

	certs := make(map[string][]byte, len(files))

	for _, f := range files {
		host := strings.TrimSuffix(f.Name(), ".crt")

		if certs[host], err = os.ReadFile("./src/certs/list/" + f.Name()); err != nil {
			return nil, err
		}
	}

	return certs, nil
}

func Generate(host string) ([]byte, error) {
	c, err := exec.Command("./src/certs/gen.sh", host).Output()
	if err != nil {
		return nil, err
	}

	err = os.WriteFile("./src/certs/list/"+host+".crt", c, 0644)
	if err != nil {
		return nil, err
	}

	return c, nil
}
