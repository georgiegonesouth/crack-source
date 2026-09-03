package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type profile struct {
	Name   string            `json:"name"`
	Tokens map[string]string `json:"tokens"`
}

func profilePath(name string) string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".elluminos", "profiles", name+".json")
}

func loadProfile(name string) profile {
	p := profile{Name: name, Tokens: map[string]string{}}
	data, err := os.ReadFile(profilePath(name))
	if err != nil {
		return p
	}
	_ = json.Unmarshal(data, &p)
	if p.Tokens == nil {
		p.Tokens = map[string]string{}
	}
	return p
}

func (m *Model) commitTokens() {
	for i, k := range tokenKeys {
		m.profile.Tokens[k] = m.tokenInputs[i].Value()
	}
	saveProfile(m.profile)
}

func saveProfile(p profile) {
	path := profilePath(p.Name)
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return
	}
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile(path, data, 0600)
}
