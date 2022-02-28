package cfg

import (
	"fmt"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func fetchEnv() UserConfig {
	user_config := UserConfig{}
	user_config.VERBOSE = os.Getenv("VERBOSE") == "true"
	user_config.STATUS = os.Getenv("STATUS") == "true"
	user_config.ADDRESS = os.Getenv("ADDRESS")

	endpoint := os.Getenv("ENDPOINT")

	listen_port := os.Getenv("LISTEN_PORT")
	if listen_port != "" {
		// running in server mode
		listen_port, err := strconv.Atoi(listen_port)
		if err != nil {
			log.WithField("Error", err).Fatal("Error parsing LISTEN_PORT as an integer!")
		}
		user_config.SERVER.LISTEN_PORT = uint16(listen_port)
		user_config.SERVER.PRIVATE_KEY = os.Getenv("PRIVATE_KEY")
		user_config.SERVER.CLIENTS = make(map[string]Clients)
		for i := 1; ; i++ {
			allowed_ip := os.Getenv(fmt.Sprintf("ALLOWED_IP_%d", i))
			if allowed_ip == "" {
				break
			}
			client_name := fmt.Sprintf("client_%d", i)
			dns := os.Getenv(fmt.Sprintf("DNS_%d", i))
			public_key := os.Getenv(fmt.Sprintf("PUBLIC_KEY_%d", i))
			preshared_key := os.Getenv(fmt.Sprintf("PRESHARED_KEY_%d", i))
			persistent_keep_alive := os.Getenv(fmt.Sprintf("PERSISTENT_KEEP_ALIVE_%d", i))
			var persistent_keep_alive_int uint8 = 0
			if persistent_keep_alive != "" {
				tmp, err := strconv.Atoi(persistent_keep_alive)
				if err != nil {
					log.WithField("Error", err).Fatal("Error parsing PERSISTENT_KEEP_ALIVE as an integer!")
				}
				persistent_keep_alive_int = uint8(tmp)
			}
			port_fwds := make([]PortForward, 0)
			for j := 1; ; j++ {
				port_fwd_mode := os.Getenv(fmt.Sprintf("PORT_FORWARD_MODE_%d_%d", i, j))
				if port_fwd_mode == "" {
					break
				}
				port_fwd_port := os.Getenv(fmt.Sprintf("PORT_FORWARD_PORT_%d_%d", i, j))
				port_fwd_port_int, err := strconv.Atoi(port_fwd_port)
				if err != nil {
					log.WithField("Error", err).Fatal(fmt.Sprintf("Error parsing PORT_FORWARD_PORT_%d_%d as an integer!", i, j))
				}
				port_fwds = append(port_fwds, PortForward{
					MODE: port_fwd_mode,
					PORT: uint16(port_fwd_port_int),
				})
			}

			client := Clients{
				ALLOWED_IP:            allowed_ip,
				DNS:                   dns,
				PUBLIC_KEY:            public_key,
				PRESHARED_KEY:         preshared_key,
				PERSISTENT_KEEP_ALIVE: persistent_keep_alive_int,
				PORT_FORWARD:          port_fwds,
				ENDPOINT:              endpoint,
			}
			user_config.SERVER.CLIENTS[client_name] = client
		}
	}

	return user_config
}
