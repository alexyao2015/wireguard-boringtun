package cmd

import (
	"fmt"
	"strings"

	"github.com/alexyao2015/wireguard-boringtun/cmd/cfg"
	"github.com/alexyao2015/wireguard-boringtun/resources"
)

var client, _ = resources.EmbedFS.ReadFile("client.conf")

func clientConfigBackend(cfg cfg.Clients, address string, allowed_ip string) string {
	var dns string = ""
	if cfg.DNS != "" {
		dns = "\nDNS = " + cfg.DNS
	}
	var preshared_key string = ""
	if cfg.PRESHARED_KEY != "" {
		preshared_key = "\nPresharedKey = " + cfg.PRESHARED_KEY
	}
	var persistent_keep_alive string = ""
	if cfg.PERSISTENT_KEEP_ALIVE != 0 {
		persistent_keep_alive = "\nPersistentKeepalive = " + fmt.Sprintf("%d", cfg.PERSISTENT_KEEP_ALIVE)
	}

	replacer := strings.NewReplacer(
		"${ADDRESS}", address,
		"${PRIVATE_KEY}", cfg.PRIVATE_KEY,
		"${DNS}", dns,
		"${PUBLIC_KEY}", cfg.PUBLIC_KEY,
		"${PRESHARED_KEY}", preshared_key,
		"${ENDPOINT}", cfg.ENDPOINT,
		"${ALLOWED_IP}", allowed_ip,
		"${PERSISTENT_KEEP_ALIVE}", persistent_keep_alive,
	)
	return replacer.Replace(string(client))
}

func genClient(user_cfg cfg.UserConfig) string {
	return clientConfigBackend(user_cfg.CLIENT, user_cfg.ADDRESS, user_cfg.CLIENT.ALLOWED_IP)
}
