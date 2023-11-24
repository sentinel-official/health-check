package types

import (
	"crypto/tls"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	NodeTypeUnknown NodeType = 0x00 + iota
	NodeTypeWireGuard
	NodeTypeV2Ray
)

type NodeType byte

func NewNodeTypeFromString(v string) NodeType {
	v = strings.TrimSpace(v)
	v = strings.ToLower(v)

	switch v {
	case "wire_guard":
		return NodeTypeWireGuard
	case "v2ray":
		return NodeTypeV2Ray
	default:
		return NodeTypeUnknown
	}
}

func NewNodeTypeFromUInt64(v uint64) NodeType {
	switch v {
	case 1:
		return NodeTypeWireGuard
	case 2:
		return NodeTypeV2Ray
	default:
		return NodeTypeUnknown
	}
}

func (z NodeType) String() string {
	switch z {
	case NodeTypeWireGuard:
		return "wire_guard"
	case NodeTypeV2Ray:
		return "v2ray"
	default:
		return "unknown"
	}
}

type (
	NodeInfo struct {
		Type uint64 `json:"type"`
	}
)

func FetchNewNodeInfo(remoteURL string, timeout time.Duration) (*NodeInfo, error) {
	urlPath, err := url.JoinPath(remoteURL, "status")
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: timeout,
	}

	resp, err := client.Get(urlPath)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var m map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, err
	}

	buf, err := json.Marshal(m["result"])
	if err != nil {
		return nil, err
	}

	var v NodeInfo
	if err := json.Unmarshal(buf, &v); err != nil {
		return nil, err
	}

	return &v, nil
}
