package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	MongoURI string `json:"mongo_uri"`
	Epgdump  string
	Recpt1   string
	Channels map[string][]struct {
		Ch  string
		Sid string
	}
}

func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var config Config
	if err := json.NewDecoder(f).Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
