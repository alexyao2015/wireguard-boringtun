package cfg

import (
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	log "github.com/sirupsen/logrus"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

var validate = validator.New()

func init() {
	validate.RegisterValidation("notblank", validators.NotBlank)
}

// validates and fills in defaults for the user config
// Creates private and public keys for clients if not defined
func validateCfg(cfg UserConfig) bool {
	var is_client = cfg.CLIENT.PRIVATE_KEY != ""
	var is_server = cfg.SERVER.LISTEN_PORT != 0

	if is_client && is_server {
		log.Error("Cannot have both client and server config keys!")
		log.Error("Remove either client or server keys!")
		log.Fatal("Exiting...")
	} else if !(is_client || is_server) {
		log.Error("You must define either the client or server config key!")
		log.Fatal("Exiting...")
	}

	var errors []error
	if is_client {
		log.Debug("Running in client mode")
		errors = append(errors, validate.Struct(cfg.CLIENT))
	} else {
		log.Debug("Running in server mode")

		// create server private key if one doesn't exist
		if cfg.SERVER.PRIVATE_KEY == "" {
			log.Info("Generating server private key")
			private_key, err := wgtypes.GeneratePrivateKey()
			if err != nil {
				log.WithField("Error", err).Fatal("Error generating private key!")
			}
			cfg.SERVER.PRIVATE_KEY = string(private_key[:])
			log.Info("Saving server private key to server.key")
			os.WriteFile("server.key", []byte(cfg.SERVER.PRIVATE_KEY), 0744)
		}

		errors = append(errors, validate.Struct(cfg.SERVER))
		if !(len(cfg.SERVER.CLIENTS) > 0) {
			log.Fatal("When in server mode, clients must be greater than 0!")
		}
		// validate each client
		for name, client := range cfg.SERVER.CLIENTS {
			// Generate a new pub/private key if one is not provided
			if client.PUBLIC_KEY == "" {
				if client.PRIVATE_KEY == "" {
					log.Info("Generating private key for ", name)
					private_key, err := wgtypes.GeneratePrivateKey()
					if err != nil {
						log.WithField("Error", err).Fatal("Error generating private key!")
					}
					psk, err := wgtypes.GenerateKey()
					if err != nil {
						log.WithField("Error", err).Fatal("Error generating PSK!")
					}
					client.PRIVATE_KEY = string(private_key.String())
					client.PRESHARED_KEY = string(psk.String())
					client.GENERATED = true
				}
				log.Debug("Generating public key for ", name)
				key, err := wgtypes.ParseKey(client.PRIVATE_KEY)
				if err != nil {
					log.WithField("Error", err).Fatal("Error parsing private key!")
				}
				client.PUBLIC_KEY = string(key.PublicKey().String())
				cfg.SERVER.CLIENTS[name] = client
			}

			errors = append(errors, validate.Struct(client))

			// if port forwards are listed, validate them as well
			if len(client.PORT_FORWARD) > 0 {
				for _, port_fwd := range client.PORT_FORWARD {
					errors = append(errors, validate.Struct(port_fwd))
				}
			}
		}
	}

	for _, err := range errors {
		if err != nil {
			log.WithField("Error", err).Fatal("Error validating config!")
		}
	}

	return is_client
}
