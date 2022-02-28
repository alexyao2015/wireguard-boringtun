package cfg

import (
	"os"

	log "github.com/sirupsen/logrus"
)

var configFile = "config.yaml"

// This returns the user config and a bool representing isClient .
func ParseUser() (UserConfig, bool) {
	// if file exists do yaml else do env
	var user_config UserConfig
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		log.Info("No config file found, using environment variables")
		user_config = fetchEnv()
	} else {
		log.Info("Config file found, using config file")
		user_config = fetchYaml()
	}
	fromClientConf(&user_config)
	fromServerKey(&user_config)

	if user_config.VERBOSE {
		log.SetLevel(log.DebugLevel)
	}

	is_client := validateCfg(&user_config)
	return user_config, is_client
}
