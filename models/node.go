package models

import (
	"github.com/sentinel-official/health-check/types"
)

type Node struct {
	Addr      string         `json:"addr,omitempty" bson:"addr,omitempty"`
	RemoteURL string         `json:"remote_url,omitempty" bson:"remote_url,omitempty"`
	Status    types.Status   `json:"status,omitempty" bson:"status,omitempty"`
	Type      types.NodeType `json:"type,omitempty" bson:"type,omitempty"`
}
