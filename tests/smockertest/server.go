package smockertest

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

const (
	image = "thiht/smocker:0.18.2"

	port      = 8080
	adminPort = 8081
)

type containerID struct {
	id string
}

func MustStart() *containerID {
	container, err := Start()
	if err != nil {
		log.Fatalln(err)
	}
	return container
}

func Start() (*containerID, error) {
	cmd := exec.Command(
		"docker", "run", "-d", "-p", fmt.Sprintf("%d:%d", port, port), "-p", fmt.Sprintf("%d:%d", adminPort, adminPort), "--name", "dyndns-smocker", image,
	)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("%v%v", stderr.String(), err)
	}

	id := strings.TrimSpace(stdout.String())
	if id == "" {
		return nil, errors.New("unexpected empty output from `docker run`")
	}

	return &containerID{id}, nil
}

func (c containerID) Nuke() {
	if err := c.Kill(); err != nil {
		log.Fatalln(err)
	}

	if err := c.Remove(); err != nil {
		log.Fatalln(err)
	}
}

func (c containerID) Kill() error {
	return exec.Command("docker", "kill", c.id).Run()
}

func (c containerID) Remove() error {
	return exec.Command("docker", "rm", c.id).Run()
}
