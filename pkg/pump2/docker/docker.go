package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"os"
)

func ImageRemove(ImageName string, Force bool) ([]types.ImageDeleteResponseItem, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithVersion("1.39"))
	if err != nil {
		panic(err)
	}
	out, err := cli.ImageRemove(ctx, ImageName, types.ImageRemoveOptions{Force: Force})
	if err != nil {
		panic(err)
	}
	return out, err
}

func ImageBuild(ImageName string) (types.ImageBuildResponse, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.WithVersion("1.39"))
	if err != nil {
		panic(err)
	}
	bctx, err := os.Open("/path/to/tar/archieve.tar")
	out, err := cli.ImageBuild(ctx, bctx, types.ImageBuildOptions{ImageName:ImageName})
	if err != nil {
		panic(err)
	}
	return out, err
}