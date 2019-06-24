package main

import (
	"os"
	"fmt"
	"github.com/daniellok/rps/client"
	"github.com/daniellok/rps/server"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("come on")
		return
	}

	if args[1] != "client" && args[1] != "server" {
		fmt.Println("client or server?")
	}

	if args[1] == "server" {
		server.RunServer()
	} else {
		client.RunClient()
	}
}
