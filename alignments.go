package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
)

// alignmentToolInput defines the input structure for the alignments tool.
type alignmentToolInput struct {
	Name string `json:"name" mcp:"description=The index of the alignment to retrieve (e.g., 'chaotic-good')."`
}

// alignmentListAPIResponse defines the structure for a single alignment in the list response.
type alignmentListAPIResponse struct {
	Index string `json:"index"`
	Name  string `json:"name"`
	URL   string `json:"url"`
}

// alignmentDetail defines the structure for a detailed alignment response.
type alignmentDetail struct {
	Index string `json:"index"`
	Name  string `json:"name"`
	Desc  string `json:"desc"`
	URL   string `json:"url"`
}

// alignmentToolOutput defines the output structure for the alignments tool.
type alignmentToolOutput struct {
	Count     int                        `json:"count,omitempty"`
	Results   []alignmentListAPIResponse `json:"results,omitempty"`
	Alignment *alignmentDetail           `json:"alignment,omitempty"`
}

// fetchAlignmentByNameResult fetches an alignment by index and returns an MCP tool result.
func fetchAlignmentByNameResult(
	client *http.Client,
	input alignmentToolInput,
	fetchByName func(*http.Client, endpoint, string, any) error,
) (*mcp.CallToolResult, error) {
	alignment := &alignmentDetail{}
	err := fetchByName(client, alignments, input.Name, alignment)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("failed to fetch alignment", err), err
	}
	output := alignmentToolOutput{Alignment: alignment}
	jsonData, err := json.Marshal(output)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to marshal alignment output", err), err
	}
	return mcp.NewToolResultText(string(jsonData)), nil
}

// fetchAlignmentListResult fetches a list of alignments and returns an MCP tool result.
func fetchAlignmentListResult(
	client *http.Client,
	fetchList func(*http.Client, endpoint, any, string) error,
) (*mcp.CallToolResult, error) {
	var results []alignmentListAPIResponse
	err := fetchList(client, alignments, &results, "")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to fetch alignment list", err), err
	}
	output := alignmentToolOutput{Count: len(results), Results: results}
	jsonData, err := json.Marshal(output)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to marshal alignment list output", err), err
	}
	return mcp.NewToolResultText(string(jsonData)), nil
}

// runAlignmentTool executes the core logic for the alignments tool.
func runAlignmentTool(
	input alignmentToolInput,
	fetchByName func(*http.Client, endpoint, string, any) error,
	fetchList func(*http.Client, endpoint, any, string) error,
) (*mcp.CallToolResult, error) {
	client := http.DefaultClient
	if input.Name != "" {
		return fetchAlignmentByNameResult(client, input, fetchByName)
	}
	return fetchAlignmentListResult(client, fetchList)
}

// handleAlignmentTool is the MCP handler for the alignments tool.
func handleAlignmentTool(ctx context.Context, req mcp.CallToolRequest, input alignmentToolInput) (*mcp.CallToolResult, error) {
	return runAlignmentTool(input, fetchByName, fetchList)
}
