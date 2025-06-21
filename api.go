package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"
)

// endpoint represents a specific API endpoint in the D&D 5e API.
type endpoint string

const (
	apiBaseURL            = "https://www.dnd5eapi.co/api"
	requestTimeoutSeconds = 10

	abilityScores       endpoint = "ability-scores"
	alignments          endpoint = "alignments"
	backgrounds         endpoint = "backgrounds"
	classes             endpoint = "classes"
	conditions          endpoint = "conditions"
	damageTypes         endpoint = "damage-types"
	equipment           endpoint = "equipment"
	equipmentCategories endpoint = "equipment-categories"
	feats               endpoint = "feats"
	features            endpoint = "features"
	languages           endpoint = "languages"
	magicItems          endpoint = "magic-items"
	magicSchools        endpoint = "magic-schools"
	monsters            endpoint = "monsters"
	proficiencies       endpoint = "proficiencies"
	races               endpoint = "races"
	ruleSections        endpoint = "rule-sections"
	rules               endpoint = "rules"
	skills              endpoint = "skills"
	spells              endpoint = "spells"
	subclasses          endpoint = "subclasses"
	subraces            endpoint = "subraces"
	traits              endpoint = "traits"
	weaponProperties    endpoint = "weapon-properties"
)

// listResponse defines the structure of the response for a list endpoint.
type listResponse struct {
	Count   int                      `json:"count"`
	Results []map[string]interface{} `json:"results"`
}

// fetchAPIItem fetches a single item by endpoint and item from the D&D 5e API.
func fetchAPIItem(client *http.Client, endpoint endpoint, item string) (map[string]interface{}, error) {
	resp, err := client.Get(fmt.Sprintf("%s/%s/%s", apiBaseURL, endpoint, url.PathEscape(item)))
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
	return data, nil
}

// fetchByName fetches an item by name from the D&D 5e API and unmarshals it into the provided variable.
func fetchByName(client *http.Client, e endpoint, name string, v any) error {
	logrus.WithFields(logrus.Fields{"endpoint": e, "name": name}).Debug("fetchByName called")
	index := toKebabCase(name)
	data, err := fetchAPIItem(client, e, index)
	if err != nil {
		logrus.WithError(err).Error("fetchAPIItem failed in fetchByName")
		return err
	}
	b, err := json.Marshal(data)
	if err != nil {
		logrus.WithError(err).Error("Marshal failed in fetchByName")
		return err
	}
	if err := json.Unmarshal(b, v); err != nil {
		logrus.WithError(err).Error("Unmarshal failed in fetchByName")
		return err
	}
	logrus.WithField("name", name).Debug("fetchByName succeeded")
	return nil
}

// fetchAPIList fetches a list of items for the given endpoint from the D&D 5e API.
func fetchAPIList(client *http.Client, endpoint endpoint, filter string) (listResponse, error) {
	url := fmt.Sprintf("%s/%s", apiBaseURL, endpoint)
	if filter != "" {
		url = fmt.Sprintf("%s?%s", url, filter)
	}
	resp, err := client.Get(url)
	if err != nil {
		return listResponse{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return listResponse{}, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return listResponse{}, err
	}

	var data listResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return listResponse{}, err
	}
	return data, nil
}

// fetchList fetches a list of items from the D&D 5e API and unmarshals it into the provided variable.
func fetchList(client *http.Client, e endpoint, v any, filter string) error {
	logrus.WithFields(logrus.Fields{"endpoint": e, "filter": filter}).Debug("fetchList called")
	spells, err := fetchAPIList(client, e, filter)
	if err != nil {
		logrus.WithError(err).Error("fetchAPIList failed in fetchList")
		return err
	}
	jsonData, err := json.Marshal(spells.Results)
	if err != nil {
		logrus.WithError(err).Error("Marshal failed in fetchList")
		return err
	}
	if err := json.Unmarshal(jsonData, v); err != nil {
		logrus.WithError(err).Error("Unmarshal failed in fetchList")
		return err
	}
	logrus.Debug("fetchList succeeded")
	return nil
}
