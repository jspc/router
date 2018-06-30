package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

const (
	exposePortVar   = "ROUTER_PORT"
	exposeSchemeVar = "ROUTER_SCHEME"
)

var (
	defaultScheme = "http"
)

type DockerClient interface {
	ContainerList(context.Context, types.ContainerListOptions) ([]types.Container, error)
	ContainerInspect(context.Context, string) (types.ContainerJSON, error)
}

type Docker struct {
	Client DockerClient
}

func NewDocker() (d Docker, err error) {
	d.Client, err = client.NewEnvClient()

	return
}

func (d Docker) GetContainerAddress(name string) (scheme *string, cn string, err error) {
	containers, err := d.Client.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return
	}

	var containerDetails types.ContainerJSON

	for _, container := range containers {
		for _, n := range container.Names {
			if n == fmt.Sprintf("/%s", name) {
				containerDetails, err = d.Client.ContainerInspect(context.Background(), container.ID)
				if err != nil {
					return
				}

				port := varsSearch(containerDetails.Config.Env, exposePortVar)
				if port == nil {
					err = fmt.Errorf("Service %q doesn't have a port exported to the router", name)

					return
				}

				scheme = varsSearch(containerDetails.Config.Env, exposeSchemeVar)
				if scheme == nil {
					scheme = &defaultScheme
				}

				cn = fmt.Sprintf("%s:%s",
					containerDetails.NetworkSettings.Networks["bridge"].IPAddress,
					*port,
				)

				return
			}
		}
	}

	err = fmt.Errorf("No container %q found", name)

	return
}

func varsSearch(v []string, key string) *string {
	for _, variable := range v {
		elems := strings.Split(variable, "=")

		if elems[0] == key {
			return &elems[1]
		}
	}

	return nil
}
