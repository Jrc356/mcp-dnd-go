package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterTools(s *server.MCPServer) {
	client := http.DefaultClient

	s.AddTool(mcp.NewTool("get_categories",
		mcp.WithDescription("List all available D&D 5e API categories for browsing the official content."),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := FetchCategories(client, APIBaseURL)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		jsonBytes, err := json.Marshal(result)
		if err != nil {
			return mcp.NewToolResultError("Failed to marshal result: " + err.Error()), nil
		}
		return mcp.NewToolResultText(string(jsonBytes)), nil
	})

	s.AddTool(mcp.NewTool("get_items",
		mcp.WithDescription("Retrieve a list of all items available in a specific D&D 5e API category."),
		mcp.WithString("category", mcp.Required(), mcp.Description("The D&D API category to retrieve items from (e.g., 'spells', 'monsters', 'equipment')")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		category, err := req.RequireString("category")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		if _, ok := CategoryDescriptions[category]; !ok {
			return mcp.NewToolResultError("Invalid category: " + category), nil
		}
		result, err := FetchItems(client, APIBaseURL, category)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		jsonBytes, err := json.Marshal(result)
		if err != nil {
			return mcp.NewToolResultError("Failed to marshal result: " + err.Error()), nil
		}
		return mcp.NewToolResultText(string(jsonBytes)), nil
	})

	s.AddTool(mcp.NewTool("get_item",
		mcp.WithDescription("Retrieve detailed information about a specific D&D 5e item by its category and index."),
		mcp.WithString("category", mcp.Required(), mcp.Description("The D&D API category the item belongs to (e.g., 'spells', 'monsters', 'equipment')")),
		mcp.WithString("index", mcp.Required(), mcp.Description("The unique identifier for the specific item (e.g., 'fireball', 'adult-red-dragon')")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		category, err := req.RequireString("category")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		index, err := req.RequireString("index")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, err := FetchItem(client, APIBaseURL, category, index)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		jsonBytes, err := json.Marshal(result)
		if err != nil {
			return mcp.NewToolResultError("Failed to marshal result: " + err.Error()), nil
		}
		return mcp.NewToolResultText(string(jsonBytes)), nil
	})
}
