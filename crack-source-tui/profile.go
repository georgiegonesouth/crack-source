package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type profile struct {
	Tokens map[string]string `json:"tokens"`
}

func profilePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".crack-source", "profiles", "default.json")
}

func loadProfile() profile {
	p := profile{Tokens: map[string]string{}}
	data, err := os.ReadFile(profilePath())
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
	path := profilePath()
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return
	}
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile(path, data, 0600)
}
