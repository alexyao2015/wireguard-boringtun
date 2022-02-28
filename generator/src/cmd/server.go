package cmd

import (
	"fmt"
	"strings"

	"github.com/alexyao2015/wireguard-boringtun/cmd/cfg"
	"github.com/alexyao2015/wireguard-boringtun/resources"
)

var server, _ = resources.EmbedFS.ReadFile("server.conf")
var server_peer, _ = resources.EmbedFS.ReadFile("server_peer.conf")

func addRules(inRules string, genRules []string) string {
	for _, rule := range genRules {
		inRules += rule + ";"
	}
	return inRules
}

func genServer(userCfg cfg.UserConfig) string {
	var rules_temp string = ""
	var found_ipv4, foundIPv6 bool = false, false
	splitAddr := strings.Split(userCfg.ADDRESS, ",")

	// Add default iptable rules
	for _, addr := range splitAddr {
		is_ipv6 := strings.Contains(addr, ":")
		if !found_ipv4 && !is_ipv6 {
			rules_temp = addRules(rules_temp, generator_config.IPV4.DEFAULT)
			found_ipv4 = true
		} else if !foundIPv6 && is_ipv6 {
			rules_temp = addRules(rules_temp, generator_config.IPV6.DEFAULT)
			foundIPv6 = true
		}
	}

	// add port forward rules
	for _, client := range userCfg.SERVER.CLIENTS {
		split_addr := strings.Split(client.ALLOWED_IP, ",")
		for _, addr := range split_addr {
			addr = strings.TrimSpace(addr)

			is_ipv6 := strings.Contains(addr, ":")
			for _, port_fwd := range client.PORT_FORWARD {
				var new_rules string
				if is_ipv6 {
					new_rules = addRules("", generator_config.IPV6.PORT_FORWARD)
				} else {
					new_rules = addRules("", generator_config.IPV4.PORT_FORWARD)
				}
				replacer := strings.NewReplacer(
					"${protocol}", port_fwd.MODE,
					"${port_number}", fmt.Sprintf("%d", port_fwd.PORT),
					"${address}", addr,
				)
				new_rules = replacer.Replace(new_rules)
				rules_temp += new_rules
			}
		}
	}

	rules_temp = strings.Replace(
		rules_temp,
		"${listen_port}", fmt.Sprintf("%d", userCfg.SERVER.LISTEN_PORT),
		-1,
	)

	post_up := strings.Replace(
		rules_temp,
		"${add_remove}", "A",
		-1,
	)

	post_down := strings.Replace(
		rules_temp,
		"${add_remove}", "D",
		-1,
	)

	var server_peer_cfg string = ""
	// create peers
	for name, peer := range userCfg.SERVER.CLIENTS {
		var preshared_key, persistent_keep_alive string = "", ""
		if peer.PRESHARED_KEY != "" {
			preshared_key = "\nPresharedKey = " + peer.PRESHARED_KEY
		}
		if peer.PERSISTENT_KEEP_ALIVE != 0 {
			persistent_keep_alive = "\nPersistentKeepalive = " + fmt.Sprintf("%d", peer.PERSISTENT_KEEP_ALIVE)
		}

		replacer := strings.NewReplacer(
			"${COMMENT_STRING}", "# "+name,
			"${PUBLIC_KEY}", peer.PUBLIC_KEY,
			"${ALLOWED_IP}", peer.ALLOWED_IP,
			"${PRESHARED_KEY}", preshared_key,
			"${PERSISTENT_KEEP_ALIVE}", persistent_keep_alive,
		)
		server_peer_cfg += replacer.Replace(string(server_peer))
		server_peer_cfg += "\n"
	}

	server_peer_cfg = strings.TrimRight(server_peer_cfg, "\n")

	replacer := strings.NewReplacer(
		"${ADDRESS}", userCfg.ADDRESS,
		"${LISTEN_PORT}", fmt.Sprintf("%d", userCfg.SERVER.LISTEN_PORT),
		"${PRIVATE_KEY}", userCfg.SERVER.PRIVATE_KEY,
		"${POST_UP}", post_up,
		"${POST_DOWN}", post_down,
		"${PEERS}", server_peer_cfg,
	)
	server_cfg := replacer.Replace(string(server))

	return server_cfg
}

func genServerClients(userCfg cfg.UserConfig) map[string]string {
	client_configs := make(map[string]string)
	for name, client_cfg := range userCfg.SERVER.CLIENTS {
		if !client_cfg.GENERATED {
			continue
		}

		// set a default allowed_ip for the config if it doesn't exist
		allowed_ip := client_cfg.ALLOWED_IP_CONFIG
		if allowed_ip == "" {
			if strings.Contains(client_cfg.ALLOWED_IP, ":") {
				allowed_ip = "0.0.0.0/0, ::/0"
			} else {
				allowed_ip = "0.0.0.0/0"
			}
		}
		// for the client config, the server side allowed ip becomes the client side address
		str_cfg := clientConfigBackend(client_cfg, client_cfg.ALLOWED_IP, allowed_ip)
		client_configs[name] = str_cfg
	}
	return client_configs
}
