package main

import (
	"fmt"
)

type StatusCLI struct {
	Version StatusVersionCmd `kong:"cmd,help='get the Aptly server version'"`
	Storage StatusStorageCmd `kong:"cmd,help='get used information how full the storage is, since Aptly 1.6.0'"`
}

type StatusVersionCmd struct{}

func (c *StatusVersionCmd) Run(ctx *Context) error {
	ver, err := ctx.client.Version()
	if err != nil {
		return err
	}
	fmt.Printf("Aptly Server: %s\n", ver)
	return nil
}

type StatusStorageCmd struct{}

func (c *StatusStorageCmd) Run(ctx *Context) error {
	storage, err := ctx.client.StorageUsage()
	if err != nil {
		return err
	}
	fmt.Printf("Total: %d MiB\n", storage.Total)
	fmt.Printf("Free: %d MiB\n", storage.Free)
	fmt.Printf("Percent: %f%%\n", storage.PercentFull)
	return nil
}
