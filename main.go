package main

import (
	"context"
	"fmt"
	"github.com/ServiceWeaver/weaver"
	"github.com/lemon-1997/weaver/api"
	"os"
)

//go:generate weaver generate ./...

func main() {
	ctx := context.Background()
	root := weaver.Init(ctx)
	server, err := api.NewServer(root)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating server: ", err)
		os.Exit(1)
	}
	if err = server.Run("localhost:12345"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
