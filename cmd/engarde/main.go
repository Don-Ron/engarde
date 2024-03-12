package main

import (
	"os"

	"github.com/kamushadenes/engarde/v2"
)

// Version is passed by the compiler
var Version = "UNOFFICIAL BUILD"

func printVersion() {
	if Version != "" {
		print("engarde-client ver. " + Version + "\r\n")
	}
}

func printUsage() {
	if _, err := os.Stderr.WriteString("Usage: engarde <server|client> [config_file]\n"); err != nil {
		panic(err)
	}
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
	}

	configName := "engarde.yml"

	if len(os.Args) > 2 {
		configName = os.Args[2]
	}

	printVersion()

	switch os.Args[1] {
	case "server":
		engarde.RunServer(configName)
	case "client":
		engarde.RunClient(configName)
	case "list-interfaces":
		engarde.ListInterfaces()
	}
}
