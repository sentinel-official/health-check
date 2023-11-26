package wireguard

const (
	configTemplate = `
[Interface]
Address = {{ .Interface.Address }}
PrivateKey = {{ .Interface.PrivateKey }}
DNS = {{ .Interface.DNS }}

[Peer]
PublicKey = {{ .Peer.PublicKey }}
AllowedIPs = {{ .Peer.AllowedIPs }}
Endpoint = {{ .Peer.Endpoint }}
PersistentKeepalive = {{ .Peer.PersistentKeepalive }}
`
)
