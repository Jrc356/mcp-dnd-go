package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

// monsterListAPIResponse defines the structure for a single monster in the list response.
type monsterListAPIResponse struct {
	Index string `json:"index"`
	Name  string `json:"name"`
	URL   string `json:"url"`
}

// monsterToolInput defines the input structure for the monster tool.
type monsterToolInput struct {
	Name            string    `json:"name" mcp:"description=The index of the monster to retrieve."`
	ChallengeRating []float64 `json:"challenge_rating" mcp:"description=The challenge rating(s) to filter on."`
}

// buildQueryString constructs a query string from the monsterFilter fields for use in API requests.
func (f *monsterToolInput) buildQueryString() string {
	if f == nil || len(f.ChallengeRating) == 0 {
		return ""
	}
	crs := make([]string, len(f.ChallengeRating))
	for i, cr := range f.ChallengeRating {
		crs[i] = strings.TrimRight(strings.TrimRight(strconv.FormatFloat(cr, 'f', 2, 64), "0"), ".")
	}
	return "challenge_rating=" + strings.Join(crs, ",")
}

// monsterToolOutput defines the output structure for the monster tool.
type monsterToolOutput struct {
	Count   int                      `json:"count,omitempty"`
	Results []monsterListAPIResponse `json:"results,omitempty"`
	Monster *monsterDetail           `json:"monster,omitempty"`
}

// monsterDetail defines the structure for a detailed monster response.
type monsterDetail struct {
	Index      string `json:"index"`
	Name       string `json:"name"`
	Size       string `json:"size"`
	Type       string `json:"type"`
	Alignment  string `json:"alignment"`
	ArmorClass []struct {
		Type  string `json:"type"`
		Value int    `json:"value"`
	} `json:"armor_class"`
	HitPoints     int    `json:"hit_points"`
	HitDice       string `json:"hit_dice"`
	HitPointsRoll string `json:"hit_points_roll"`
	Speed         struct {
		Walk string `json:"walk"`
		Swim string `json:"swim"`
	} `json:"speed"`
	Strength      int `json:"strength"`
	Dexterity     int `json:"dexterity"`
	Constitution  int `json:"constitution"`
	Intelligence  int `json:"intelligence"`
	Wisdom        int `json:"wisdom"`
	Charisma      int `json:"charisma"`
	Proficiencies []struct {
		Value       int `json:"value"`
		Proficiency struct {
			Index string `json:"index"`
			Name  string `json:"name"`
			URL   string `json:"url"`
		} `json:"proficiency"`
	} `json:"proficiencies"`
	DamageVulnerabilities []string `json:"damage_vulnerabilities"`
	DamageResistances     []string `json:"damage_resistances"`
	DamageImmunities      []string `json:"damage_immunities"`
	ConditionImmunities   []string `json:"condition_immunities"`
	Senses                struct {
		Darkvision        string `json:"darkvision"`
		PassivePerception int    `json:"passive_perception"`
	} `json:"senses"`
	Languages        string  `json:"languages"`
	ChallengeRating  float64 `json:"challenge_rating"`
	ProficiencyBonus int     `json:"proficiency_bonus"`
	XP               int     `json:"xp"`
	SpecialAbilities []struct {
		Name string `json:"name"`
		Desc string `json:"desc"`
		DC   struct {
			DCType struct {
				Index string `json:"index"`
				Name  string `json:"name"`
				URL   string `json:"url"`
			} `json:"dc_type"`
			DCValue     int    `json:"dc_value"`
			SuccessType string `json:"success_type"`
		} `json:"dc"`
		Damage []struct {
			DamageType struct {
				Index string `json:"index"`
				Name  string `json:"name"`
				URL   string `json:"url"`
			} `json:"damage_type"`
			DamageDice string `json:"damage_dice"`
		} `json:"damage"`
	} `json:"special_abilities"`
	Actions []struct {
		Name        string `json:"name"`
		Desc        string `json:"desc"`
		AttackBonus int    `json:"attack_bonus"`
		DC          struct {
			DCType struct {
				Index string `json:"index"`
				Name  string `json:"name"`
				URL   string `json:"url"`
			} `json:"dc_type"`
			DCValue     int    `json:"dc_value"`
			SuccessType string `json:"success_type"`
		} `json:"dc"`
		Damage []struct {
			DamageType struct {
				Index string `json:"index"`
				Name  string `json:"name"`
				URL   string `json:"url"`
			} `json:"damage_type"`
			DamageDice string `json:"damage_dice"`
		} `json:"damage"`
		Actions []struct {
			ActionName string `json:"action_name"`
			Count      string `json:"count"`
			Type       string `json:"type"`
		} `json:"actions"`
	} `json:"actions"`
	LegendaryActions []struct {
		Name   string `json:"name"`
		Desc   string `json:"desc"`
		Damage []struct {
			DamageType struct {
				Index string `json:"index"`
				Name  string `json:"name"`
				URL   string `json:"url"`
			} `json:"damage_type"`
			DamageDice string `json:"damage_dice"`
		} `json:"damage"`
	} `json:"legendary_actions"`
	Image     string `json:"image"`
	URL       string `json:"url"`
	UpdatedAt string `json:"updated_at"`
}

// fetchMonsterByNameResult fetches a monster by index and returns an MCP tool result.
func fetchMonsterByNameResult(
	client *http.Client,
	input monsterToolInput,
	fetchByName func(*http.Client, endpoint, string, any) error,
) (*mcp.CallToolResult, error) {
	monster := &monsterDetail{}
	err := fetchByName(client, monsters, input.Name, monster)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("failed to fetch monster", err), err
	}
	output := monsterToolOutput{Monster: monster}
	jsonData, err := json.Marshal(output)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to marshal monster output", err), err
	}
	return mcp.NewToolResultText(string(jsonData)), nil
}

// fetchMonsterListResult fetches a list of monsters with optional filtering and returns an MCP tool result.
func fetchMonsterListResult(
	client *http.Client,
	input monsterToolInput,
	fetchList func(*http.Client, endpoint, any, string) error,
) (*mcp.CallToolResult, error) {
	var results []monsterListAPIResponse
	err := fetchList(client, monsters, &results, input.buildQueryString())
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to fetch monster list", err), err
	}
	output := monsterToolOutput{Count: len(results), Results: results}
	jsonData, err := json.Marshal(output)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to marshal monster list output", err), err
	}
	return mcp.NewToolResultText(string(jsonData)), nil
}

// runMonsterTool executes the core logic for the monster tool, using injected fetchByName and fetchList dependencies for testability.
func runMonsterTool(
	input monsterToolInput,
	fetchByName func(*http.Client, endpoint, string, any) error,
	fetchList func(*http.Client, endpoint, any, string) error,
) (*mcp.CallToolResult, error) {
	client := http.DefaultClient
	if input.Name != "" {
		return fetchMonsterByNameResult(client, input, fetchByName)
	}
	return fetchMonsterListResult(client, input, fetchList)
}

// handleMonsterTool is the MCP handler for the monster tool.
func handleMonsterTool(ctx context.Context, req mcp.CallToolRequest, input monsterToolInput) (*mcp.CallToolResult, error) {
	return runMonsterTool(input, fetchByName, fetchList)
}
