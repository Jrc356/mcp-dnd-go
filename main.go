package main

import (
	"fmt"
	"net/http"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer(
		"D&D 5e Knowledge Base",
		"1.0.0",
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)

	httpClient := http.DefaultClient
	apiBaseURL := APIBaseURL
	registerTools(s, httpClient, apiBaseURL)

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
