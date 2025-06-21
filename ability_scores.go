package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
)

// abilityScoreToolInput defines the input structure for the ability-scores tool.
type abilityScoreToolInput struct {
	Name string `json:"name" mcp:"description=The index of the ability score to retrieve (e.g., 'str', 'dex')."`
}

// abilityScoreListAPIResponse defines the structure for a single ability score in the list response.
type abilityScoreListAPIResponse struct {
	Index string `json:"index"`
	Name  string `json:"name"`
	URL   string `json:"url"`
}

// abilityScoreDetail defines the structure for a detailed ability score response.
type abilityScoreDetail struct {
	Index    string   `json:"index"`
	Name     string   `json:"name"`
	FullName string   `json:"full_name"`
	Desc     []string `json:"desc"`
	Skills   []struct {
		Index string `json:"index"`
		Name  string `json:"name"`
		URL   string `json:"url"`
	} `json:"skills"`
	URL string `json:"url"`
}

// abilityScoreToolOutput defines the output structure for the ability-scores tool.
type abilityScoreToolOutput struct {
	Count        int                           `json:"count,omitempty"`
	Results      []abilityScoreListAPIResponse `json:"results,omitempty"`
	AbilityScore *abilityScoreDetail           `json:"ability_score,omitempty"`
}

// fetchAbilityScoreByNameResult fetches an ability score by index and returns an MCP tool result.
func fetchAbilityScoreByNameResult(
	client *http.Client,
	input abilityScoreToolInput,
	fetchByName func(*http.Client, endpoint, string, any) error,
) (*mcp.CallToolResult, error) {
	abilityScore := &abilityScoreDetail{}
	err := fetchByName(client, abilityScores, input.Name, abilityScore)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("failed to fetch ability score", err), err
	}
	output := abilityScoreToolOutput{AbilityScore: abilityScore}
	jsonData, err := json.Marshal(output)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to marshal ability score output", err), err
	}
	return mcp.NewToolResultText(string(jsonData)), nil
}

// fetchAbilityScoreListResult fetches a list of ability scores and returns an MCP tool result.
func fetchAbilityScoreListResult(
	client *http.Client,
	fetchList func(*http.Client, endpoint, any, string) error,
) (*mcp.CallToolResult, error) {
	var results []abilityScoreListAPIResponse
	err := fetchList(client, abilityScores, &results, "")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to fetch ability score list", err), err
	}
	output := abilityScoreToolOutput{Count: len(results), Results: results}
	jsonData, err := json.Marshal(output)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to marshal ability score list output", err), err
	}
	return mcp.NewToolResultText(string(jsonData)), nil
}

// runAbilityScoreTool executes the core logic for the ability-scores tool.
func runAbilityScoreTool(
	input abilityScoreToolInput,
	fetchByName func(*http.Client, endpoint, string, any) error,
	fetchList func(*http.Client, endpoint, any, string) error,
) (*mcp.CallToolResult, error) {
	client := http.DefaultClient
	if input.Name != "" {
		return fetchAbilityScoreByNameResult(client, input, fetchByName)
	}
	return fetchAbilityScoreListResult(client, fetchList)
}

// handleAbilityScoreTool is the MCP handler for the ability-scores tool.
func handleAbilityScoreTool(ctx context.Context, req mcp.CallToolRequest, input abilityScoreToolInput) (*mcp.CallToolResult, error) {
	return runAbilityScoreTool(input, fetchByName, fetchList)
}
