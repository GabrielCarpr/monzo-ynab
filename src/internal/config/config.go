package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

const configFileLocation = "/etc/monzo-ynab/config.json"

var configFileCache configMap

type configMap map[string]string

// NewConfig generates a new configuration
func NewConfig() Config {
	return Config{
		YNABToken:     requiredString("YNAB_TOKEN"),
		YNABAccountID: requiredString("ACCOUNT_ID"),
		YNABBudgetID:  requiredString("BUDGET_ID"),
	}
}

// Config contains all app configuration
type Config struct {
	YNABToken     string // YNAB Personal Access Token
	YNABAccountID string // YNAB Account ID to sync
	YNABBudgetID  string // YNAB Budget ID to sync
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
