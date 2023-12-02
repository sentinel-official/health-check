package record

import (
	hubtypes "github.com/sentinel-official/hub/types"

	"github.com/sentinel-official/health-check/models"
)

type Record struct {
	Addr                string          `json:"addr"`
	ConfigExchangeError string          `json:"config_exchange_error"`
	InfoFetchError      string          `json:"info_fetch_error"`
	LocationFetchError  string          `json:"location_fetch_error"`
	Status              hubtypes.Status `json:"status"`
}

func NewRecord(v *models.Record) *Record {
	return &Record{
		Addr:                v.Addr,
		ConfigExchangeError: v.ConfigExchangeError,
		InfoFetchError:      v.InfoFetchError,
		LocationFetchError:  v.LocationFetchError,
		Status:              v.Status,
	}
}

type Records []*Record

func NewRecords(v []*models.Record) []*Record {
	var items = make([]*Record, 0, len(v))
	for i := 0; i < len(v); i++ {
		items = append(items, NewRecord(v[i]))
	}

	return items
}
