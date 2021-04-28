package main

import (
	"encoding/base64"
	"encoding/json"
	"os"
)

type Config struct {
	Pattern string `json:"pattern"`
	Replace string `json:"replace"`
}

func get(envVarName string) (ret []Config) {
	decoded, decodeErr := base64.StdEncoding.DecodeString(os.Getenv(envVarName))
	if decodeErr != nil {
		log.Error("Could not decode config. Ensure config is provided in base64 format. %v\n", decodeErr)
		return
	}
	jsonParseErr := json.Unmarshal(decoded, &ret)
	if jsonParseErr != nil {
		log.Error("Could not parse json config. Ensure config is valid json. %v\n", jsonParseErr)
		return nil
	}
	return ret
}
