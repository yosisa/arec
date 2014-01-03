package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	configFile := "arec.sample.json"
	config, err := LoadConfig(&configFile)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, config.MongoURI, "mongodb://localhost/arec")
}
