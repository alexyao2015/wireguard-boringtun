package helpers

import (
	"flag"

	log "github.com/sirupsen/logrus"
)

var (
	DryRun          = flag.Bool("dry-run", false, "Run without outputting files.")
	Verbose         = flag.Bool("verbose", false, "Enable verbose logging.")
	ConfigDir       = flag.String("config-dir", ".", "Location of config directory.")
	ClientsDir      = flag.String("clients-dir", "clients", "Location of clients directory relative to config dir.")
	WireguardConfig = flag.String("wireguard-config", "wg0.conf", "Name of wireguard config file.")
)

func Main() {
	flag.Parse()
	log.Info("Settings:")

	log.Info("Dry run: ", *DryRun)
	log.Info("Verbose: ", *Verbose)
	log.Info("Config dir: ", *ConfigDir)
	log.Info("Clients dir: ", *ClientsDir)
	log.Info("Wireguard config: ", *WireguardConfig)

	if *Verbose {
		log.SetLevel(log.DebugLevel)
	}
}
