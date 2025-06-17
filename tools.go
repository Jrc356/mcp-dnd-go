package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerTools(s *server.MCPServer, httpClient *http.Client, apiBaseURL string) {
	registerListResourcesTool(s, httpClient, apiBaseURL, "list_resources", "List all available D&D 5e API resource categories (API root).")
	registerListCategoryItemsTool(s, httpClient, apiBaseURL, "list_category_items", "List all items/resources in a D&D 5e API category (e.g., spells, monsters, equipment).")
	registerGetResourceByIndexTool(s, httpClient, apiBaseURL, "get_resource_by_index", "Get a specific D&D 5e resource by category and index (e.g., a spell, monster, or item).")
	registerFilterSpellsTool(s, httpClient, apiBaseURL, "filter_spells", "List spells with optional filters (e.g., by class, level, school, etc.).")
	registerFilterMonstersTool(s, httpClient, apiBaseURL, "filter_monsters", "List monsters with optional filters (e.g., by type, challenge rating, etc.).")
	registerRandomMonsterTool(s, httpClient, apiBaseURL, "random_monster", "Get a random monster from the D&D 5e API.")
	registerRandomSpellTool(s, httpClient, apiBaseURL, "random_spell", "Get a random spell from the D&D 5e API.")
	registerListNamesInCategoryTool(s, httpClient, apiBaseURL, "list_names_in_category", "List all names in a D&D 5e API category (e.g., all monster names, spell names, etc.).")
	registerGetClassFeaturesTool(s, httpClient, apiBaseURL, "get_class_features", "Get all features for a D&D 5e class (by class index).")
	registerGetMonsterActionsTool(s, httpClient, apiBaseURL, "get_monster_actions", "Get all actions for a D&D 5e monster (by monster index).")
	registerSummarizeMonsterTool(s, httpClient, apiBaseURL, "summarize_monster", "Get a concise stat block summary for a D&D 5e monster (by monster index).")
	registerSummarizeSpellTool(s, httpClient, apiBaseURL, "summarize_spell", "Get a concise summary for a D&D 5e spell (by spell index).")
	registerFilterItemsTool(s, httpClient, apiBaseURL, "filter_items", "List items with optional filters (e.g., by equipment_category, cost, etc.).")
	registerEquipmentByTypeTool(s, httpClient, apiBaseURL, "equipment_by_type", "List all equipment of a specific type (e.g., weapons, armor, tools).")
	registerSpellSlotsTableTool(s, httpClient, apiBaseURL, "spell_slots_table", "Return the spell slots table for a given class and level.")
	registerMonsterByCRTool(s, httpClient, apiBaseURL, "monster_by_cr", "List or fetch a random monster by challenge rating (CR).")
	registerSpellsBySchoolTool(s, httpClient, apiBaseURL, "spells_by_school", "List all spells from a specific school of magic.")
	registerRacesAndTraitsTool(s, httpClient, apiBaseURL, "races_and_traits", "List all playable races and their traits.")
	registerBackgroundsTool(s, httpClient, apiBaseURL, "backgrounds", "List all backgrounds and their features.")
	registerFeatsTool(s, httpClient, apiBaseURL, "feats", "List all feats and their details.")
	registerConditionsTool(s, httpClient, apiBaseURL, "conditions", "List all conditions (e.g., blinded, stunned) and their effects.")
	registerDamageTypesTool(s, httpClient, apiBaseURL, "damage_types", "List all damage types and their descriptions.")
	registerMagicItemsTool(s, httpClient, apiBaseURL, "get_magic_items", "List all magic items, optionally filterable by rarity or type.")
}

func handleGetClassFeatures(req mcp.CallToolRequest, client *http.Client, apiBaseURL string) (interface{}, error) {
	classIndex, err := req.RequireString("class_index")
	if err != nil {
		return nil, err
	}
	features, err := FetchClassFeatures(client, apiBaseURL, classIndex)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"class_index": classIndex, "features": features}, nil
}

func registerGetClassFeaturesTool(s *server.MCPServer, client *http.Client, apiBaseURL string, name, description string) {
	s.AddTool(mcp.NewTool(name,
		mcp.WithDescription(description),
		mcp.WithString("class_index", mcp.Required(), mcp.Description("The class index, e.g., 'wizard', 'fighter'")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := handleGetClassFeatures(req, client, apiBaseURL)
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

func handleGetMonsterActions(req mcp.CallToolRequest, client *http.Client, apiBaseURL string) (interface{}, error) {
	monsterIndex, err := req.RequireString("monster_index")
	if err != nil {
		return nil, err
	}
	data, err := FetchItem(client, apiBaseURL, "monsters", monsterIndex)
	if err != nil {
		return nil, err
	}
	actions, ok := data["actions"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("actions not found or invalid for monster %s", monsterIndex)
	}
	return map[string]interface{}{"monster_index": monsterIndex, "actions": actions}, nil
}

func registerGetMonsterActionsTool(s *server.MCPServer, client *http.Client, apiBaseURL string, name, description string) {
	s.AddTool(mcp.NewTool(name,
		mcp.WithDescription(description),
		mcp.WithString("monster_index", mcp.Required(), mcp.Description("The monster index, e.g., 'adult-red-dragon'")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := handleGetMonsterActions(req, client, apiBaseURL)
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

func handleSummarizeMonster(req mcp.CallToolRequest, client *http.Client, apiBaseURL string) (interface{}, error) {
	monsterIndex, err := req.RequireString("monster_index")
	if err != nil {
		return nil, err
	}
	data, err := FetchItem(client, apiBaseURL, "monsters", monsterIndex)
	if err != nil {
		return nil, err
	}
	summary := map[string]interface{}{
		"name":             data["name"],
		"size":             data["size"],
		"type":             data["type"],
		"alignment":        data["alignment"],
		"armor_class":      data["armor_class"],
		"hit_points":       data["hit_points"],
		"hit_dice":         data["hit_dice"],
		"challenge_rating": data["challenge_rating"],
		"xp":               data["xp"],
		"speed":            data["speed"],
		"abilities": map[string]interface{}{
			"strength":     data["strength"],
			"dexterity":    data["dexterity"],
			"constitution": data["constitution"],
			"intelligence": data["intelligence"],
			"wisdom":       data["wisdom"],
			"charisma":     data["charisma"],
		},
		"proficiencies":     data["proficiencies"],
		"senses":            data["senses"],
		"languages":         data["languages"],
		"special_abilities": data["special_abilities"],
		"actions":           data["actions"],
		"legendary_actions": data["legendary_actions"],
	}
	return summary, nil
}

func registerSummarizeMonsterTool(s *server.MCPServer, client *http.Client, apiBaseURL string, name, description string) {
	s.AddTool(mcp.NewTool(name,
		mcp.WithDescription(description),
		mcp.WithString("monster_index", mcp.Required(), mcp.Description("The monster index, e.g., 'adult-red-dragon'")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := handleSummarizeMonster(req, client, apiBaseURL)
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

func handleSummarizeSpell(req mcp.CallToolRequest, client *http.Client, apiBaseURL string) (interface{}, error) {
	spellIndex, err := req.RequireString("spell_index")
	if err != nil {
		return nil, err
	}
	data, err := FetchItem(client, apiBaseURL, "spells", spellIndex)
	if err != nil {
		return nil, err
	}
	summary := map[string]interface{}{
		"name":          data["name"],
		"level":         data["level"],
		"school":        data["school"],
		"casting_time":  data["casting_time"],
		"range":         data["range"],
		"components":    data["components"],
		"duration":      data["duration"],
		"concentration": data["concentration"],
		"ritual":        data["ritual"],
		"desc":          data["desc"],
		"higher_level":  data["higher_level"],
	}
	return summary, nil
}

func registerSummarizeSpellTool(s *server.MCPServer, client *http.Client, apiBaseURL string, name, description string) {
	s.AddTool(mcp.NewTool(name,
		mcp.WithDescription(description),
		mcp.WithString("spell_index", mcp.Required(), mcp.Description("The spell index, e.g., 'fireball'")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := handleSummarizeSpell(req, client, apiBaseURL)
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

func handleFilterItems(req mcp.CallToolRequest, client *http.Client, apiBaseURL string) (interface{}, error) {
	filter := req.GetString("filter", "")
	url := apiBaseURL + "/equipment"
	if filter != "" {
		url += "?" + filter
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}
	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

func registerFilterItemsTool(s *server.MCPServer, client *http.Client, apiBaseURL string, name, description string) {
	s.AddTool(mcp.NewTool(name,
		mcp.WithDescription(description),
		mcp.WithString("filter", mcp.Description("Query string for filtering items, e.g., 'equipment_category=weapon&cost=1gp'")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := handleFilterItems(req, client, apiBaseURL)
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

func handleListResources(client *http.Client, apiBaseURL string) (interface{}, error) {
	result, err := FetchCategories(client, apiBaseURL)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func registerListResourcesTool(s *server.MCPServer, client *http.Client, apiBaseURL string, name, description string) {
	s.AddTool(mcp.NewTool(name,
		mcp.WithDescription(description),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := handleListResources(client, apiBaseURL)
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

func handleListCategoryItems(req mcp.CallToolRequest, client *http.Client, apiBaseURL string) (interface{}, error) {
	category, err := req.RequireString("category")
	if err != nil {
		return nil, err
	}
	result, err := FetchItems(client, apiBaseURL, category)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func registerListCategoryItemsTool(s *server.MCPServer, client *http.Client, apiBaseURL string, name, description string) {
	s.AddTool(mcp.NewTool(name,
		mcp.WithDescription(description),
		mcp.WithString("category", mcp.Required(), mcp.Description("The D&D API category to retrieve items from (e.g., 'spells', 'monsters', 'equipment')")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := handleListCategoryItems(req, client, apiBaseURL)
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

func handleGetResourceByIndex(req mcp.CallToolRequest, client *http.Client, apiBaseURL string) (interface{}, error) {
	category, err := req.RequireString("category")
	if err != nil {
		return nil, err
	}
	index, err := req.RequireString("index")
	if err != nil {
		return nil, err
	}
	result, err := FetchItem(client, apiBaseURL, category, index)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func registerGetResourceByIndexTool(s *server.MCPServer, client *http.Client, apiBaseURL string, name, description string) {
	s.AddTool(mcp.NewTool(name,
		mcp.WithDescription(description),
		mcp.WithString("category", mcp.Required(), mcp.Description("The D&D API category the item belongs to (e.g., 'spells', 'monsters', 'equipment')")),
		mcp.WithString("index", mcp.Required(), mcp.Description("The unique identifier for the specific item (e.g., 'fireball', 'adult-red-dragon')")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := handleGetResourceByIndex(req, client, apiBaseURL)
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

func handleFilterSpells(req mcp.CallToolRequest, client *http.Client, apiBaseURL string) (interface{}, error) {
	filter := req.GetString("filter", "")
	url := apiBaseURL + "/spells"
	if filter != "" {
		url += "?" + filter
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}
	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

func registerFilterSpellsTool(s *server.MCPServer, client *http.Client, apiBaseURL string, name, description string) {
	s.AddTool(mcp.NewTool(name,
		mcp.WithDescription(description),
		mcp.WithString("filter", mcp.Description("Query string for filtering spells, e.g., 'level=3&school=evocation'")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := handleFilterSpells(req, client, apiBaseURL)
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

func handleFilterMonsters(req mcp.CallToolRequest, client *http.Client, apiBaseURL string) (interface{}, error) {
	filter := req.GetString("filter", "")
	url := apiBaseURL + "/monsters"
	if filter != "" {
		url += "?" + filter
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}
	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

func registerFilterMonstersTool(s *server.MCPServer, client *http.Client, apiBaseURL string, name, description string) {
	s.AddTool(mcp.NewTool(name,
		mcp.WithDescription(description),
		mcp.WithString("filter", mcp.Description("Query string for filtering monsters, e.g., 'type=dragon&challenge_rating=10'")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := handleFilterMonsters(req, client, apiBaseURL)
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

func handleRandomMonster(client *http.Client, apiBaseURL string) (interface{}, error) {
	result, err := FetchItems(client, apiBaseURL, "monsters")
	if err != nil {
		return nil, err
	}
	items, ok := result["items"].([]Item)
	if !ok || len(items) == 0 {
		return nil, fmt.Errorf("no monsters found")
	}
	idx := int(time.Now().UnixNano()) % len(items)
	monster := items[idx]
	data, err := FetchItem(client, apiBaseURL, "monsters", monster.Index)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func registerRandomMonsterTool(s *server.MCPServer, client *http.Client, apiBaseURL string, name, description string) {
	s.AddTool(mcp.NewTool(name,
		mcp.WithDescription(description),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := handleRandomMonster(client, apiBaseURL)
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

func handleRandomSpell(client *http.Client, apiBaseURL string) (interface{}, error) {
	result, err := FetchItems(client, apiBaseURL, "spells")
	if err != nil {
		return nil, err
	}
	items, ok := result["items"].([]Item)
	if !ok || len(items) == 0 {
		return nil, fmt.Errorf("no spells found")
	}
	idx := int(time.Now().UnixNano()) % len(items)
	spell := items[idx]
	data, err := FetchItem(client, apiBaseURL, "spells", spell.Index)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func registerRandomSpellTool(s *server.MCPServer, client *http.Client, apiBaseURL string, name, description string) {
	s.AddTool(mcp.NewTool(name,
		mcp.WithDescription(description),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := handleRandomSpell(client, apiBaseURL)
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

func registerListNamesInCategoryTool(s *server.MCPServer, client *http.Client, apiBaseURL string, name, description string) {
	s.AddTool(mcp.NewTool(name,
		mcp.WithDescription(description),
		mcp.WithString("category", mcp.Required(), mcp.Description("The D&D API category to retrieve names from (e.g., 'spells', 'monsters', 'equipment')")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		category, err := req.RequireString("category")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		result, err := FetchItems(client, apiBaseURL, category)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		items, ok := result["items"].([]Item)
		if !ok {
			return mcp.NewToolResultError("No items found in category: " + category), nil
		}
		names := make([]string, 0, len(items))
		for _, item := range items {
			names = append(names, item.Name)
		}
		jsonBytes, err := json.Marshal(map[string]interface{}{"category": category, "names": names})
		if err != nil {
			return mcp.NewToolResultError("Failed to marshal result: " + err.Error()), nil
		}
		return mcp.NewToolResultText(string(jsonBytes)), nil
	})
}

func handleEquipmentByType(req mcp.CallToolRequest, client *http.Client, apiBaseURL string) (interface{}, error) {
	equipmentType, err := req.RequireString("equipment_type")
	if err != nil {
		return nil, err
	}
	result, err := FetchItems(client, apiBaseURL, "equipment")
	if err != nil {
		return nil, err
	}
	items, ok := result["items"].([]Item)
	if !ok {
		return nil, fmt.Errorf("no equipment found")
	}
	filtered := []Item{}
	for _, item := range items {
		itemData, err := FetchItem(client, apiBaseURL, "equipment", item.Index)
		if err == nil {
			if t, ok := itemData["equipment_category"].(map[string]interface{}); ok {
				if t["index"] == equipmentType {
					filtered = append(filtered, item)
				}
			}
		}
	}
	return map[string]interface{}{"equipment_type": equipmentType, "items": filtered}, nil
}

func registerEquipmentByTypeTool(s *server.MCPServer, client *http.Client, apiBaseURL string, name, description string) {
	s.AddTool(mcp.NewTool(name,
		mcp.WithDescription(description),
		mcp.WithString("equipment_type", mcp.Required(), mcp.Description("The equipment type, e.g., 'weapon', 'armor', 'tool'")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := handleEquipmentByType(req, client, apiBaseURL)
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

func handleSpellSlotsTable(req mcp.CallToolRequest, client *http.Client, apiBaseURL string) (interface{}, error) {
	classIndex, err := req.RequireString("class_index")
	if err != nil {
		return nil, err
	}
	levelStr, err := req.RequireString("level")
	if err != nil {
		return nil, err
	}
	classData, err := FetchItem(client, apiBaseURL, "classes", classIndex)
	if err != nil {
		return nil, err
	}
	levels, ok := classData["spellcasting"]
	if !ok {
		return nil, fmt.Errorf("spellcasting not found for class %s", classIndex)
	}
	return map[string]interface{}{"class_index": classIndex, "level": levelStr, "spell_slots": levels}, nil
}

func registerSpellSlotsTableTool(s *server.MCPServer, client *http.Client, apiBaseURL string, name, description string) {
	s.AddTool(mcp.NewTool(name,
		mcp.WithDescription(description),
		mcp.WithString("class_index", mcp.Required(), mcp.Description("The class index, e.g., 'wizard', 'cleric'")),
		mcp.WithString("level", mcp.Required(), mcp.Description("The class level (1-20)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := handleSpellSlotsTable(req, client, apiBaseURL)
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

func handleMonsterByCR(req mcp.CallToolRequest, client *http.Client, apiBaseURL string) (interface{}, error) {
	crStr, err := req.RequireString("cr")
	if err != nil {
		return nil, err
	}
	result, err := FetchItems(client, apiBaseURL, "monsters")
	if err != nil {
		return nil, err
	}
	items, ok := result["items"].([]Item)
	if !ok {
		return nil, fmt.Errorf("no monsters found")
	}
	filtered := []Item{}
	for _, item := range items {
		itemData, err := FetchItem(client, apiBaseURL, "monsters", item.Index)
		if err == nil {
			if fmt.Sprintf("%v", itemData["challenge_rating"]) == crStr {
				filtered = append(filtered, item)
			}
		}
	}
	return map[string]interface{}{"cr": crStr, "monsters": filtered}, nil
}

func registerMonsterByCRTool(s *server.MCPServer, client *http.Client, apiBaseURL string, name, description string) {
	s.AddTool(mcp.NewTool(name,
		mcp.WithDescription(description),
		mcp.WithString("cr", mcp.Required(), mcp.Description("The challenge rating (e.g., 1, 2, 5, 10)")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := handleMonsterByCR(req, client, apiBaseURL)
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

func handleSpellsBySchool(req mcp.CallToolRequest, client *http.Client, apiBaseURL string) (interface{}, error) {
	schoolIndex, err := req.RequireString("school_index")
	if err != nil {
		return nil, err
	}
	result, err := FetchItems(client, apiBaseURL, "spells")
	if err != nil {
		return nil, err
	}
	items, ok := result["items"].([]Item)
	if !ok {
		return nil, fmt.Errorf("no spells found")
	}
	filtered := []Item{}
	for _, item := range items {
		itemData, err := FetchItem(client, apiBaseURL, "spells", item.Index)
		if err == nil {
			if school, ok := itemData["school"].(map[string]interface{}); ok {
				if school["index"] == schoolIndex {
					filtered = append(filtered, item)
				}
			}
		}
	}
	return map[string]interface{}{"school_index": schoolIndex, "spells": filtered}, nil
}

func registerSpellsBySchoolTool(s *server.MCPServer, client *http.Client, apiBaseURL string, name, description string) {
	s.AddTool(mcp.NewTool(name,
		mcp.WithDescription(description),
		mcp.WithString("school_index", mcp.Required(), mcp.Description("The school index, e.g., 'evocation', 'illusion'")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := handleSpellsBySchool(req, client, apiBaseURL)
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

func handleRacesAndTraits(client *http.Client, apiBaseURL string) (interface{}, error) {
	result, err := FetchItems(client, apiBaseURL, "races")
	if err != nil {
		return nil, err
	}
	items, ok := result["items"].([]Item)
	if !ok {
		return nil, fmt.Errorf("no races found")
	}
	races := []map[string]interface{}{}
	for _, item := range items {
		itemData, err := FetchItem(client, apiBaseURL, "races", item.Index)
		if err == nil {
			race := map[string]interface{}{"name": itemData["name"], "traits": itemData["traits"]}
			races = append(races, race)
		}
	}
	return map[string]interface{}{"races": races}, nil
}

func registerRacesAndTraitsTool(s *server.MCPServer, client *http.Client, apiBaseURL string, name, description string) {
	s.AddTool(mcp.NewTool(name,
		mcp.WithDescription(description),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := handleRacesAndTraits(client, apiBaseURL)
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

func handleBackgrounds(client *http.Client, apiBaseURL string) (interface{}, error) {
	result, err := FetchItems(client, apiBaseURL, "backgrounds")
	if err != nil {
		return nil, err
	}
	items, ok := result["items"].([]Item)
	if !ok {
		return nil, fmt.Errorf("no backgrounds found")
	}
	backgrounds := []map[string]interface{}{}
	for _, item := range items {
		itemData, err := FetchItem(client, apiBaseURL, "backgrounds", item.Index)
		if err == nil {
			background := map[string]interface{}{"name": itemData["name"], "feature": itemData["feature"]}
			backgrounds = append(backgrounds, background)
		}
	}
	return map[string]interface{}{"backgrounds": backgrounds}, nil
}

func registerBackgroundsTool(s *server.MCPServer, client *http.Client, apiBaseURL string, name, description string) {
	s.AddTool(mcp.NewTool(name,
		mcp.WithDescription(description),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := handleBackgrounds(client, apiBaseURL)
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

func handleFeats(client *http.Client, apiBaseURL string) (interface{}, error) {
	result, err := FetchItems(client, apiBaseURL, "feats")
	if err != nil {
		return nil, err
	}
	items, ok := result["items"].([]Item)
	if !ok {
		return nil, fmt.Errorf("no feats found")
	}
	feats := []map[string]interface{}{}
	for _, item := range items {
		itemData, err := FetchItem(client, apiBaseURL, "feats", item.Index)
		if err == nil {
			feats = append(feats, itemData)
		}
	}
	return map[string]interface{}{"feats": feats}, nil
}

func registerFeatsTool(s *server.MCPServer, client *http.Client, apiBaseURL string, name, description string) {
	s.AddTool(mcp.NewTool(name,
		mcp.WithDescription(description),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := handleFeats(client, apiBaseURL)
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

func handleConditions(client *http.Client, apiBaseURL string) (interface{}, error) {
	result, err := FetchItems(client, apiBaseURL, "conditions")
	if err != nil {
		return nil, err
	}
	items, ok := result["items"].([]Item)
	if !ok {
		return nil, fmt.Errorf("no conditions found")
	}
	conditions := []map[string]interface{}{}
	for _, item := range items {
		itemData, err := FetchItem(client, apiBaseURL, "conditions", item.Index)
		if err == nil {
			conditions = append(conditions, itemData)
		}
	}
	return map[string]interface{}{"conditions": conditions}, nil
}

func registerConditionsTool(s *server.MCPServer, client *http.Client, apiBaseURL string, name, description string) {
	s.AddTool(mcp.NewTool(name,
		mcp.WithDescription(description),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := handleConditions(client, apiBaseURL)
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

func handleDamageTypes(client *http.Client, apiBaseURL string) (interface{}, error) {
	result, err := FetchItems(client, apiBaseURL, "damage-types")
	if err != nil {
		return nil, err
	}
	items, ok := result["items"].([]Item)
	if !ok {
		return nil, fmt.Errorf("no damage types found")
	}
	damageTypes := []map[string]interface{}{}
	for _, item := range items {
		itemData, err := FetchItem(client, apiBaseURL, "damage-types", item.Index)
		if err == nil {
			damageTypes = append(damageTypes, itemData)
		}
	}
	return map[string]interface{}{"damage_types": damageTypes}, nil
}

func registerDamageTypesTool(s *server.MCPServer, client *http.Client, apiBaseURL string, name, description string) {
	s.AddTool(mcp.NewTool(name,
		mcp.WithDescription(description),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := handleDamageTypes(client, apiBaseURL)
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

func handleMagicItems(req mcp.CallToolRequest, client *http.Client, apiBaseURL string) (interface{}, error) {
	rarity := req.GetString("rarity", "")
	typeFilter := req.GetString("type", "")
	result, err := FetchItems(client, apiBaseURL, "magic-items")
	if err != nil {
		return nil, err
	}
	items, ok := result["items"].([]Item)
	if !ok {
		return nil, fmt.Errorf("no magic items found")
	}
	filtered := []Item{}
	for _, item := range items {
		itemData, err := FetchItem(client, apiBaseURL, "magic-items", item.Index)
		if err == nil {
			match := true
			if rarity != "" {
				if r, ok := itemData["rarity"].(map[string]interface{}); ok {
					if r["name"] != rarity {
						match = false
					}
				}
			}
			if typeFilter != "" {
				if t, ok := itemData["equipment_category"].(map[string]interface{}); ok {
					if t["index"] != typeFilter {
						match = false
					}
				}
			}
			if match {
				filtered = append(filtered, item)
			}
		}
	}
	return map[string]interface{}{"rarity": rarity, "type": typeFilter, "items": filtered}, nil
}

func registerMagicItemsTool(s *server.MCPServer, client *http.Client, apiBaseURL string, name, description string) {
	s.AddTool(mcp.NewTool(name,
		mcp.WithDescription(description),
		mcp.WithString("rarity", mcp.Description("Optional rarity filter, e.g., 'rare', 'legendary'")),
		mcp.WithString("type", mcp.Description("Optional type filter, e.g., 'weapon', 'armor'")),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := handleMagicItems(req, client, apiBaseURL)
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
