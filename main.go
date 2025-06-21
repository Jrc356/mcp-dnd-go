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
	opts = append(opts, makeToolOptions(input)...)
	readonly := true
	opts = append(opts, mcp.WithToolAnnotation(mcp.ToolAnnotation{ReadOnlyHint: &readonly, OpenWorldHint: &readonly}))
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
		server.WithLogging(),
	)

	logrus.Info("Creating tools...")
	tools := []server.ServerTool{
		newAPITool(
			spells,
			"Fetches information about D&D 5e spells.",
			spellToolInput{},
			handleSpellTool,
		),
		newAPITool(
			monsters,
			"Fetches information about D&D 5e monsters.",
			monsterToolInput{},
			handleMonsterTool,
		),
		newAPITool(
			abilityScores,
			"Fetches information about D&D 5e ability scores.",
			abilityScoreToolInput{},
			handleAbilityScoreTool,
		),
		newAPITool(
			alignments,
			"Fetches information about D&D 5e alignments.",
			alignmentToolInput{},
			handleAlignmentTool,
		),
		newAPITool(
			backgrounds,
			"Fetches information about D&D 5e backgrounds.",
			backgroundToolInput{},
			handleBackgroundTool,
		),
		newAPITool(
			classes,
			"Fetches information about D&D 5e classes.",
			classToolInput{},
			handleClassTool,
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
