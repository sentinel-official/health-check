package models

import (
	"github.com/sentinel-official/health-check/types"
)

type Session struct {
	ID             uint64       `json:"id,omitempty" bson:"id,omitempty"`
	SubscriptionID uint64       `json:"subscription_id,omitempty" bson:"subscription_id,omitempty"`
	Status         types.Status `json:"status,omitempty" bson:"status,omitempty"`
}
