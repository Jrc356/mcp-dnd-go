package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
)

// backgroundToolInput defines the input structure for the backgrounds tool.
type backgroundToolInput struct {
	Name string `json:"name" mcp:"description=The index of the background to retrieve (e.g., 'acolyte')."`
}

// backgroundListAPIResponse defines the structure for a single background in the list response.
type backgroundListAPIResponse struct {
	Index string `json:"index"`
	Name  string `json:"name"`
	URL   string `json:"url"`
}

// backgroundDetail defines the structure for a detailed background response.
type backgroundDetail struct {
	Index string `json:"index"`
	Name  string `json:"name"`
	URL   string `json:"url"`
	// Add more fields as needed based on the API response
}

// backgroundToolOutput defines the output structure for the backgrounds tool.
type backgroundToolOutput struct {
	Count      int                         `json:"count,omitempty"`
	Results    []backgroundListAPIResponse `json:"results,omitempty"`
	Background *backgroundDetail           `json:"background,omitempty"`
}

// fetchBackgroundByNameResult fetches a background by index and returns an MCP tool result.
func fetchBackgroundByNameResult(
	client *http.Client,
	input backgroundToolInput,
	fetchByName func(*http.Client, endpoint, string, any) error,
) (*mcp.CallToolResult, error) {
	background := &backgroundDetail{}
	err := fetchByName(client, backgrounds, input.Name, background)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("failed to fetch background", err), err
	}
	output := backgroundToolOutput{Background: background}
	jsonData, err := json.Marshal(output)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to marshal background output", err), err
	}
	return mcp.NewToolResultText(string(jsonData)), nil
}

// fetchBackgroundListResult fetches a list of backgrounds and returns an MCP tool result.
func fetchBackgroundListResult(
	client *http.Client,
	fetchList func(*http.Client, endpoint, any, string) error,
) (*mcp.CallToolResult, error) {
	var results []backgroundListAPIResponse
	err := fetchList(client, backgrounds, &results, "")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to fetch background list", err), err
	}
	output := backgroundToolOutput{Count: len(results), Results: results}
	jsonData, err := json.Marshal(output)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to marshal background list output", err), err
	}
	return mcp.NewToolResultText(string(jsonData)), nil
}

// runBackgroundTool executes the core logic for the backgrounds tool.
func runBackgroundTool(
	input backgroundToolInput,
	fetchByName func(*http.Client, endpoint, string, any) error,
	fetchList func(*http.Client, endpoint, any, string) error,
) (*mcp.CallToolResult, error) {
	client := http.DefaultClient
	if input.Name != "" {
		return fetchBackgroundByNameResult(client, input, fetchByName)
	}
	return fetchBackgroundListResult(client, fetchList)
}

// handleBackgroundTool is the MCP handler for the backgrounds tool.
func handleBackgroundTool(ctx context.Context, req mcp.CallToolRequest, input backgroundToolInput) (*mcp.CallToolResult, error) {
	return runBackgroundTool(input, fetchByName, fetchList)
}
