package main

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/sentinel-official/health-check/types"
)

type Config struct {
	App              *types.AppConfig
	Chain            *types.ChainConfig
	Database         *types.DatabaseConfig
	Key              *types.KeyConfig
	Query            *types.QueryConfig
	Tx               *types.TxConfig
	MaxGigabytePrice int64         `json:"max_gigabyte_price,omitempty"`
	PaymentDenom     string        `json:"payment_denom,omitempty"`
	RequestTimeout   time.Duration `json:"request_timeout,omitempty"`
}

func AddConfigFlagsToCmd(cmd *cobra.Command) {
	types.AddAppConfigFlagsToCmd(appName, cmd)
	types.AddChainConfigFlagsToCmd(cmd)
	types.AddDatabaseConfigFlagsToCmd(cmd)
	types.AddKeyConfigFlagsToCmd(cmd)
	types.AddQueryConfigFlagsToCmd(cmd)
	types.AddTxConfigFlagsToCmd(cmd)

	cmd.Flags().Int64("max_gigabyte_price", 25_000_000, "Max gigabyte price in udvpn denomination")
	cmd.Flags().Duration("request_timeout", 15*time.Second, "HTTP request timeout")
	cmd.Flags().String("payment_denom", "udvpn", "Payment denomination")
}

func NewConfigFromFlags(flags *pflag.FlagSet) (*Config, error) {
	appConfig, err := types.NewAppConfigFromFlags(flags)
	if err != nil {
		return nil, err
	}

	chainConfig, err := types.NewChainConfigFromFlags(flags)
	if err != nil {
		return nil, err
	}

	databaseConfig, err := types.NewDatabaseConfigFromFlags(flags)
	if err != nil {
		return nil, err
	}

	keyConfig, err := types.NewKeyConfigFromFlags(flags)
	if err != nil {
		return nil, err
	}

	queryConfig, err := types.NewQueryConfigFromFlags(flags)
	if err != nil {
		return nil, err
	}

	txConfig, err := types.NewTxConfigFromFlags(flags)
	if err != nil {
		return nil, err
	}

	maxGigabytePrice, err := flags.GetInt64("max_gigabyte_price")
	if err != nil {
		return nil, err
	}

	paymentDenom, err := flags.GetString("payment_denom")
	if err != nil {
		return nil, err
	}

	requestTimeout, err := flags.GetDuration("request_timeout")
	if err != nil {
		return nil, err
	}

	return &Config{
		App:              appConfig,
		Chain:            chainConfig,
		Database:         databaseConfig,
		Key:              keyConfig,
		Query:            queryConfig,
		Tx:               txConfig,
		MaxGigabytePrice: maxGigabytePrice,
		PaymentDenom:     paymentDenom,
		RequestTimeout:   requestTimeout,
	}, nil
}
