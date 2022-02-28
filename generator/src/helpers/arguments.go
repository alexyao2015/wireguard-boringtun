package helpers

import (
	"flag"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

var (
	Dry_run          = flag.Bool("dry-run", false, "Run without outputting files.")
	cerbose          = flag.Bool("verbose", false, "Enable verbose logging.")
	Config_dir       = flag.String("config-dir", ".", "Location of config directory.")
	Clients_dir      = flag.String("clients-dir", "clients", "Location of clients directory relative to config dir.")
	Wireguard_config = flag.String("wireguard-config", "wg0.conf", "Name of wireguard config file.")
	Server_key       = flag.String("server-key", "server.key", "Name of server key file.")
)

var (
	Config_path           string
	Clients_path          string
	Server_key_path       string
	Wireguard_config_path string
)

func Main() {
	flag.Parse()
	log.Info("Settings:")

	log.Info("Dry run: ", *Dry_run)
	log.Info("Verbose: ", *cerbose)
	log.Info("Config dir: ", *Config_dir)
	log.Info("Clients dir: ", *Clients_dir)
	log.Info("Wireguard config: ", *Wireguard_config)
	log.Info("Server key: ", *Server_key)

	Config_path = *Config_dir
	Clients_path = filepath.Join(Config_path, *Clients_dir)
	Server_key_path = filepath.Join(Config_path, *Server_key)
	Wireguard_config_path = filepath.Join(Config_path, *Wireguard_config)

	if *cerbose {
		log.SetLevel(log.DebugLevel)
	}
}
