package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/sentinel-official/health-check/types"
)

type Config struct {
	App        *types.AppConfig
	Database   *types.DatabaseConfig
	ListenAddr string `json:"listen_addr,omitempty"`
}

func AddConfigFlagsToCmd(cmd *cobra.Command) {
	types.AddAppConfigFlagsToCmd(appName, cmd)
	types.AddDatabaseConfigFlagsToCmd(cmd)

	cmd.Flags().String("listen-addr", ":8080", "Listen address")
}

func NewConfigFromFlags(flags *pflag.FlagSet) (*Config, error) {
	appConfig, err := types.NewAppConfigFromFlags(flags)
	if err != nil {
		return nil, err
	}

	databaseConfig, err := types.NewDatabaseConfigFromFlags(flags)
	if err != nil {
		return nil, err
	}

	listenAddr, err := flags.GetString("listen-addr")
	if err != nil {
		return nil, err
	}

	return &Config{
		App:        appConfig,
		Database:   databaseConfig,
		ListenAddr: listenAddr,
	}, nil
}
