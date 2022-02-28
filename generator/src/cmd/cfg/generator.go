package cfg

import (
	"github.com/alexyao2015/wireguard-boringtun/helpers"
	"github.com/alexyao2015/wireguard-boringtun/resources"
	"gopkg.in/yaml.v2"
)

type DefaultRules struct {
	DEFAULT      []string
	PORT_FORWARD []string
}

type GeneratorConfig struct {
	IPV4 DefaultRules
	IPV6 DefaultRules
}

var rules, _ = resources.EmbedFS.ReadFile("rules.yaml")

// Parses embedded iptables rules.yaml file
func ParseGenerator() GeneratorConfig {
	config := GeneratorConfig{}
	helpers.WrapError(
		func() ([]byte, error) {
			err := yaml.Unmarshal([]byte(rules), &config)
			return nil, err
		},
		true,
		"rules.yaml",
	)
	return config
}
