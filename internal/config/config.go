package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Settings loaded from the configuration file.
type Settings struct {
	// Application's port
	Port int

	// Database connection string
	DatabaseConnString string
}

// ReadFile loads the settings from a configuration file.
// Source https://github.com/plifk/market/blob/c97fb9123dc35141c0f08ae9c411b6bd1a603fb7/internal/config/config.go#L71
func ReadFile(path string) (s Settings, err error) {
	f, err := os.Open(path) // #nosec
	if err != nil {
		return s, err
	}
	if err = json.NewDecoder(f).Decode(&s); err != nil {
		return s, fmt.Errorf("cannot load gaivota configuration: %w", err)
	}
	return s, nil
}
