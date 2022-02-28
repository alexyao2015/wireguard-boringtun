package main

import (
	"github.com/alexyao2015/wireguard-boringtun/cmd"
	"github.com/alexyao2015/wireguard-boringtun/helpers"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: true,
	})
}

func main() {
	helpers.Main()
	cmd.Main()
}
