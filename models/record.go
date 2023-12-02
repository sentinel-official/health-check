package models

import (
	"time"

	hubtypes "github.com/sentinel-official/hub/types"

	"github.com/sentinel-official/health-check/types"
)

type Record struct {
	Addr                    string          `json:"addr" bson:"addr"`
	ClientConfig            []byte          `json:"client_config" bson:"client_config"`
	ConfigExchangeError     string          `json:"config_exchange_error" bson:"config_exchange_error"`
	ConfigExchangeTimestamp time.Time       `json:"config_exchange_timestamp" bson:"config_exchange_timestamp"`
	GigabytePrice           int64           `json:"gigabyte_price" bson:"gigabyte_price"`
	InfoFetchError          string          `json:"info_fetch_error" bson:"info_fetch_error"`
	InfoFetchTimestamp      time.Time       `json:"info_fetch_timestamp" bson:"info_fetch_timestamp"`
	LocationFetchError      string          `json:"location_fetch_error" bson:"location_fetch_error"`
	LocationFetchTimestamp  time.Time       `json:"location_fetch_timestamp" bson:"location_fetch_timestamp"`
	RemoteURL               string          `json:"remote_url" bson:"remote_url"`
	ServerConfig            []byte          `json:"server_config" bson:"server_config"`
	SessionID               uint64          `json:"session_id" bson:"session_id"`
	Status                  hubtypes.Status `json:"status" bson:"status"`
	SubscriptionID          uint64          `json:"subscription_id" bson:"subscription_id"`
	Type                    types.NodeType  `json:"type" bson:"type"`
}
