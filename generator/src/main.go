package main

import (
	"github.com/alexyao2015/wireguard-boringtun/cmd"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: true,
	})
	log.SetLevel(log.DebugLevel)
}

func main() {
	cmd.Main()
}
