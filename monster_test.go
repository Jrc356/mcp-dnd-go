package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestMonsterFilter_buildQueryString(t *testing.T) {
	cases := []struct {
		name   string
		filter *monsterFilter
		want   string
	}{
		{"nil filter", nil, ""},
		{"empty filter", &monsterFilter{}, ""},
		{"single CR", &monsterFilter{ChallengeRating: []float64{1}}, "challenge_rating=1"},
		{"multiple CRs", &monsterFilter{ChallengeRating: []float64{1, 2.5}}, "challenge_rating=1,2.5"},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got := ""
			if tt.filter != nil {
				got = tt.filter.buildQueryString()
			}
			if got != tt.want {
				t.Errorf("buildQueryString() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRunMonsterTool(t *testing.T) {
	monster := monsterDetail{Name: "Goblin", Index: "goblin"}
	monsters := []monsterListAPIResponse{{Name: "Goblin", Index: "goblin"}, {Name: "Orc", Index: "orc"}}
	cases := []struct {
		name       string
		input      monsterToolInput
		mockByName func(*http.Client, endpoint, string, any) error
		mockList   func(*http.Client, endpoint, any, filter) error
		wantOutput monsterToolOutput
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:  "by name",
			input: monsterToolInput{Name: "goblin"},
			mockByName: func(_ *http.Client, _ endpoint, name string, v any) error {
				ptr, ok := v.(*monsterDetail)
				if !ok {
					return errors.New("wrong type")
				}
				*ptr = monster
				return nil
			},
			mockList:   func(_ *http.Client, _ endpoint, v any, _ filter) error { return nil },
			wantOutput: monsterToolOutput{Monster: &monster},
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name:       "list",
			input:      monsterToolInput{},
			mockByName: func(_ *http.Client, _ endpoint, name string, v any) error { return nil },
			mockList: func(_ *http.Client, _ endpoint, v any, _ filter) error {
				ptr, ok := v.(*[]monsterListAPIResponse)
				if !ok {
					return errors.New("wrong type")
				}
				*ptr = monsters
				return nil
			},
			wantOutput: monsterToolOutput{Count: len(monsters), Results: monsters},
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name:  "fetchByName error",
			input: monsterToolInput{Name: "fail"},
			mockByName: func(_ *http.Client, _ endpoint, name string, v any) error {
				return errors.New("fetchByName failed")
			},
			mockList:   func(_ *http.Client, _ endpoint, v any, _ filter) error { return nil },
			wantOutput: monsterToolOutput{},
			wantErr:    true,
			wantErrMsg: "fetchByName failed",
		},
		{
			name:       "fetchList error",
			input:      monsterToolInput{},
			mockByName: func(_ *http.Client, _ endpoint, name string, v any) error { return nil },
			mockList: func(_ *http.Client, _ endpoint, v any, _ filter) error {
				return errors.New("fetchList failed")
			},
			wantOutput: monsterToolOutput{},
			wantErr:    true,
			wantErrMsg: "fetchList failed",
		},
		{
			name:       "empty list",
			input:      monsterToolInput{},
			mockByName: func(_ *http.Client, _ endpoint, name string, v any) error { return nil },
			mockList: func(_ *http.Client, _ endpoint, v any, _ filter) error {
				ptr, ok := v.(*[]monsterListAPIResponse)
				if !ok {
					return errors.New("wrong type")
				}
				*ptr = nil
				return nil
			},
			wantOutput: monsterToolOutput{Count: 0, Results: nil},
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name:  "nil monster from fetchByName",
			input: monsterToolInput{Name: "empty"},
			mockByName: func(_ *http.Client, _ endpoint, name string, v any) error {
				ptr, ok := v.(*monsterDetail)
				if !ok {
					return errors.New("wrong type")
				}
				*ptr = monsterDetail{} // zero value
				return nil
			},
			mockList:   func(_ *http.Client, _ endpoint, v any, _ filter) error { return nil },
			wantOutput: monsterToolOutput{Monster: &monsterDetail{}},
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name:       "tool error result",
			input:      monsterToolInput{Name: ""},
			mockByName: func(_ *http.Client, _ endpoint, name string, v any) error { return nil },
			mockList: func(_ *http.Client, _ endpoint, v any, _ filter) error {
				return errors.New("test Go error")
			},
			wantOutput: monsterToolOutput{},
			wantErr:    true,
			wantErrMsg: "test Go error",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := runMonsterTool(tc.input, tc.mockByName, tc.mockList)
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
			var out monsterToolOutput
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
				if out.Monster == nil || out.Monster.Name != tc.wantOutput.Monster.Name {
					t.Errorf("expected monster name %q, got %+v", tc.wantOutput.Monster.Name, out.Monster)
				}
			} else {
				if out.Count != tc.wantOutput.Count || len(out.Results) != len(tc.wantOutput.Results) {
					t.Errorf("expected %d monsters, got count=%d, results=%d", tc.wantOutput.Count, out.Count, len(out.Results))
				}
				for i, m := range tc.wantOutput.Results {
					if out.Results[i].Name != m.Name {
						t.Errorf("expected monster %d to be %q, got %q", i, m.Name, out.Results[i].Name)
					}
				}
			}
		})
	}
}
