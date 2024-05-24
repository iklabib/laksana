package containers

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"strings"
)

var TENANT_IMAGE = "quay.io/iklabib/markisa:tenant"

type ContainerClient struct {
	Engine string
}

func ContainerEngine() string {
	engine := os.Getenv("MARKISA_CONTAINER_ENGINE")
	if engine != "" {
		return engine
	}

	engine, _ = exec.LookPath("docker")
	if engine != "" {
		return engine
	}

	engine, _ = exec.LookPath("podman")
	if engine != "" {
		return engine
	}

	panic(errors.New("no container engine found in path"))
}

func NewDefault() *ContainerClient {
	return &ContainerClient{
		Engine: ContainerEngine(),
	}
}

func NewClient() *ContainerClient {
	return &ContainerClient{}
}

func (c ContainerClient) SpawnTenant() (string, error) {
	var stdoutBuff bytes.Buffer
	cmd := exec.Command(c.Engine, "run", "--rm", "-dt", TENANT_IMAGE)
	cmd.Stdout = &stdoutBuff
	err := cmd.Run()
	return stdoutBuff.String(), err
}

func (c ContainerClient) ExecTenant(id string, req []byte) ([]byte, error) {
	var stdoutBuff bytes.Buffer
	cmd := exec.Command(c.Engine, "exec", "-i", id, ".local/bin/commander")
	cmd.Stdin = bytes.NewBuffer(req)
	cmd.Stdout = &stdoutBuff
	err := cmd.Run()
	return stdoutBuff.Bytes(), err
}


func (c ContainerClient) CleanContainers() error {
	var stdoutBuff bytes.Buffer
	cmd := exec.Command(c.Engine, "ps", "-a", "-f", "ancestor=quay.io/iklabib/markisa:tenant", "--format", "'{{.ID}}'")
	cmd.Stdout = &stdoutBuff
	err := cmd.Run()
	if err != nil {
		return err
	}

	ids := strings.Split(stdoutBuff.String(), "\n")
	for _, id :=range ids {
		err := exec.Command(c.Engine, "rm", "-f", id).Run()
		if err != nil {
			return err
		}
	}
	return nil
}