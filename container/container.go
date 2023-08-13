package container

import (
	"context"
	"fmt"
	"os/user"

	"github.com/containers/common/libnetwork/types"
	"github.com/containers/podman/v4/pkg/bindings"
	"github.com/containers/podman/v4/pkg/bindings/containers"
	"github.com/containers/podman/v4/pkg/specgen"
)


func Connect() context.Context {
  user, err := user.Current()
  if err != nil {
    panic(err)
  }
  uri := fmt.Sprintf("unix://var/run/user/%s/podman/podman.sock", user.Uid)

	conn, err := bindings.NewConnection(context.Background(), uri)
	if err != nil {
    panic(err)
	}
  return conn
}

type Container struct {
  Connection context.Context
}

func Init() Container {
  return Container {
    Connection: Connect(),
  }
}

func (c Container) Create(spec string, hostPort, containerPort uint16) string {
	s := specgen.NewSpecGenerator(spec, false)
  s.NetNS = specgen.Namespace{
    NSMode: specgen.Slirp,
  }
  s.ContainerNetworkConfig.PortMappings = []types.PortMapping {
    {
      HostIP: "127.0.0.1",
      HostPort: hostPort,
      ContainerPort: containerPort,
    },
  }

	resp, err := containers.CreateWithSpec(c.Connection, s, nil)
	if err != nil {
    panic(err)
	}
  return resp.ID
}

func (c Container) Run(id string) {
  opts := &containers.StartOptions{}
  if err := containers.Start(c.Connection, id, opts); err != nil {
    panic(err)
  }
}

func (c Container) Stop(id string) {
  if err := containers.Stop(c.Connection, id, nil); err != nil{
    panic(err)
  }
}

func (c Container) Remove(id string) {
  _, err := containers.Remove(c.Connection, id, nil);
  if (err != nil) {
    panic(err)
  }
}


