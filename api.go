package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Category struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URI         string `json:"uri"`
}

type Item struct {
	Name  string `json:"name"`
	Index string `json:"index"`
	URI   string `json:"uri"`
}

func FetchCategories(client *http.Client, apiBaseURL string) (map[string]interface{}, error) {
	resp, err := client.Get(apiBaseURL + "/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	categories := []Category{}
	for key := range data {
		desc, ok := CategoryDescriptions[key]
		if !ok {
			desc = fmt.Sprintf("Collection of D&D 5e %s", key)
		}
		categories = append(categories, Category{
			Name:        key,
			Description: desc,
			URI:         fmt.Sprintf("resource://dnd/items/%s", key),
		})
	}
	return map[string]interface{}{
		"categories": categories,
		"count":      len(categories),
	}, nil
}

func FetchItems(client *http.Client, apiBaseURL, category string) (map[string]interface{}, error) {
	resp, err := client.Get(fmt.Sprintf("%s/%s", apiBaseURL, category))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Category '%s' not found or API request failed", category)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	results, ok := data["results"].([]interface{})
	if !ok {
		results = []interface{}{}
	}
	items := []Item{}
	for _, v := range results {
		itemMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		name, _ := itemMap["name"].(string)
		index, _ := itemMap["index"].(string)
		items = append(items, Item{
			Name:  name,
			Index: index,
			URI:   fmt.Sprintf("resource://dnd/item/%s/%s", category, index),
		})
	}
	return map[string]interface{}{
		"category": category,
		"count":    len(items),
		"items":    items,
		"source":   "D&D 5e API (www.dnd5eapi.co)",
	}, nil
}

func FetchItem(client *http.Client, apiBaseURL, category, index string) (map[string]interface{}, error) {
	resp, err := client.Get(fmt.Sprintf("%s/%s/%s", apiBaseURL, category, index))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Item '%s' not found in category '%s' or API request failed", index, category)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	return data, nil
}
