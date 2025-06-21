package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestSpellInput_buildQueryString(t *testing.T) {
	tests := []struct {
		name  string
		input *spellToolInput
		want  string
	}{
		{"nil filter", nil, ""},
		{"empty filter", &spellToolInput{}, ""},
		{"level only", &spellToolInput{Level: 3}, "level=3"},
		{"school only", &spellToolInput{School: "evocation"}, "school=evocation"},
		{"level and school", &spellToolInput{Level: 3, School: "evocation"}, "level=3&school=evocation"},
		{"all fields", &spellToolInput{Level: 3, School: "evocation"}, "level=3&school=evocation"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.buildQueryString()
			if got != tt.want {
				t.Errorf("buildQueryString() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRunSpellTool(t *testing.T) {
	// Load test data
	spellData, err := os.ReadFile("testdata/spell_by_name.json")
	if err != nil {
		t.Fatalf("failed to read spell_by_name.json: %v", err)
	}
	var spell spellAPIResponse
	if err := json.Unmarshal(spellData, &spell); err != nil {
		t.Fatalf("failed to unmarshal spell_by_name.json: %v", err)
	}

	listData, err := os.ReadFile("testdata/spell_list.json")
	if err != nil {
		t.Fatalf("failed to read spell_list.json: %v", err)
	}
	var spells []spellListAPIResponse
	if err := json.Unmarshal(listData, &spells); err != nil {
		t.Fatalf("failed to unmarshal spell_list.json: %v", err)
	}

	cases := []struct {
		name       string
		input      spellToolInput
		mockByName func(*http.Client, endpoint, string, any) error
		mockList   func(*http.Client, endpoint, any, string) error
		wantOutput spellToolOutput
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:  "by name",
			input: spellToolInput{Name: spell.Name},
			mockByName: func(_ *http.Client, _ endpoint, name string, v any) error {
				ptr, ok := v.(*spellAPIResponse)
				if !ok {
					return errors.New("wrong type")
				}
				*ptr = spell
				return nil
			},
			mockList:   func(_ *http.Client, _ endpoint, v any, _ string) error { return nil },
			wantOutput: spellToolOutput{Spell: &spell},
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name:       "list",
			input:      spellToolInput{},
			mockByName: func(_ *http.Client, _ endpoint, name string, v any) error { return nil },
			mockList: func(_ *http.Client, _ endpoint, v any, _ string) error {
				ptr, ok := v.(*[]spellListAPIResponse)
				if !ok {
					return errors.New("wrong type")
				}
				*ptr = spells
				return nil
			},
			wantOutput: spellToolOutput{Count: len(spells), Results: spells},
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name:  "fetchByName error",
			input: spellToolInput{Name: "fail"},
			mockByName: func(_ *http.Client, _ endpoint, name string, v any) error {
				return errors.New("fetchByName failed")
			},
			mockList:   func(_ *http.Client, _ endpoint, v any, _ string) error { return nil },
			wantOutput: spellToolOutput{},
			wantErr:    true,
			wantErrMsg: "fetchByName failed",
		},
		{
			name:       "fetchList error",
			input:      spellToolInput{},
			mockByName: func(_ *http.Client, _ endpoint, name string, v any) error { return nil },
			mockList: func(_ *http.Client, _ endpoint, v any, _ string) error {
				return errors.New("fetchList failed")
			},
			wantOutput: spellToolOutput{},
			wantErr:    true,
			wantErrMsg: "fetchList failed",
		},
		{
			name:       "empty list",
			input:      spellToolInput{},
			mockByName: func(_ *http.Client, _ endpoint, name string, v any) error { return nil },
			mockList: func(_ *http.Client, _ endpoint, v any, _ string) error {
				ptr, ok := v.(*[]spellListAPIResponse)
				if !ok {
					return errors.New("wrong type")
				}
				*ptr = nil
				return nil
			},
			wantOutput: spellToolOutput{Count: 0, Results: nil},
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name:  "nil spell from fetchByName",
			input: spellToolInput{Name: "empty"},
			mockByName: func(_ *http.Client, _ endpoint, name string, v any) error {
				ptr, ok := v.(*spellAPIResponse)
				if !ok {
					return errors.New("wrong type")
				}
				*ptr = spellAPIResponse{} // zero value
				return nil
			},
			mockList:   func(_ *http.Client, _ endpoint, v any, _ string) error { return nil },
			wantOutput: spellToolOutput{Spell: &spellAPIResponse{}},
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name:       "tool error result",
			input:      spellToolInput{Name: ""},
			mockByName: func(_ *http.Client, _ endpoint, name string, v any) error { return nil },
			mockList: func(_ *http.Client, _ endpoint, v any, _ string) error {
				return errors.New("test Go error")
			},
			wantOutput: spellToolOutput{},
			wantErr:    true,
			wantErrMsg: "test Go error",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := runSpellTool(tc.input, tc.mockByName, tc.mockList)
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
			var out spellToolOutput
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
				if out.Spell == nil || out.Spell.Name != tc.wantOutput.Spell.Name {
					t.Errorf("expected spell name %q, got %+v", tc.wantOutput.Spell.Name, out.Spell)
				}
			} else {
				if out.Count != tc.wantOutput.Count || len(out.Results) != len(tc.wantOutput.Results) {
					t.Errorf("expected %d spells, got count=%d, results=%d", tc.wantOutput.Count, out.Count, len(out.Results))
				}
				for i, s := range tc.wantOutput.Results {
					if out.Results[i].Name != s.Name {
						t.Errorf("expected spell %d to be %q, got %q", i, s.Name, out.Results[i].Name)
					}
				}
			}
		})
	}
}

func TestFetchSpellByNameResult(t *testing.T) {
	spell := spellAPIResponse{Name: "Magic Missile"}
	client := &http.Client{}
	cases := []struct {
		name       string
		input      spellToolInput
		mockFn     func(*http.Client, endpoint, string, any) error
		wantErr    bool
		wantErrMsg string
		wantName   string
	}{
		{
			name:  "success",
			input: spellToolInput{Name: "Magic Missile"},
			mockFn: func(_ *http.Client, _ endpoint, name string, v any) error {
				ptr := v.(*spellAPIResponse)
				*ptr = spell
				return nil
			},
			wantErr:  false,
			wantName: "Magic Missile",
		},
		{
			name:  "fetch error",
			input: spellToolInput{Name: "fail"},
			mockFn: func(_ *http.Client, _ endpoint, name string, v any) error {
				return errors.New("fail fetch")
			},
			wantErr:    true,
			wantErrMsg: "fail fetch",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := fetchSpellByNameResult(client, tc.input, tc.mockFn)
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				if res == nil || !res.IsError {
					t.Fatalf("expected MCP error result, got %+v", res)
				}
				if len(res.Content) > 0 && tc.wantErrMsg != "" {
					txt, _ := mcp.AsTextContent(res.Content[0])
					if !strings.Contains(txt.Text, tc.wantErrMsg) {
						t.Errorf("expected error message to contain %q, got %q", tc.wantErrMsg, txt.Text)
					}
				}
				return
			}
			if res == nil {
				t.Fatalf("unexpected nil result")
			}
			if len(res.Content) == 0 {
				t.Fatalf("no content in result")
			}
			txt, _ := mcp.AsTextContent(res.Content[0])
			var out spellToolOutput
			if err := json.Unmarshal([]byte(txt.Text), &out); err != nil {
				t.Fatalf("unmarshal output: %v", err)
			}
			if out.Spell == nil || out.Spell.Name != tc.wantName {
				t.Errorf("expected spell name %q, got %+v", tc.wantName, out.Spell)
			}
		})
	}
}

func TestFetchSpellListResult(t *testing.T) {
	client := &http.Client{}
	spells := []spellListAPIResponse{{Name: "A"}, {Name: "B"}}
	cases := []struct {
		name       string
		input      spellToolInput
		mockFn     func(*http.Client, endpoint, any, string) error
		wantErr    bool
		wantErrMsg string
		wantCount  int
	}{
		{
			name:  "success",
			input: spellToolInput{},
			mockFn: func(_ *http.Client, _ endpoint, v any, _ string) error {
				ptr := v.(*[]spellListAPIResponse)
				*ptr = spells
				return nil
			},
			wantErr:   false,
			wantCount: 2,
		},
		{
			name:  "fetch error",
			input: spellToolInput{},
			mockFn: func(_ *http.Client, _ endpoint, v any, _ string) error {
				return errors.New("fail list")
			},
			wantErr:    true,
			wantErrMsg: "fail list",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := fetchSpellListResult(client, tc.input, tc.mockFn)
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				if res == nil || !res.IsError {
					t.Fatalf("expected MCP error result, got %+v", res)
				}
				if len(res.Content) > 0 && tc.wantErrMsg != "" {
					txt, _ := mcp.AsTextContent(res.Content[0])
					if !strings.Contains(txt.Text, tc.wantErrMsg) {
						t.Errorf("expected error message to contain %q, got %q", tc.wantErrMsg, txt.Text)
					}
				}
				return
			}
			if res == nil {
				t.Fatalf("unexpected nil result")
			}
			if len(res.Content) == 0 {
				t.Fatalf("no content in result")
			}
			txt, _ := mcp.AsTextContent(res.Content[0])
			var out spellToolOutput
			if err := json.Unmarshal([]byte(txt.Text), &out); err != nil {
				t.Fatalf("unmarshal output: %v", err)
			}
			if out.Count != tc.wantCount {
				t.Errorf("expected count %d, got %d", tc.wantCount, out.Count)
			}
		})
	}
}
