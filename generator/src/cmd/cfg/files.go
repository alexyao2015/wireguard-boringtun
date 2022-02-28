package cfg

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexyao2015/wireguard-boringtun/helpers"
	log "github.com/sirupsen/logrus"

	ini "gopkg.in/ini.v1"
)

// parse private key from clients folder
func fromClientConf(cfg *UserConfig) {
	if _, err := os.Stat(helpers.Clients_path); os.IsNotExist(err) {
		return
	}
	dirListings, err := os.ReadDir(helpers.Clients_path)
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
		log.Info("Found client config: " + client_name)
		_, exists := cfg.SERVER.CLIENTS[client_name]
		if !exists {
			log.Warn("Client config not found in server config! Skipping...")
			continue
		}
		item_path := filepath.Join(helpers.Clients_path, item.Name())
		ini_cfg, err := ini.Load(item_path)
		if err != nil {
			log.WithField("Error", err).Fatal(fmt.Sprintf("Error reading client config: %s", item.Name()))
			continue
		}

		client_config := cfg.SERVER.CLIENTS[client_name]
		private_key := ini_cfg.Section("Interface").Key("PrivateKey").String()
		if private_key == "" {
			log.Warn("Client config has no private key! Skipping...")
			continue
		} else {
			log.Info("Found private key for client: " + client_name)
			client_config.PRIVATE_KEY = private_key
		}

		psk := ini_cfg.Section("Interface").Key("PresharedKey").String()
		if psk == "" {
			log.Warn("Client config has no preshared key. Skipping...")
			continue
		} else {
			log.Info("Found preshared key for client: " + client_name)
			client_config.PRESHARED_KEY = psk
		}
		cfg.SERVER.CLIENTS[client_name] = client_config
	}
}

func fromServerKey(cfg *UserConfig) {
	if _, err := os.Stat(helpers.Server_key_path); os.IsNotExist(err) {
		return
	}
	if cfg.SERVER.PRIVATE_KEY != "" {
		return
	}
	data, err := os.ReadFile(helpers.Server_key_path)
	if err != nil {
		log.WithField("Error", err).Fatal(fmt.Sprintf("Error reading %s file!", *helpers.Server_key))
	}
	cfg.SERVER.PRIVATE_KEY = string(data)
}
