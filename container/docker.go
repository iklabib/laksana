package container

// https://pkg.go.dev/github.com/docker/docker/client
import (
	"github.com/docker/docker/client"
)

func CreateClient() *client.Client {
    docker_api_ver := client.WithVersion("1.43")
    docker_host := client.WithHost("unix:///run/user/1000/podman/podman.sock")

    cli, err := client.NewClientWithOpts(docker_host, docker_api_ver)
    if err != nil {
        panic(err)
    }
    return cli
}
