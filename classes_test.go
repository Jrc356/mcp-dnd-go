package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestRunClassTool(t *testing.T) {
	class := classDetail{
		Index:  "barbarian",
		Name:   "Barbarian",
		HitDie: 12,
		URL:    "/api/2014/classes/barbarian",
	}
	list := []classListAPIResponse{{Index: "barbarian", Name: "Barbarian", URL: "/api/2014/classes/barbarian"}}

	cases := []struct {
		name       string
		input      classToolInput
		mockByName func(*http.Client, endpoint, string, any) error
		mockList   func(*http.Client, endpoint, any, string) error
		wantOutput classToolOutput
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:  "by name",
			input: classToolInput{Name: "barbarian"},
			mockByName: func(_ *http.Client, _ endpoint, name string, v any) error {
				ptr, ok := v.(*classDetail)
				if !ok {
					return errors.New("wrong type")
				}
				*ptr = class
				return nil
			},
			mockList:   func(_ *http.Client, _ endpoint, v any, _ string) error { return nil },
			wantOutput: classToolOutput{Class: &class},
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name:       "list",
			input:      classToolInput{},
			mockByName: func(_ *http.Client, _ endpoint, name string, v any) error { return nil },
			mockList: func(_ *http.Client, _ endpoint, v any, _ string) error {
				ptr, ok := v.(*[]classListAPIResponse)
				if !ok {
					return errors.New("wrong type")
				}
				*ptr = list
				return nil
			},
			wantOutput: classToolOutput{Count: len(list), Results: list},
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name:  "fetchByName error",
			input: classToolInput{Name: "fail"},
			mockByName: func(_ *http.Client, _ endpoint, name string, v any) error {
				return errors.New("fetchByName failed")
			},
			mockList:   func(_ *http.Client, _ endpoint, v any, _ string) error { return nil },
			wantOutput: classToolOutput{},
			wantErr:    true,
			wantErrMsg: "fetchByName failed",
		},
		{
			name:       "fetchList error",
			input:      classToolInput{},
			mockByName: func(_ *http.Client, _ endpoint, name string, v any) error { return nil },
			mockList: func(_ *http.Client, _ endpoint, v any, _ string) error {
				return errors.New("fetchList failed")
			},
			wantOutput: classToolOutput{},
			wantErr:    true,
			wantErrMsg: "fetchList failed",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := runClassTool(tc.input, tc.mockByName, tc.mockList)
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
			var out classToolOutput
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
				if out.Class == nil || out.Class.Name != tc.wantOutput.Class.Name {
					t.Errorf("expected class name %q, got %+v", tc.wantOutput.Class.Name, out.Class)
				}
			} else {
				if out.Count != tc.wantOutput.Count || len(out.Results) != len(tc.wantOutput.Results) {
					t.Errorf("expected %d classes, got count=%d, results=%d", tc.wantOutput.Count, out.Count, len(out.Results))
				}
				for i, a := range tc.wantOutput.Results {
					if out.Results[i].Name != a.Name {
						t.Errorf("expected class %d to be %q, got %q", i, a.Name, out.Results[i].Name)
					}
				}
			}
		})
	}
}
