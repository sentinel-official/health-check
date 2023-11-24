package types

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdkstd "github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authvesting "github.com/cosmos/cosmos-sdk/x/auth/vesting"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/sentinel-official/hub/x/vpn"
)

type EncodingConfig struct {
	Amino             *codec.LegacyAmino
	Codec             codec.Codec
	InterfaceRegistry codectypes.InterfaceRegistry
	TxConfig          client.TxConfig
}

func NewEncodingConfig() *EncodingConfig {
	var (
		amino             = codec.NewLegacyAmino()
		interfaceRegistry = codectypes.NewInterfaceRegistry()
		cdc               = codec.NewProtoCodec(interfaceRegistry)
		txConfig          = authtx.NewTxConfig(cdc, authtx.DefaultSignModes)
	)

	return &EncodingConfig{
		Amino:             amino,
		Codec:             cdc,
		InterfaceRegistry: interfaceRegistry,
		TxConfig:          txConfig,
	}
}

func DefaultEncodingConfig() *EncodingConfig {
	var (
		config  = NewEncodingConfig()
		modules = module.NewBasicManager(
			auth.AppModuleBasic{},
			authvesting.AppModuleBasic{},
			authzmodule.AppModuleBasic{},
			bank.AppModuleBasic{},
			feegrantmodule.AppModuleBasic{},
			vpn.AppModuleBasic{},
		)
	)

	sdkstd.RegisterLegacyAminoCodec(config.Amino)
	sdkstd.RegisterInterfaces(config.InterfaceRegistry)
	modules.RegisterLegacyAminoCodec(config.Amino)
	modules.RegisterInterfaces(config.InterfaceRegistry)

	return config
}
