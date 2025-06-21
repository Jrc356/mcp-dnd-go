package main

import (
	"reflect"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"
)

// structToToolOptions converts a struct to a slice of MCP ToolOptions.
// It inspects the struct fields, extracts JSON and MCP tags, and creates ToolOptions accordingly.
// Fields without a JSON tag or with a JSON tag of "-" are ignored.
// The MCP tag is expected to contain a description in the format "description=...".
// If the MCP tag is not present or does not contain a description, the field is skipped.
// If a field type is unsupported, it logs a warning and skips that field.
// Example usage:
//
//	type MyStruct struct {
//	    Name  string `json:"name" mcp:"description=The name of the item."`
//	    Count int    `json:"count" mcp:"description=The number of items."`
//	    Active bool   `json:"active" mcp:"description=Whether the item is active."`
//	    Tags  []string `json:"tags" mcp:"description=Tags associated with the item."`
//	}
func structToToolOptions(s any) []mcp.ToolOption {
	logrus.Debugf("Converting struct %T to MCP ToolOptions", s)
	var opts []mcp.ToolOption
	t := reflect.TypeOf(s)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		logrus.Debugf("Processing field: %s, type: %s", field.Name, field.Type)
		jsonTag := field.Tag.Get("json")
		mcpTag := field.Tag.Get("mcp")
		if jsonTag == "" {
			continue
		}
		name := strings.Split(jsonTag, ",")[0]
		if name == "-" {
			continue
		}

		mcpTagParts := strings.Split(mcpTag, ",")
		if len(mcpTagParts) == 0 || !strings.Contains(mcpTagParts[0], "=") {
			continue
		}
		description := strings.SplitN(mcpTagParts[0], "=", 2)[1]
		opt := fieldToToolOption(name, description, field.Type)
		opts = append(opts, opt)
	}
	return opts
}

// fieldToToolOption converts a struct field to an MCP ToolOption.
// It determines the field type and creates the appropriate ToolOption based on its kind.
// Supported types include string, number (int/float), boolean, struct, and slice.
func fieldToToolOption(name string, description string, fieldType reflect.Type) mcp.ToolOption {
	logrus.Debugf("Converting field '%s' of type '%s' to MCP ToolOption", name, fieldType)
	for fieldType.Kind() == reflect.Ptr {
		fieldType = fieldType.Elem()
	}
	switch fieldType.Kind() {
	case reflect.String:
		logrus.Debugf("field %s is a string", name)
		return mcp.WithString(name, mcp.Description(description))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64:
		logrus.Debugf("field %s is a number", name)
		return mcp.WithNumber(name, mcp.Description(description))
	case reflect.Bool:
		logrus.Debugf("field %s is a boolean", name)
		return mcp.WithBoolean(name, mcp.Description(description))
	case reflect.Struct:
		logrus.Debugf("field %s is a struct", name)
		return mcp.WithObject(name, mcp.Description(description), mcp.Properties(structToProperties(reflect.Zero(fieldType).Interface())))
	case reflect.Slice:
		logrus.Debugf("field %s is a slice", name)
		elemType := fieldType.Elem()
		for elemType.Kind() == reflect.Ptr {
			elemType = elemType.Elem()
		}
		var itemType string
		switch elemType.Kind() {
		case reflect.String:
			itemType = "string"
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64:
			itemType = "number"
		case reflect.Bool:
			itemType = "boolean"
		default:
			logrus.Warnf("%s is an unsupported slice element type: %s", name, elemType.Kind())
			return nil
		}
		return mcp.WithArray(name, mcp.Description(description), mcp.Items(map[string]any{"type": itemType}))
	default:
		logrus.Warnf("field %s is an unsupported type: %s", name, fieldType.Kind())
		return nil
	}
}

// fieldToProperty converts a struct field to a property map for use in MCP schemas.
// It extracts the field type and description from the MCP tag, and constructs a property map accordingly.
// If the field type is unsupported, it logs a warning and returns nil.
func fieldToProperty(description string, field reflect.StructField) map[string]any {
	logrus.Debugf("Converting field '%s' to property", field.Name)
	fieldType := field.Type
	for fieldType.Kind() == reflect.Ptr {
		fieldType = fieldType.Elem()
	}
	prop := map[string]any{
		"description": description,
	}
	switch fieldType.Kind() {
	case reflect.String:
		logrus.Debugf("field %s is a string", field.Name)
		prop["type"] = "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64:
		logrus.Debugf("field %s is a number", field.Name)
		prop["type"] = "number"
	case reflect.Bool:
		logrus.Debugf("field %s is a boolean", field.Name)
		prop["type"] = "boolean"
	case reflect.Struct:
		logrus.Debugf("field %s is a struct", field.Name)
		prop["type"] = "object"
		prop["properties"] = structToProperties(reflect.Zero(fieldType).Interface())
	case reflect.Slice:
		logrus.Debugf("field %s is a slice", field.Name)
		elemType := fieldType.Elem()
		for elemType.Kind() == reflect.Ptr {
			elemType = elemType.Elem()
		}
		var itemType string
		switch elemType.Kind() {
		case reflect.String:
			itemType = "string"
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64:
			itemType = "number"
		case reflect.Bool:
			itemType = "boolean"
		default:
			logrus.Warnf("%s is an unsupported slice element type: %s", field.Name, elemType.Kind())
			return nil
		}
		prop["type"] = "array"
		prop["items"] = map[string]any{"type": itemType}
	default:
		logrus.Warnf("%s is an unsupported type: %s", field.Name, fieldType.Kind())
		return nil
	}
	return prop
}

// structToProperties converts a struct to a map of properties for use in MCP schemas.
// It inspects the struct fields, extracts JSON and MCP tags, and creates a properties map.
// Fields without a JSON tag or with a JSON tag of "-" are ignored.
// The MCP tag is expected to contain a description in the format "description=...".
// If the MCP tag is not present or does not contain a description, the field is skipped.
// If a field type is unsupported, it logs a warning and skips that field.
func structToProperties(s any) map[string]any {
	logrus.Debugf("Converting struct %T to properties", s)
	props := make(map[string]any)
	t := reflect.TypeOf(s)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		mcpTag := field.Tag.Get("mcp")
		if jsonTag == "" {
			continue
		}
		name := strings.Split(jsonTag, ",")[0]
		if name == "-" {
			continue
		}

		mcpTagParts := strings.Split(mcpTag, ",")
		description := strings.Split(mcpTagParts[0], "=")[1]
		prop := fieldToProperty(description, field)
		props[name] = prop
	}
	return props
}

// toKebabCase converts a string to kebab-case.
// It converts the string to lowercase, trims whitespace, and replaces spaces with hyphens.
func toKebabCase(s string) string {
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, " ", "-")
	return s
}
