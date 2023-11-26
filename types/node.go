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
	NodeTypeUnspecified NodeType = 0x00 + iota
	NodeTypeWireGuard
	NodeTypeV2Ray
)

type NodeType byte

func NewNodeTypeFromString(v string) NodeType {
	v = strings.TrimSpace(v)
	v = strings.ToLower(v)

	switch v {
	case "wireguard":
		return NodeTypeWireGuard
	case "v2ray":
		return NodeTypeV2Ray
	default:
		return NodeTypeUnspecified
	}
}

func NewNodeTypeFromUInt64(v uint64) NodeType {
	switch v {
	case 1:
		return NodeTypeWireGuard
	case 2:
		return NodeTypeV2Ray
	default:
		return NodeTypeUnspecified
	}
}

func (z NodeType) String() string {
	switch z {
	case NodeTypeWireGuard:
		return "wireguard"
	case NodeTypeV2Ray:
		return "v2ray"
	default:
		return "unspecified"
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
