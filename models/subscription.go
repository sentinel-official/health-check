package models

import (
	"github.com/sentinel-official/health-check/types"
)

type Subscription struct {
	ID       uint64       `json:"id,omitempty" bson:"id,omitempty"`
	NodeAddr string       `json:"node_addr,omitempty" bson:"node_addr,omitempty"`
	Status   types.Status `json:"status,omitempty" bson:"status,omitempty"`
}
