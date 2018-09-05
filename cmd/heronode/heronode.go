package main

import (
	"context"
	//	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

const imageName = "ethereum/client-go"

func main() {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithVersion("1.38"))
	if err != nil {
		panic(err)
	}

	images, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	existed := false
	for _, image := range images {
		for _, name := range image.RepoTags {
			if strings.Contains(name, imageName) {
				existed = true
			}
		}
	}

	if !existed {
		reader, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
		if err != nil {
			panic(err)
		}
		defer reader.Close()
		io.Copy(os.Stdout, reader)
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}
	existed = false
	for _, c := range containers {
		if strings.Contains(c.Image, imageName) {
			existed = true
		}
	}

	if !existed {
		hostConfig := &container.HostConfig{
			PortBindings: nat.PortMap{
				"8545/udp": []nat.PortBinding{
					{
						HostIP:   "0.0.0,0",
						HostPort: "8545",
					},
				},
				"30303/udp": []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: "8545",
					},
				},
			},
		}
		resp, err := cli.ContainerCreate(ctx, &container.Config{
			Image: imageName,
			Cmd:   []string{"--rpc", "rpcaddr 0.0.0.0", "--ws", "--cache 1024", "--rpccorsdomain *"},
			Tty:   true,
			ExposedPorts: nat.PortSet{
				"8545/udp":  struct{}{},
				"30303/udp": struct{}{},
			},
		}, hostConfig, nil, "")
		if err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
			panic(err)
		}

		statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
		select {
		case err := <-errCh:
			if err != nil {
				panic(err)
			}
		case <-statusCh:
		}

		out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
		if err != nil {
			panic(err)
		}

		io.Copy(os.Stdout, out)
	}
}
