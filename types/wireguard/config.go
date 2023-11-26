package wireguard

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"text/template"
)

type (
	interfaceConfig struct {
		Address    string
		PrivateKey string
		DNS        string
	}

	peerConfig struct {
		PublicKey           string
		AllowedIPs          string
		Endpoint            string
		PersistentKeepalive int64
	}
)

type Config struct {
	Interface interfaceConfig
	Peer      peerConfig
}

func (z *Config) WriteToFile(path string) error {
	t, err := template.New("WireGuardConfig").Parse(configTemplate)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, z); err != nil {
		return err
	}

	return os.WriteFile(path, buf.Bytes(), 0600)
}

func NewConfigFromBytes(privKey, buf []byte) (*Config, error) {
	if len(buf) != 58 {
		return nil, fmt.Errorf("incorrect buffer size %d", len(buf))
	}

	return &Config{
		Interface: interfaceConfig{
			Address:    fmt.Sprintf("%s/32", net.IP(buf[0:4])),
			PrivateKey: NewKey(privKey).String(),
			DNS:        "1.0.0.1, 1.1.1.1",
		},
		Peer: peerConfig{
			PublicKey:           NewKey(buf[26:58]).String(),
			AllowedIPs:          "0.0.0.0/0",
			Endpoint:            fmt.Sprintf("%s:%d", net.IP(buf[20:24]), binary.BigEndian.Uint16(buf[24:26])),
			PersistentKeepalive: 15,
		},
	}, nil
}
