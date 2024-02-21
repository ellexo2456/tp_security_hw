package certs

import (
	"math/rand"
	"os"
	"os/exec"
	"strconv"
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
	c, err := exec.Command("./gen_cert.sh", host, strconv.Itoa(rand.Intn(1000000000000))).Output()
	if err != nil {
		return nil, err
	}

	err = os.WriteFile("./src/certs/list/"+host+".crt", c, 0644)
	if err != nil {
		return nil, err
	}

	return c, nil
}
