package context

import (
	"context"
	"io"
	"reflect"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sentinel-official/health-check/types"
)

type Context struct {
	context.Context
	appName       string
	ctx           client.Context
	database      *mongo.Database
	queryMaxTries int64
	txf           tx.Factory
}

func NewContext(ctx context.Context) *Context {
	return &Context{
		Context: ctx,
	}
}

func NewDefaultContext() *Context {
	encCfg := types.DefaultEncodingConfig()
	c := NewContext(context.TODO()).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithTxConfig(encCfg.TxConfig)

	c.ctx = c.ctx.WithCodec(encCfg.Codec).
		WithInterfaceRegistry(encCfg.InterfaceRegistry).
		WithLegacyAmino(encCfg.Amino).
		WithOutputFormat("text").
		WithOutput(io.Discard).
		WithSkipConfirmation(true)

	return c
}

func (c *Context) WithKey(mnemonic string, keyType, account, index uint32, bip39Passphrase string) (keyring.Info, error) {
	return c.ctx.Keyring.NewAccount(
		c.ctx.FromName, mnemonic, bip39Passphrase,
		hd.CreateHDPath(keyType, account, index).String(), hd.Secp256k1,
	)
}

func (c *Context) Database() *mongo.Database {
	return c.database
}

func (c *Context) FromAddr() (sdk.AccAddress, error) {
	key, err := c.ctx.Keyring.Key(c.ctx.FromName)
	if err != nil {
		return nil, err
	}

	return key.GetAddress(), nil
}

func (c *Context) PrepareDatabase(username, password, uri, name string) (*mongo.Database, error) {
	dbClient, err := c.PrepareDatabaseClient(username, password, uri)
	if err != nil {
		return nil, err
	}

	return dbClient.Database(name), nil
}

func (c *Context) PrepareDatabaseClient(username, password, uri string) (*mongo.Client, error) {
	registry := bson.NewRegistry()
	registry.RegisterTypeMapEntry(bson.TypeDateTime, reflect.TypeOf(time.Time{}))
	registry.RegisterTypeMapEntry(bson.TypeEmbeddedDocument, reflect.TypeOf(bson.M{}))

	opts := options.Client().
		SetAppName(c.appName).
		ApplyURI(uri).
		SetRegistry(registry).
		SetMaxPoolSize(0)

	if username != "" && password != "" {
		opts = opts.SetAuth(
			options.Credential{
				Username: username,
				Password: password,
			},
		)
	}

	return mongo.Connect(c, opts)
}

func (c *Context) Sign(buf []byte) ([]byte, cryptotypes.PubKey, error) {
	return c.ctx.Keyring.Sign(c.ctx.FromName, buf)
}

func (c *Context) WithAccountRetriever(v client.AccountRetriever) *Context {
	c.ctx = c.ctx.WithAccountRetriever(v)
	c.txf = c.txf.WithAccountRetriever(v)

	return c
}

func (c *Context) WithAppName(v string) *Context {
	c.appName = v

	return c
}

func (c *Context) WithBroadcastMode(v string) *Context {
	c.ctx = c.ctx.WithBroadcastMode(v)

	return c
}

func (c *Context) WithChainID(v string) *Context {
	c.ctx = c.ctx.WithChainID(v)
	c.txf = c.txf.WithChainID(v)

	return c
}

func (c *Context) WithDatabase(v *mongo.Database) *Context {
	c.database = v

	return c
}

func (c *Context) WithFeeGranterAddr(v string) (*Context, error) {
	if v == "" {
		return c, nil
	}

	addr, err := sdk.AccAddressFromBech32(v)
	if err != nil {
		return nil, err
	}

	c.ctx = c.ctx.WithFeeGranterAddress(addr)
	return c, nil
}

func (c *Context) WithFees(v string) *Context {
	c.txf = c.txf.WithFees(v)

	return c
}

func (c *Context) WithFromName(v string) *Context {
	c.ctx = c.ctx.WithFrom(v).
		WithFromName(v)

	return c
}

func (c *Context) WithGas(v uint64) *Context {
	c.txf = c.txf.WithGas(v)

	return c
}

func (c *Context) WithGasAdjustment(v float64) *Context {
	c.txf = c.txf.WithGasAdjustment(v)

	return c
}

func (c *Context) WithGasPrices(v string) *Context {
	c.txf = c.txf.WithGasPrices(v)

	return c
}

func (c *Context) WithInMemoryKeyring(opts ...keyring.Option) *Context {
	return c.WithKeyring(keyring.NewInMemory(opts...))
}

func (c *Context) WithKeyring(v keyring.Keyring) *Context {
	c.ctx = c.ctx.WithKeyring(v)
	c.txf = c.txf.WithKeybase(v)

	return c
}

func (c *Context) WithMemo(v string) *Context {
	c.txf = c.txf.WithMemo(v)

	return c
}

func (c *Context) WithQueryMaxTries(v int64) *Context {
	c.queryMaxTries = v

	return c
}

func (c *Context) WithRPCAddr(v string) (*Context, error) {
	rpcClient, err := rpchttp.New(v, "/websocket")
	if err != nil {
		return nil, err
	}

	c.ctx = c.ctx.WithClient(rpcClient).
		WithNodeURI(v)

	return c, nil
}

func (c *Context) WithSignMode(v string) *Context {
	mode := signing.SignMode_SIGN_MODE_UNSPECIFIED
	switch v {
	case flags.SignModeDirect:
		mode = signing.SignMode_SIGN_MODE_DIRECT
	case flags.SignModeLegacyAminoJSON:
		mode = signing.SignMode_SIGN_MODE_LEGACY_AMINO_JSON
	}

	c.ctx = c.ctx.WithSignModeStr(v)
	c.txf = c.txf.WithSignMode(mode)

	return c
}

func (c *Context) WithSimulateAndExecute(v bool) *Context {
	c.txf = c.txf.WithSimulateAndExecute(v)

	return c
}

func (c *Context) WithTimeoutHeight(v uint64) *Context {
	c.txf = c.txf.WithTimeoutHeight(v)

	return c
}

func (c *Context) WithTxConfig(v client.TxConfig) *Context {
	c.ctx = c.ctx.WithTxConfig(v)
	c.txf = c.txf.WithTxConfig(v)

	return c
}
