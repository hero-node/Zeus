package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/mitchellh/go-homedir"
)

const eth_imageName = "ethereum/client-go"
const ipfs_imageName = "ipfs/go-ipfs"

func main() {

	stop := false
	if len(os.Args) > 1 {
		if os.Args[1] == "stop" {
			stop = true
		}
	}
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithVersion("1.38"))
	if err != nil {
		panic(err)
	}

	images, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	eth_existed := false
	ipfs_existed := false
	for _, image := range images {
		for _, name := range image.RepoTags {
			if strings.Contains(name, eth_imageName) {
				eth_existed = true
			}
			if strings.Contains(name, ipfs_imageName) {
				ipfs_existed = true
			}
		}
	}

	if !eth_existed {
		reader, err := cli.ImagePull(ctx, eth_imageName, types.ImagePullOptions{})
		if err != nil {
			panic(err)
		}
		defer reader.Close()
		io.Copy(os.Stdout, reader)
	}

	if !ipfs_existed {
		reader, err := cli.ImagePull(ctx, ipfs_imageName, types.ImagePullOptions{})
		if err != nil {
			panic(err)
		}
		defer reader.Close()
		io.Copy(os.Stdout, reader)
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All:true})
	if err != nil {
		panic(err)
	}
	eth_container := types.Container{ID:"empty"}
	ipfs_container := types.Container{ID:"empty"}
	for _, c := range containers {
		if strings.Contains(c.Image, eth_imageName) {
			eth_container = c
			continue
		}
		if strings.Contains(c.Image, ipfs_imageName) {
			ipfs_container = c
			continue
		}
	}
	if !stop {
		if eth_container.ID == "empty" {
			hostConfig := &container.HostConfig{
				PortBindings: nat.PortMap{
					"8545/tcp": []nat.PortBinding{
						{
							HostIP:   "0.0.0.0",
							HostPort: "8545",
						},
					},
					"30303/tcp": []nat.PortBinding{
						{
							HostIP:   "0.0.0.0",
							HostPort: "30303",
						},
					},
				},
			}
			resp, err := cli.ContainerCreate(ctx, &container.Config{
				Image: eth_imageName,
				Cmd:   []string{"--rpc", "--cache=768", "--maxpeers=128", "--rpcaddr=0.0.0.0", "--ws", "--rpccorsdomain=*"},
				Tty:   true,
				ExposedPorts: nat.PortSet{
					"8545/tcp":  struct{}{},
					"30303/tcp": struct{}{},
				},
			}, hostConfig, nil, "")
			if err != nil {
				panic(err)
			}
			fmt.Println("ethID:", resp.ID)
			if err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
				panic(err)
			}

			out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
			if err != nil {
				panic(err)
			}

			io.Copy(os.Stdout, out)
			fmt.Println("----------- Eth node started ------------")
		} else {
			if strings.Contains(eth_container.Status, "Exited") {
				if err := cli.ContainerStart(ctx, eth_container.ID, types.ContainerStartOptions{}); err != nil {
					panic(err)
				}
				fmt.Println("ethID:", eth_container.ID)
				fmt.Println("----------- Eth node started ------------")
			}
		}

		if ipfs_container.ID == "empty" {
			homeDir, err := homedir.Dir()
			if err != nil {
				panic(err)
			}

			ipfsPath := filepath.Join(homeDir, "ipfs")
			stagePath := filepath.Join(ipfsPath, "staging")
			dataPath := filepath.Join(ipfsPath, "data")

			hostConfig := &container.HostConfig{
				PortBindings: nat.PortMap{
					"8080/tcp": []nat.PortBinding{
						{
							HostIP:   "0.0.0.0",
							HostPort: "8080",
						},
					},
					"4001/tcp": []nat.PortBinding{
						{
							HostIP:   "0.0.0.0",
							HostPort: "4001",
						},
					},
					"5001/tcp": []nat.PortBinding{
						{
							HostIP:   "0.0.0.0",
							HostPort: "5001",
						},
					},
				},
				Binds: []string{stagePath + ":/export", dataPath + ":/data/ipfs"},
			}
			resp, err := cli.ContainerCreate(ctx, &container.Config{
				Image: ipfs_imageName,
				Tty:   true,
				ExposedPorts: nat.PortSet{
					"8080/tcp": struct{}{},
					"4001/tcp": struct{}{},
					"5001/tcp": struct{}{},
				},
			}, hostConfig, nil, "")
			fmt.Println("ipfsID:", resp.ID)
			if err != nil {
				panic(err)
			}

			if err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
				panic(err)
			}

			out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
			if err != nil {
				panic(err)
			}

			io.Copy(os.Stdout, out)

			fmt.Println("----------- Ipfs node started ------------")
		} else {
			if strings.Contains(ipfs_container.Status, "Exited") {
				if err := cli.ContainerStart(ctx, ipfs_container.ID, types.ContainerStartOptions{}); err != nil {
					panic(err)
				}
				fmt.Println("ipfsID:", ipfs_container.ID)
				fmt.Println("----------- Ipfs node started ------------")
			}
		}
	}

	if stop {
		if eth_container.ID != "empty" {
			if err := cli.ContainerStop(ctx, eth_container.ID, nil); err != nil {
				panic(err)
			}
		}
		if ipfs_container.ID != "empty" {
			if err = cli.ContainerStop(ctx, ipfs_container.ID, nil); err != nil {
				panic(err)
			}
		}
		fmt.Println("Stop Hero Node success")
	}
}
