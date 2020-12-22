package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
)

const configFileLocation = "/etc/monzo-ynab/config.json"

var configFileCache configMap

type configMap map[string]string

// NewConfig generates a new configuration
func NewConfig() Config {
	return Config{
		YNABToken:        requiredString("YNAB_TOKEN"),
		YNABAccountID:    requiredString("ACCOUNT_ID"),
		YNABBudgetID:     requiredString("BUDGET_ID"),
		MonzoAccountID:   requiredString("MONZO_ACCOUNT_ID"),
		MonzoAccessToken: requiredString("MONZO_ACCESS_TOKEN"),
		BaseURL:          defaultString("BASE_URL", "http://localhost:8080"),
	}
}

// Config contains all app configuration
type Config struct {
	YNABToken     string `json:"YNAB_TOKEN"` // YNAB Personal Access Token
	YNABAccountID string `json:"ACCOUNT_ID"` // YNAB Account ID to sync
	YNABBudgetID  string `json:"BUDGET_ID"`  // YNAB Budget ID to sync

	MonzoAccountID   string `json:"MONZO_ACCOUNT_ID"`   // Monzo's account ID
	MonzoAccessToken string `json:"MONZO_ACCESS_TOKEN"` // Access Token for Monzo

	BaseURL string `json:"BASE_URL"` // The sync app's URL
}

// Set sets a value in the config
func (c *Config) Set(key string, value string) error {
	point := reflect.ValueOf(c)
	struc := point.Elem()
	field := struc.FieldByName(key)
	if field.IsValid() && field.CanSet() {
		field.SetString(value)
	} else {
		return fmt.Errorf("Could not set %s", key)
	}

	err := c.Persist()
	if err != nil {
		return err
	}
	return nil
}

// Persist saves the config to the config file in `configFileLocation`.
func (c Config) Persist() error {
	bytes, err := json.Marshal(c)
	if err != nil {
		return err
	}

	ioutil.WriteFile(configFileLocation, bytes, 0644)
	return nil
}

// defaultString gets a single value, and uses the default if it isn't found.
func defaultString(key string, dflt string) string {
	val := getValue(key)
	if val == "" {
		return dflt
	}
	return val
}

// requiredString gets a single value and fails if it isn't found
func requiredString(key string) string {
	val := getValue(key)
	if val == "" {
		log.Fatalf("Config value %s is required but was not found", key)
	}
	return val
}

// getValue gets a single value from configuration
func getValue(key string) string {
	fileData := loadConfigFile()

	var val string
	val = fileData[key]
	if val != "" {
		return val
	}

	val = os.Getenv(key)
	return val
}

// loadConfigFile get the config file and returns a configMap
func loadConfigFile() configMap {
	if len(configFileCache) > 0 {
		return configFileCache
	}

	dat, err := ioutil.ReadFile(configFileLocation)
	if err != nil {
		return configMap{} // File doesn't exist/isn't readable. No issue.
	}

	var cm configMap
	err = json.Unmarshal(dat, &cm)
	if err != nil {
		log.Fatalf("Config file unreadable: %s", err)
	}

	configFileCache = cm
	return cm
}
