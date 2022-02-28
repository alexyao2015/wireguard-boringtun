package cfg

import (
	"os"

	"github.com/alexyao2015/wireguard-boringtun/helpers"

	"gopkg.in/yaml.v2"
)

func fetchYaml() UserConfig {
	user_config := UserConfig{}
	data, _ := helpers.WrapError(
		func() ([]byte, error) {
			return os.ReadFile(helpers.Config_file_path)
		},
		true,
		helpers.Config_file_path,
	)

	_, _ = helpers.WrapError(
		func() ([]byte, error) {
			err := yaml.Unmarshal(data, &user_config)
			return nil, err
		},
		true,
		helpers.Config_file_path,
	)
	return user_config
}
