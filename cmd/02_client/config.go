package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/sentinel-official/health-check/types"
)

type Config struct {
	App      *types.AppConfig
	Database *types.DatabaseConfig
	Address  string `json:"address,omitempty"`
}

func AddConfigFlagsToCmd(cmd *cobra.Command) {
	types.AddAppConfigFlagsToCmd(appName, cmd)
	types.AddDatabaseConfigFlagsToCmd(cmd)

	cmd.Flags().String("address", "", "Node address")
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

	address, err := flags.GetString("address")
	if err != nil {
		return nil, err
	}

	return &Config{
		App:      appConfig,
		Database: databaseConfig,
		Address:  address,
	}, nil
}
