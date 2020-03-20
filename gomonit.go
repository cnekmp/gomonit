package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/cnekmp/gomonit/client"
	"github.com/cnekmp/gomonit/server"
)

func main() {

	gomonitClient := flag.Bool("client", false, "Start gomonit CLIENT")
	gomonitServer := flag.Bool("server", false, "Start gomonit SERVER")
	flag.Parse()

	switch {
	case *gomonitServer:
		fmt.Println("Starting GOMonit SERVER...")
		server.Run()
	case *gomonitClient:
		fmt.Println("Starting GOMonit CLIENT...")
		client.Run()
	default:
		fmt.Println("No parameters provided. Please run with '-h' for details.")
		os.Exit(1)
	}
}
