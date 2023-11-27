package v2ray

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"text/template"

	"github.com/hashicorp/go-uuid"

	"github.com/sentinel-official/health-check/utils"
)

type (
	apiConfig struct {
		Port uint16
	}

	proxyConfig struct {
		Port uint16
	}

	vMessConfig struct {
		Address   string
		ID        string
		Port      uint16
		Transport string
	}
)

type Config struct {
	API   apiConfig
	Proxy proxyConfig
	VMess vMessConfig
}

func NewConfigFromBytes(uid, buf []byte) (*Config, error) {
	if len(buf) != 7 {
		return nil, fmt.Errorf("incorrect buffer size %d", len(buf))
	}

	id, err := uuid.FormatUUID(uid)
	if err != nil {
		return nil, err
	}

	apiPort, err := utils.GetFreeTCPPort()
	if err != nil {
		return nil, err
	}

	return &Config{
		API: apiConfig{
			Port: apiPort,
		},
		Proxy: proxyConfig{
			Port: 1080,
		},
		VMess: vMessConfig{
			Address:   net.IP(buf[0:4]).String(),
			ID:        id,
			Port:      binary.BigEndian.Uint16(buf[4:6]),
			Transport: VMessTransportFromByte(buf[6]),
		},
	}, nil
}

func (z *Config) WriteToFile(path string) error {
	t, err := template.New("V2RayConfig").Parse(configTemplate)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, z); err != nil {
		return err
	}

	return os.WriteFile(path, buf.Bytes(), 0600)
}
