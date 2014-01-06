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

func LoadConfig(path *string) (Config, error) {
	var config Config
	f, err := os.Open(*path)
	if err != nil {
		return config, err
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	return config, dec.Decode(&config)
}
