package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
)

// classToolInput defines the input structure for the classes tool.
type classToolInput struct {
	Name string `json:"name" mcp:"description=The index of the class to retrieve (e.g., 'barbarian')."`
}

// classListAPIResponse defines the structure for a single class in the list response.
type classListAPIResponse struct {
	Index string `json:"index"`
	Name  string `json:"name"`
	URL   string `json:"url"`
}

// classDetail defines the structure for a detailed class response.
type classDetail struct {
	Index                    string      `json:"index"`
	Name                     string      `json:"name"`
	HitDie                   int         `json:"hit_die"`
	ProficiencyChoices       interface{} `json:"proficiency_choices"`
	Proficiencies            interface{} `json:"proficiencies"`
	SavingThrows             interface{} `json:"saving_throws"`
	StartingEquipment        interface{} `json:"starting_equipment"`
	StartingEquipmentOptions interface{} `json:"starting_equipment_options"`
	ClassLevels              string      `json:"class_levels"`
	MultiClassing            interface{} `json:"multi_classing"`
	Subclasses               interface{} `json:"subclasses"`
	URL                      string      `json:"url"`
	UpdatedAt                string      `json:"updated_at"`
	// Add more fields as needed based on the API response
}

// classToolOutput defines the output structure for the classes tool.
type classToolOutput struct {
	Count   int                    `json:"count,omitempty"`
	Results []classListAPIResponse `json:"results,omitempty"`
	Class   *classDetail           `json:"class,omitempty"`
}

// fetchClassByNameResult fetches a class by index and returns an MCP tool result.
func fetchClassByNameResult(
	client *http.Client,
	input classToolInput,
	fetchByName func(*http.Client, endpoint, string, any) error,
) (*mcp.CallToolResult, error) {
	class := &classDetail{}
	err := fetchByName(client, classes, input.Name, class)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("failed to fetch class", err), err
	}
	output := classToolOutput{Class: class}
	jsonData, err := json.Marshal(output)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to marshal class output", err), err
	}
	return mcp.NewToolResultText(string(jsonData)), nil
}

// fetchClassListResult fetches a list of classes and returns an MCP tool result.
func fetchClassListResult(
	client *http.Client,
	fetchList func(*http.Client, endpoint, any, string) error,
) (*mcp.CallToolResult, error) {
	var results []classListAPIResponse
	err := fetchList(client, classes, &results, "")
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to fetch class list", err), err
	}
	output := classToolOutput{Count: len(results), Results: results}
	jsonData, err := json.Marshal(output)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to marshal class list output", err), err
	}
	return mcp.NewToolResultText(string(jsonData)), nil
}

// runClassTool executes the core logic for the classes tool.
func runClassTool(
	input classToolInput,
	fetchByName func(*http.Client, endpoint, string, any) error,
	fetchList func(*http.Client, endpoint, any, string) error,
) (*mcp.CallToolResult, error) {
	client := http.DefaultClient
	if input.Name != "" {
		return fetchClassByNameResult(client, input, fetchByName)
	}
	return fetchClassListResult(client, fetchList)
}

// handleClassTool is the MCP handler for the classes tool.
func handleClassTool(ctx context.Context, req mcp.CallToolRequest, input classToolInput) (*mcp.CallToolResult, error) {
	return runClassTool(input, fetchByName, fetchList)
}
