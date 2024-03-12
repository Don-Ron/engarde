package main

import (
	"os"
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
		server(configName)
	case "client":
		client(configName)
	case "list-interfaces":
		listInterfaces()
	}
}
