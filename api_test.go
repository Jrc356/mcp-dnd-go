package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchCategories(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"spells": "/api/spells", "monsters": "/api/monsters"}`)
	}))
	defer ts.Close()
	client := ts.Client()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"basic categories fetch", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FetchCategories(client, ts.URL)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchCategories() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got == nil {
				t.Errorf("FetchCategories() got = nil, want non-nil")
			}
		})
	}
}

func TestFetchItems(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/spells" {
			fmt.Fprint(w, `{"results": [{"name": "Fireball", "index": "fireball"}]}`)
		} else {
			w.WriteHeader(404)
			io.WriteString(w, "not found")
		}
	}))
	defer ts.Close()
	client := ts.Client()
	tests := []struct {
		name     string
		category string
		wantErr  bool
	}{
		{"valid category", "spells", false},
		{"invalid category", "notacat", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FetchItems(client, ts.URL, tt.category)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchItems() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got == nil {
				t.Errorf("FetchItems() got = nil, want non-nil")
			}
		})
	}
}

func TestFetchItem(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/spells/fireball" {
			fmt.Fprint(w, `{"name": "Fireball", "index": "fireball"}`)
		} else {
			w.WriteHeader(404)
			io.WriteString(w, "not found")
		}
	}))
	defer ts.Close()
	client := ts.Client()
	tests := []struct {
		name     string
		category string
		index    string
		wantErr  bool
	}{
		{"valid item", "spells", "fireball", false},
		{"invalid item", "spells", "notaspell", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FetchItem(client, ts.URL, tt.category, tt.index)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchItem() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got == nil {
				t.Errorf("FetchItem() got = nil, want non-nil")
			}
		})
	}
}
