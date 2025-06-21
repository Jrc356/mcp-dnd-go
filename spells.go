package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"
)

// spellToolInput defines the input structure for the spell tool.
type spellToolInput struct {
	Name   string `json:"name" mcp:"description=The name of the spell to retrieve."`
	Level  int    `json:"level" mcp:"description=The level of the spell."`
	School string `json:"school" mcp:"description=The school of magic the spell belongs to."`
}

// buildQueryString constructs a query string from the spellToolInput fields for use in API requests.
func (f *spellToolInput) buildQueryString() string {
	if f == nil {
		return ""
	}
	params := []string{}
	if f.Level != 0 {
		params = append(params, "level="+strconv.Itoa(f.Level))
	}
	if f.School != "" {
		params = append(params, "school="+f.School)
	}
	if len(params) == 0 {
		return ""
	}
	return strings.Join(params, "&")
}

// spellToolOutput defines the output structure for the spell tool.
type spellToolOutput struct {
	Count   int                    `json:"count,omitempty"`
	Results []spellListAPIResponse `json:"results,omitempty"`
	Spell   *spellAPIResponse      `json:"spell,omitempty"`
}

// spellListAPIResponse defines the structure for a single spell in the list response.
type spellListAPIResponse struct {
	Index string `json:"index"`
	Name  string `json:"name"`
	Level int    `json:"level"`
	URL   string `json:"url"`
}

// spellAPIResponse defines the structure for a detailed spell response.
type spellAPIResponse struct {
	Index         string   `json:"index"`
	Name          string   `json:"name"`
	Desc          []string `json:"desc"`
	Range         string   `json:"range"`
	Components    []string `json:"components"`
	Ritual        bool     `json:"ritual"`
	Duration      string   `json:"duration"`
	Concentration bool     `json:"concentration"`
	CastingTime   string   `json:"casting_time"`
	Level         int      `json:"level"`
	DC            struct {
		DCType struct {
			Index string `json:"index"`
			Name  string `json:"name"`
			URL   string `json:"url"`
		} `json:"dc_type"`
		DCSuccess string `json:"dc_success"`
	} `json:"dc"`
	AreaOfEffect struct {
		Type string `json:"type"`
		Size int    `json:"size"`
	} `json:"area_of_effect"`
	School struct {
		Index string `json:"index"`
		Name  string `json:"name"`
		URL   string `json:"url"`
	} `json:"school"`
	Classes []struct {
		Index string `json:"index"`
		Name  string `json:"name"`
		URL   string `json:"url"`
	} `json:"classes"`
	Subclasses []struct {
		Index string `json:"index"`
		Name  string `json:"name"`
		URL   string `json:"url"`
	} `json:"subclasses"`
	URL       string `json:"url"`
	UpdatedAt string `json:"updated_at"`
}

// fetchSpellByNameResult handles the logic for fetching a spell by name and returning an MCP tool result.
func fetchSpellByNameResult(
	client *http.Client,
	input spellToolInput,
	fetchByName func(*http.Client, endpoint, string, any) error,
) (*mcp.CallToolResult, error) {
	spell := &spellAPIResponse{}
	err := fetchByName(client, spells, input.Name, spell)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("failed to fetch name", err), err
	}
	output := spellToolOutput{Spell: spell}
	jsonData, err := json.Marshal(output)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to unmarshal", err), err
	}
	return mcp.NewToolResultText(string(jsonData)), nil
}

// fetchSpellListResult handles the logic for fetching a list of spells and returning an MCP tool result.
func fetchSpellListResult(
	client *http.Client,
	input spellToolInput,
	fetchList func(*http.Client, endpoint, any, string) error,
) (*mcp.CallToolResult, error) {
	var results []spellListAPIResponse
	err := fetchList(client, spells, &results, input.buildQueryString())
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to fetch spell list", err), err
	}
	output := spellToolOutput{Count: len(results), Results: results}
	jsonData, err := json.Marshal(output)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to marshal output", err), err
	}
	return mcp.NewToolResultText(string(jsonData)), nil
}

// runSpellTool executes the core logic for the spell tool, using injected fetchByName and fetchList dependencies for testability.
// It returns an MCP tool result and a Go error if one occurs.
func runSpellTool(
	input spellToolInput,
	fetchByName func(*http.Client, endpoint, string, any) error,
	fetchList func(*http.Client, endpoint, any, string) error,
) (*mcp.CallToolResult, error) {
	client := http.DefaultClient
	if input.Name != "" {
		return fetchSpellByNameResult(client, input, fetchByName)
	}
	return fetchSpellListResult(client, input, fetchList)
}

// handleSpellTool is the MCP handler for the spell tool. It dispatches to the appropriate fetch function.
func handleSpellTool(ctx context.Context, req mcp.CallToolRequest, input spellToolInput) (*mcp.CallToolResult, error) {
	logrus.WithFields(logrus.Fields{"input": input}).Debug("handleSpellTool called")
	return runSpellTool(input, fetchByName, fetchList)
}
