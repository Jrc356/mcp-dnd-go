package main

import (
	"fmt"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer(
		"D&D 5e Knowledge Base",
		"1.0.0",
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)

	RegisterTools(s)

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
