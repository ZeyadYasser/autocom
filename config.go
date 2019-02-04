package main

import (
	"log"
	"os"
	"strconv"
)

// Config configuration
type Config struct {
	Port            int
	User            string
	Password        string
	EngineMaxLevels int
	EngineToLower   bool
	EngineSkipBegin bool
}

func getEnv(key, fallback string) string {
    value, exists := os.LookupEnv(key)
    if !exists {
        value = fallback
    }
    return value
}

func newConfig() *Config {
	port, err := strconv.Atoi(getEnv("PORT", "9696"))
	if err != nil {
		log.Fatalf("Incorrect 'PORT' format (Integer expected): %v\n", err)
	}
	maxLevels, err := strconv.Atoi(getEnv("ENGINE_MAXLEVELS", "1"))
	if err != nil {
		log.Fatalf("Incorrect 'ENGINE_MAXLEVELS' format (Integer expected): %v\n", err)
	}
	toLower, err := strconv.ParseBool(getEnv("ENGINE_TOLOWER", "false"))
	if err != nil {
		log.Fatalf("Incorrect 'ENGINE_TOLOWER' format (Boolean expected): %v\n", err)
	}
	skipBegin, err := strconv.ParseBool(getEnv("ENGINE_SKIPBEGIN", "false"))
	if err != nil {
		log.Fatalf("Incorrect 'ENGINE_SKIPBEGIN' format (Boolean expected): %v\n", err)
	}

	config := &Config {
		Port:            port,
		User: 	         getEnv("AUTOCOM_USER", ""),
		Password:        getEnv("AUTOCOM_PASSWORD", ""),
		EngineMaxLevels: maxLevels,
		EngineToLower:   toLower,
		EngineSkipBegin: skipBegin,
	}
	return config
}
