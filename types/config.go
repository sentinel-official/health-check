package types

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type ChainConfig struct {
	ID      string `json:"id,omitempty"`
	RPCAddr string `json:"rpc_addr,omitempty"`
}

func NewChainConfigFromFlags(flags *pflag.FlagSet) (*ChainConfig, error) {
	id, err := flags.GetString("chain.id")
	if err != nil {
		return nil, err
	}

	rpcAddr, err := flags.GetString("chain.rpc_addr")
	if err != nil {
		return nil, err
	}

	return &ChainConfig{
		ID:      id,
		RPCAddr: rpcAddr,
	}, nil
}

func AddChainConfigFlagsToCmd(cmd *cobra.Command) {
	cmd.Flags().String("chain.id", "sentinelhub-2", "Chain ID")
	cmd.Flags().String("chain.rpc_addr", "https://rpc.sentinel.co:443", "Chain RPC address")
}

type DatabaseConfig struct {
	Name string `json:"name,omitempty"`
	URI  string `json:"uri,omitempty"`
}

func NewDatabaseConfigFromFlags(flags *pflag.FlagSet) (*DatabaseConfig, error) {
	name, err := flags.GetString("database.name")
	if err != nil {
		return nil, err
	}

	uri, err := flags.GetString("database.uri")
	if err != nil {
		return nil, err
	}

	return &DatabaseConfig{
		Name: name,
		URI:  uri,
	}, nil
}

func AddDatabaseConfigFlagsToCmd(cmd *cobra.Command) {
	cmd.Flags().String("database.name", "health-check", "Database name")
	cmd.Flags().String("database.uri", "mongodb://127.0.0.1:27017", "Database URI")
}

type KeyConfig struct {
	Account         uint32 `json:"account,omitempty"`
	BIP39Passphrase string `json:"bip39_passphrase,omitempty"`
	Index           uint32 `json:"index,omitempty"`
	Name            string `json:"name,omitempty"`
	Type            uint32 `json:"type,omitempty"`
}

func NewKeyConfigFromFlags(flags *pflag.FlagSet) (*KeyConfig, error) {
	account, err := flags.GetUint32("key.account")
	if err != nil {
		return nil, err
	}

	bip39Passphrase, err := flags.GetString("key.bip39_passphrase")
	if err != nil {
		return nil, err
	}

	index, err := flags.GetUint32("key.index")
	if err != nil {
		return nil, err
	}

	name, err := flags.GetString("key.name")
	if err != nil {
		return nil, err
	}

	keyType, err := flags.GetUint32("key.type")
	if err != nil {
		return nil, err
	}

	return &KeyConfig{
		Account:         account,
		BIP39Passphrase: bip39Passphrase,
		Index:           index,
		Name:            name,
		Type:            keyType,
	}, nil
}

func AddKeyConfigFlagsToCmd(cmd *cobra.Command) {
	cmd.Flags().Uint32("key.account", 0, "Account number")
	cmd.Flags().String("key.bip39_passphrase", "", "BIP39 passphrase")
	cmd.Flags().Uint32("key.index", 0, "Key index")
	cmd.Flags().String("key.name", "key-1", "Key name")
	cmd.Flags().Uint32("key.type", 118, "Key type")
}

type QueryConfig struct {
	MaxTries int64 `json:"max_tries,omitempty"`
}

func NewQueryConfigFromFlags(flags *pflag.FlagSet) (*QueryConfig, error) {
	maxTries, err := flags.GetInt64("query.max_tries")
	if err != nil {
		return nil, err
	}

	return &QueryConfig{
		MaxTries: maxTries,
	}, nil
}

func AddQueryConfigFlagsToCmd(cmd *cobra.Command) {
	cmd.Flags().Int64("query.max_tries", 60, "Maximum number of query tries")
}

type ServerConfig struct {
	ListenPort uint16 `json:"listen_port,omitempty"`
}

func NewServerConfigFromFlags(flags *pflag.FlagSet) (*ServerConfig, error) {
	listenPort, err := flags.GetUint16("server.listen_port")
	if err != nil {
		return nil, err
	}

	return &ServerConfig{
		ListenPort: listenPort,
	}, nil
}

func AddServerConfigFlagsToCmd(cmd *cobra.Command) {
	cmd.Flags().Uint16("server.listen_port", 8080, "Server listen port")
}

type TxConfig struct {
	BroadcastMode      string  `json:"broadcast_mode,omitempty"`
	FeeGranterAddr     string  `json:"fee_granter_addr,omitempty"`
	Fees               string  `json:"fees,omitempty"`
	GasAdjustment      float64 `json:"gas_adjustment,omitempty"`
	GasPrices          string  `json:"gas_prices,omitempty"`
	Gas                uint64  `json:"gas,omitempty"`
	Memo               string  `json:"memo,omitempty"`
	SignMode           string  `json:"sign_mode,omitempty"`
	SimulateAndExecute bool    `json:"simulate_and_execute,omitempty"`
	TimeoutHeight      uint64  `json:"timeout_height,omitempty"`
}

func NewTxConfigFromFlags(flags *pflag.FlagSet) (*TxConfig, error) {
	broadcastMode, err := flags.GetString("tx.broadcast_mode")
	if err != nil {
		return nil, err
	}

	feeGranterAddr, err := flags.GetString("tx.fee_granter_addr")
	if err != nil {
		return nil, err
	}

	fees, err := flags.GetString("tx.fees")
	if err != nil {
		return nil, err
	}

	gasAdjustment, err := flags.GetFloat64("tx.gas_adjustment")
	if err != nil {
		return nil, err
	}

	gasPrices, err := flags.GetString("tx.gas_prices")
	if err != nil {
		return nil, err
	}

	gas, err := flags.GetUint64("tx.gas")
	if err != nil {
		return nil, err
	}

	memo, err := flags.GetString("tx.memo")
	if err != nil {
		return nil, err
	}

	signMode, err := flags.GetString("tx.sign_mode")
	if err != nil {
		return nil, err
	}

	simulateAndExecute, err := flags.GetBool("tx.simulate_and_execute")
	if err != nil {
		return nil, err
	}

	timeoutHeight, err := flags.GetUint64("tx.timeout_height")
	if err != nil {
		return nil, err
	}

	return &TxConfig{
		BroadcastMode:      broadcastMode,
		FeeGranterAddr:     feeGranterAddr,
		Fees:               fees,
		GasAdjustment:      gasAdjustment,
		GasPrices:          gasPrices,
		Gas:                gas,
		Memo:               memo,
		SignMode:           signMode,
		SimulateAndExecute: simulateAndExecute,
		TimeoutHeight:      timeoutHeight,
	}, nil
}

func AddTxConfigFlagsToCmd(cmd *cobra.Command) {
	cmd.Flags().String("tx.broadcast_mode", "sync", "Transaction broadcast mode")
	cmd.Flags().String("tx.fee_granter_addr", "", "Fee granter address")
	cmd.Flags().String("tx.fees", "", "Transaction fees")
	cmd.Flags().Float64("tx.gas_adjustment", 1.2, "Gas adjustment")
	cmd.Flags().String("tx.gas_prices", "0.1udvpn", "Gas prices")
	cmd.Flags().Uint64("tx.gas", 0, "Gas limit")
	cmd.Flags().String("tx.memo", "", "Transaction memo")
	cmd.Flags().String("tx.sign_mode", "", "Transaction sign mode")
	cmd.Flags().Bool("tx.simulate_and_execute", true, "Simulate and execute transaction")
	cmd.Flags().Uint64("tx.timeout_height", 0, "Transaction timeout height")
}

type AppConfig struct {
	Name string `json:"name,omitempty"`
}

func NewAppConfigFromFlags(flags *pflag.FlagSet) (*AppConfig, error) {
	name, err := flags.GetString("app.name")
	if err != nil {
		return nil, err
	}

	return &AppConfig{
		Name: name,
	}, nil
}

func AddAppConfigFlagsToCmd(name string, cmd *cobra.Command) {
	cmd.Flags().String("app.name", name, "Application name")
}

type Config struct {
}
