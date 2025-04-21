package functions

import (
	"log"

	"gopkg.in/ini.v1"
)

var (
	configPath = "config.ini"
)

/*
*
* configLoader loads the configuration from config.ini file.
* If the config file cannot be loaded, it logs the error and exits the program.
 */
func ConfigLoader() *ini.File {
	cfg, err := ini.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load config.ini file: %v", err)
	}

	return cfg
}
