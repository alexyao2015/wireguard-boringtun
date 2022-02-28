package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/alexyao2015/wireguard-boringtun/cmd/cfg"
	"github.com/alexyao2015/wireguard-boringtun/helpers"
)

var generator_config = cfg.ParseGenerator()

func Main() {
	usr_cfg, is_client := cfg.ParseUser()

	var generated_config string

	if is_client {
		generated_config = genClient(usr_cfg)
	} else {
		generated_config = genServer(usr_cfg)
		client_config := genServerClients(usr_cfg)
		err := helpers.DryRunWrapper(func() error { return os.MkdirAll(helpers.Clients_path, 755) })
		if err != nil {
			log.WithField("Error", err).Fatal("Could not create clients directory!")
		}

		for name, cfg_str := range client_config {
			client_path := filepath.Join(helpers.Clients_path, name+".conf")
			err := helpers.DryRunWrapper(func() error {
				return os.WriteFile(client_path, []byte(cfg_str), 0744)
			})
			if err != nil {
				log.WithField("Error", err).Fatal("Could not write client config!")
			}
			log.Info(fmt.Sprintf("Generated config for client %s", name))
			qr_code, err := helpers.RunCmd(cfg_str, "qrencode", "-t", "ansiutf8", "-l", "L")
			if err != nil {
				log.WithField("Error", err).Warn("Could not generate QR code! Ensure qrencode is installed.")
			} else {
				fmt.Println(qr_code)
			}
		}
	}
	generated_config = fmt.Sprintf(
		"%s\n%s\n%s",
		"# This file is generated and will be automatically overwritten.",
		"# Do not edit this file, edit the config.yaml file instead.",
		generated_config,
	)
	helpers.DryRunWrapper(func() error {
		return os.WriteFile(helpers.Wireguard_config_path, []byte(generated_config), 0744)
	})
	log.Info("Wrote config to ", helpers.Wireguard_config_path)
}
