package cfg

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	ini "gopkg.in/ini.v1"
)

// parse private key from clients folder
func fromClientConf(cfg *UserConfig) {
	if _, err := os.Stat("clients"); os.IsNotExist(err) {
		return
	}
	dirListings, err := os.ReadDir("clients")
	if err != nil {
		log.WithField("Error", err).Fatal("Error reading clients directory!")
	}
	for _, item := range dirListings {
		if item.IsDir() {
			continue
		}
		if !strings.HasSuffix(item.Name(), ".conf") {
			continue
		}
		client_name := strings.TrimSuffix(item.Name(), ".conf")
		log.Debug("Found client config: " + client_name)
		_, exists := cfg.SERVER.CLIENTS[client_name]
		if !exists {
			log.Warn("Client config not found in server config! Skipping...")
			continue
		}
		ini_cfg, err := ini.Load("clients/" + item.Name())
		if err != nil {
			log.WithField("Error", err).Fatal(fmt.Sprintf("Error reading client config: %s", item.Name()))
			continue
		}
		private_key := ini_cfg.Section("Interface").Key("PrivateKey").String()
		if private_key == "" {
			log.Warn("Client config has no private key! Skipping...")
			continue
		} else {
			log.Debug("Found private key for client: " + client_name)
			client_config := cfg.SERVER.CLIENTS[client_name]
			client_config.PRIVATE_KEY = private_key
			cfg.SERVER.CLIENTS[client_name] = client_config
		}
	}
}

func fromServerKey(cfg *UserConfig) {
	if _, err := os.Stat("server.key"); os.IsNotExist(err) {
		return
	}
	if cfg.SERVER.PRIVATE_KEY != "" {
		return
	}
	data, err := os.ReadFile("server.key")
	if err != nil {
		log.WithField("Error", err).Fatal("Error reading server.key file!")
	}
	cfg.SERVER.PRIVATE_KEY = string(data)
}
