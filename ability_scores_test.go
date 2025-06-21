package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestRunAbilityScoreTool(t *testing.T) {
	abilityScore := abilityScoreDetail{
		Index:    "str",
		Name:     "STR",
		FullName: "Strength",
		Desc:     []string{"Strength measures bodily power."},
		Skills: []struct {
			Index string `json:"index"`
			Name  string `json:"name"`
			URL   string `json:"url"`
		}{{Index: "athletics", Name: "Athletics", URL: "/api/2014/skills/athletics"}},
		URL: "/api/2014/ability-scores/str",
	}
	list := []abilityScoreListAPIResponse{{Index: "str", Name: "STR", URL: "/api/2014/ability-scores/str"}}

	cases := []struct {
		name       string
		input      abilityScoreToolInput
		mockByName func(*http.Client, endpoint, string, any) error
		mockList   func(*http.Client, endpoint, any, string) error
		wantOutput abilityScoreToolOutput
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:  "by name",
			input: abilityScoreToolInput{Name: "str"},
			mockByName: func(_ *http.Client, _ endpoint, name string, v any) error {
				ptr, ok := v.(*abilityScoreDetail)
				if !ok {
					return errors.New("wrong type")
				}
				*ptr = abilityScore
				return nil
			},
			mockList:   func(_ *http.Client, _ endpoint, v any, _ string) error { return nil },
			wantOutput: abilityScoreToolOutput{AbilityScore: &abilityScore},
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name:       "list",
			input:      abilityScoreToolInput{},
			mockByName: func(_ *http.Client, _ endpoint, name string, v any) error { return nil },
			mockList: func(_ *http.Client, _ endpoint, v any, _ string) error {
				ptr, ok := v.(*[]abilityScoreListAPIResponse)
				if !ok {
					return errors.New("wrong type")
				}
				*ptr = list
				return nil
			},
			wantOutput: abilityScoreToolOutput{Count: len(list), Results: list},
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name:  "fetchByName error",
			input: abilityScoreToolInput{Name: "fail"},
			mockByName: func(_ *http.Client, _ endpoint, name string, v any) error {
				return errors.New("fetchByName failed")
			},
			mockList:   func(_ *http.Client, _ endpoint, v any, _ string) error { return nil },
			wantOutput: abilityScoreToolOutput{},
			wantErr:    true,
			wantErrMsg: "fetchByName failed",
		},
		{
			name:       "fetchList error",
			input:      abilityScoreToolInput{},
			mockByName: func(_ *http.Client, _ endpoint, name string, v any) error { return nil },
			mockList: func(_ *http.Client, _ endpoint, v any, _ string) error {
				return errors.New("fetchList failed")
			},
			wantOutput: abilityScoreToolOutput{},
			wantErr:    true,
			wantErrMsg: "fetchList failed",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := runAbilityScoreTool(tc.input, tc.mockByName, tc.mockList)
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected Go error, got nil")
				}
				if res == nil || !res.IsError {
					t.Fatalf("expected MCP error result, got %+v", res)
				}
				if len(res.Content) > 0 {
					txt, ok := mcp.AsTextContent(res.Content[0])
					if !ok {
						t.Fatalf("content is not TextContent, got %T", res.Content[0])
					}
					if tc.wantErrMsg != "" && !strings.Contains(txt.Text, tc.wantErrMsg) {
						t.Errorf("expected error message to contain %q, got %q", tc.wantErrMsg, txt.Text)
					}
				}
				return
			}
			if res == nil {
				t.Fatalf("unexpected nil result")
			}
			var out abilityScoreToolOutput
			if len(res.Content) == 0 {
				t.Fatalf("no content in result")
			}
			txt, ok := mcp.AsTextContent(res.Content[0])
			if !ok {
				t.Fatalf("content is not TextContent, got %T", res.Content[0])
			}
			jsonStr := txt.Text

			if err := json.Unmarshal([]byte(jsonStr), &out); err != nil {
				t.Fatalf("unmarshal output: %v", err)
			}
			if tc.input.Name != "" {
				if out.AbilityScore == nil || out.AbilityScore.Name != tc.wantOutput.AbilityScore.Name {
					t.Errorf("expected ability score name %q, got %+v", tc.wantOutput.AbilityScore.Name, out.AbilityScore)
				}
			} else {
				if out.Count != tc.wantOutput.Count || len(out.Results) != len(tc.wantOutput.Results) {
					t.Errorf("expected %d ability scores, got count=%d, results=%d", tc.wantOutput.Count, out.Count, len(out.Results))
				}
				for i, a := range tc.wantOutput.Results {
					if out.Results[i].Name != a.Name {
						t.Errorf("expected ability score %d to be %q, got %q", i, a.Name, out.Results[i].Name)
					}
				}
			}
		})
	}
}
