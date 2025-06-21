package main

import (
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Name   string `json:"name" mcp:"description=The name of the item"`
	Level  int    `json:"level" mcp:"description=The level of the item"`
	Active bool   `json:"active" mcp:"description=Whether the item is active"`
}

type nestedStruct struct {
	Inner testStruct `json:"inner" mcp:"description=Inner struct"`
}

type badStructNoJSON struct {
	Name string `mcp:"description=Name"`
}

type badStructJSONIgnore struct {
	Name string `json:"-" mcp:"description=Should be ignored (json -)"`
}

type badStructMalformedMCP struct {
	Name string `json:"name" mcp:"desc"`
}

type arrayStruct struct {
	Tags   []string `json:"tags" mcp:"description=List of tags"`
	Scores []int    `json:"scores" mcp:"description=List of scores"`
}

func TestMakeToolOptions(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected mcp.Tool
	}{
		{
			"happy path - simple struct",
			testStruct{},
			mcp.NewTool("test",
				mcp.WithString("name", mcp.Description("The name of the item")),
				mcp.WithNumber("level", mcp.Description("The level of the item")),
				mcp.WithBoolean("active", mcp.Description("Whether the item is active")),
			),
		},
		{
			"happy path - nested struct",
			nestedStruct{},
			mcp.NewTool("test",
				mcp.WithObject("inner", mcp.Description("Inner struct"), mcp.Properties(makeProperties(testStruct{"Sword", 5, true}))),
			),
		},
		{
			"array fields",
			arrayStruct{},
			mcp.NewTool("test",
				mcp.WithArray("tags", mcp.Description("List of tags"), mcp.Items(map[string]any{"type": "string"})),
				mcp.WithArray("scores", mcp.Description("List of scores"), mcp.Items(map[string]any{"type": "number"})),
			),
		},
		{
			"field with no json tag",
			badStructNoJSON{},
			mcp.NewTool("test"),
		},
		{
			"field with json tag '-' (ignored)",
			badStructJSONIgnore{},
			mcp.NewTool("test"),
		},
		{
			"malformed mcp tag",
			badStructMalformedMCP{},
			mcp.NewTool("test"), // Should not panic, but will not add the field
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := makeToolOptions(tt.input)
			gotTool := mcp.NewTool("test", got...)
			expectedTool := tt.expected
			assert.Equal(t, expectedTool.InputSchema.Properties, gotTool.InputSchema.Properties)
		})
	}
}

func TestToolProperties(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected map[string]any
	}{
		{
			"simple struct",
			testStruct{"Sword", 5, true},
			map[string]any{
				"name":   map[string]any{"type": "string", "description": "The name of the item"},
				"level":  map[string]any{"type": "number", "description": "The level of the item"},
				"active": map[string]any{"type": "boolean", "description": "Whether the item is active"},
			},
		},
		{
			"nested struct",
			nestedStruct{Inner: testStruct{"Sword", 5, true}},
			map[string]any{
				"inner": map[string]any{
					"type":        "object",
					"description": "Inner struct",
					"properties": map[string]any{
						"name":   map[string]any{"type": "string", "description": "The name of the item"},
						"level":  map[string]any{"type": "number", "description": "The level of the item"},
						"active": map[string]any{"type": "boolean", "description": "Whether the item is active"},
					},
				},
			},
		},
		{
			"array fields",
			arrayStruct{Tags: []string{"a", "b"}, Scores: []int{1, 2}},
			map[string]any{
				"tags": map[string]any{
					"type":        "array",
					"description": "List of tags",
					"items":       map[string]any{"type": "string"},
				},
				"scores": map[string]any{
					"type":        "array",
					"description": "List of scores",
					"items":       map[string]any{"type": "number"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			props := makeProperties(tt.input)
			assert.Equal(t, tt.expected, props)
		})
	}
}

func TestToKebabCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello World", "hello-world"},
		{"  Leading and Trailing  ", "leading-and-trailing"},
		{"Already-Kebab", "already-kebab"},
		{"MiXeD CaSe", "mixed-case"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := toKebabCase(tt.input); got != tt.expected {
				t.Errorf("toKebabCase(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
