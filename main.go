package main

import (
	"context"
	"reflect"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"
)

// newAPITool creates a new MCP tool for a given endpoint with the specified input type and handler.
func newAPITool[T any](
	e endpoint,
	description string,
	input T,
	handler func(ctx context.Context, req mcp.CallToolRequest, input T) (*mcp.CallToolResult, error),
) server.ServerTool {
	logrus.WithFields(logrus.Fields{
		"endpoint":    e,
		"description": description,
		"inputType":   reflect.TypeOf(input),
	}).Debug("Creating new API tool")
	opts := []mcp.ToolOption{
		mcp.WithDescription(description),
	}
	opts = append(opts, structToToolOptions(input)...)
	tool := mcp.NewTool(string(e), opts...)
	logrus.Debugf("Tool Input Schema Properties: %v", tool.InputSchema.Properties)
	return server.ServerTool{
		Tool:    tool,
		Handler: mcp.NewTypedToolHandler(handler),
	}
}

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	logrus.SetLevel(logrus.DebugLevel)
	logrus.Info("Starting D&D 5e MCP server...")

	s := server.NewMCPServer(
		"D&D 5e Knowledge Base",
		"1.0.0",
		server.WithToolCapabilities(true),
		server.WithRecovery(),
	)

	logrus.Info("Creating tools...")
	tools := []server.ServerTool{
		newAPITool(
			spells,
			"fetches information about d&d 5e spells.",
			spellToolInput{
				// a value is needed to generate the schema
				Name:   "spells",
				Filter: &spellFilter{},
			},
			handleSpellTool,
		),
		newAPITool(
			monsters,
			"fetches information about d&d 5e monsters.",
			monsterToolInput{
				// a value is needed to generate the schema
				Name:   "monsters",
				Filter: &monsterFilter{},
			},
			handleMonsterTool,
		),
	}
	for _, tool := range tools {
		logrus.WithFields(logrus.Fields{
			"tool":         tool.Tool.Name,
			"description":  tool.Tool.Description,
			"input_schema": tool.Tool.InputSchema,
		}).Info("Registering tool")
		s.AddTool(tool.Tool, tool.Handler)
	}

	logrus.Info("Server setup complete. Listening for requests...")

	if err := server.ServeStdio(s); err != nil {
		logrus.WithError(err).Fatal("Server error")
	}
}
