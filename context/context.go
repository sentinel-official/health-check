package context

import (
	"io"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/sentinel-official/health-check/types"
)

type Context struct {
	ctx client.Context
	txf tx.Factory
}

func NewContext() *Context {
	return &Context{}
}

func NewDefaultContext() *Context {
	var (
		cfg = types.DefaultEncodingConfig()
		ctx = client.Context{}.
			WithCodec(cfg.Codec).
			WithInterfaceRegistry(cfg.InterfaceRegistry).
			WithLegacyAmino(cfg.Amino).
			WithOutputFormat("text").
			WithOutput(io.Discard).
			WithSkipConfirmation(true)
	)

	return NewContext().
		WithClientContext(ctx).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithTxConfig(cfg.TxConfig)
}

func (c *Context) WithClientContext(v client.Context) *Context {
	c.ctx = v
	return c
}

func (c *Context) WithAccountRetriever(v client.AccountRetriever) *Context {
	c.ctx = c.ctx.WithAccountRetriever(v)
	c.txf = c.txf.WithAccountRetriever(v)
	return c
}

func (c *Context) WithTxConfig(v client.TxConfig) *Context {
	c.ctx = c.ctx.WithTxConfig(v)
	c.txf = c.txf.WithTxConfig(v)
	return c
}

func (c *Context) Codec() codec.Codec { return c.ctx.Codec }
