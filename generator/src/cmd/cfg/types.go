package cfg

type PortForward struct {
	MODE string `validate:"eq=udp|eq=tcp"`
	PORT uint16 `validate:"notblank"`
}

type Clients struct {
	GENERATED             bool
	ALLOWED_IP            string `validate:"notblank"`
	ALLOWED_IP_CONFIG     string // used for client config generation only
	DNS                   string
	ENDPOINT              string `validate:"notblank"`
	PERSISTENT_KEEP_ALIVE uint8
	PORT_FORWARD          []PortForward
	PRIVATE_KEY           string
	PRESHARED_KEY         string
	PUBLIC_KEY            string `validate:"notblank"`
}

type UserConfig struct {
	ADDRESS string `validate:"notblank"`
	STATUS  bool
	VERBOSE bool

	SERVER struct {
		LISTEN_PORT uint16 `validate:"notblank"`
		PRIVATE_KEY string `validate:"notblank"`
		CLIENTS     map[string]Clients
	}

	CLIENT Clients
}
