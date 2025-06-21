package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestRunAlignmentTool(t *testing.T) {
	alignment := alignmentDetail{
		Index: "chaotic-good",
		Name:  "Chaotic Good",
		Desc:  "Acts as conscience directs, with little regard for what others expect.",
		URL:   "/api/2014/alignments/chaotic-good",
	}
	list := []alignmentListAPIResponse{{Index: "chaotic-good", Name: "Chaotic Good", URL: "/api/2014/alignments/chaotic-good"}}

	cases := []struct {
		name       string
		input      alignmentToolInput
		mockByName func(*http.Client, endpoint, string, any) error
		mockList   func(*http.Client, endpoint, any, string) error
		wantOutput alignmentToolOutput
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:  "by name",
			input: alignmentToolInput{Name: "chaotic-good"},
			mockByName: func(_ *http.Client, _ endpoint, name string, v any) error {
				ptr, ok := v.(*alignmentDetail)
				if !ok {
					return errors.New("wrong type")
				}
				*ptr = alignment
				return nil
			},
			mockList:   func(_ *http.Client, _ endpoint, v any, _ string) error { return nil },
			wantOutput: alignmentToolOutput{Alignment: &alignment},
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name:       "list",
			input:      alignmentToolInput{},
			mockByName: func(_ *http.Client, _ endpoint, name string, v any) error { return nil },
			mockList: func(_ *http.Client, _ endpoint, v any, _ string) error {
				ptr, ok := v.(*[]alignmentListAPIResponse)
				if !ok {
					return errors.New("wrong type")
				}
				*ptr = list
				return nil
			},
			wantOutput: alignmentToolOutput{Count: len(list), Results: list},
			wantErr:    false,
			wantErrMsg: "",
		},
		{
			name:  "fetchByName error",
			input: alignmentToolInput{Name: "fail"},
			mockByName: func(_ *http.Client, _ endpoint, name string, v any) error {
				return errors.New("fetchByName failed")
			},
			mockList:   func(_ *http.Client, _ endpoint, v any, _ string) error { return nil },
			wantOutput: alignmentToolOutput{},
			wantErr:    true,
			wantErrMsg: "fetchByName failed",
		},
		{
			name:       "fetchList error",
			input:      alignmentToolInput{},
			mockByName: func(_ *http.Client, _ endpoint, name string, v any) error { return nil },
			mockList: func(_ *http.Client, _ endpoint, v any, _ string) error {
				return errors.New("fetchList failed")
			},
			wantOutput: alignmentToolOutput{},
			wantErr:    true,
			wantErrMsg: "fetchList failed",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := runAlignmentTool(tc.input, tc.mockByName, tc.mockList)
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
			var out alignmentToolOutput
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
				if out.Alignment == nil || out.Alignment.Name != tc.wantOutput.Alignment.Name {
					t.Errorf("expected alignment name %q, got %+v", tc.wantOutput.Alignment.Name, out.Alignment)
				}
			} else {
				if out.Count != tc.wantOutput.Count || len(out.Results) != len(tc.wantOutput.Results) {
					t.Errorf("expected %d alignments, got count=%d, results=%d", tc.wantOutput.Count, out.Count, len(out.Results))
				}
				for i, a := range tc.wantOutput.Results {
					if out.Results[i].Name != a.Name {
						t.Errorf("expected alignment %d to be %q, got %q", i, a.Name, out.Results[i].Name)
					}
				}
			}
		})
	}
}
