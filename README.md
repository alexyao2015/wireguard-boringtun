# Wireguard with Boringtun

A config can be bind mounted to /etc/wireguard/wg0.conf or
environment variables can be used for configuration.

## Environment Variables

Server mode is enabled if LISTEN_PORT is defined.

If a public/private key is not defined, a new one will be generated and saved to /config.
The client config qr code will also be shown in the logs.

### General
- `DISABLE_IP6TABLES`: Disable IPv6 iptables rules when set to `true`.
- `VERBOSE`: Enable verbose logging when set to `true`. This dumps generated configs to stdout.
- `STATUS`: Enable status logging every 30 seconds when set to `true`, using `wg show`.

### `[Interface]`
- `ADDRESS`: The address of the client/subnet of the server.
- `PRIVATE_KEY`: The private key of the server/client.

#### Client
- `DNS_x`: The DNS server to be used for the client.
When used in server mode, this adds a `DNS` entry to the generated client configuration.

#### Server
- `LISTEN_PORT`: The port to listen on. (server only)

### `[Peer]`
All these variables may be repeated multiple times when in server mode. They are noted
with `x` being a number that is incremented by 1 and starting at 1.
For example, your first peer would use `PUBLIC_KEY_1`. When in client mode, the
variables are not numbered, so you would use `PUBLIC_KEY`.

- `PUBLIC_KEY_x`: The public key of the peer. (required for client)
- `PRESHARED_KEY_x`: The preshared key of the peer.
- `ALLOWED_IP_x`: The allowed IPs of the peer, required for client and server.
- `ENDPOINT`: The endpoint of the peer, excluding the port.
For server, this is used during client config generation. (required for client)

#### Server
- `PERSISTENT_KEEP_ALIVE_x`: The persistent key alive of the peer. (server only)
- `PORT_FORWARD_PORT_x_y`: The port to forward to the peer. (server only)
- `PORT_FORWARD_MODE_x_y`: The protocol to forward for the peer, `tcp` or `udp`. (server only)
- `PORT_FORWARD_ADDRESS_x_y`: The address to forward to the peer. (server only)
- `PORT_FORWARD_VERSION_x_y`: The protocol version of the address to forward, `4` or `6`. (server only)

