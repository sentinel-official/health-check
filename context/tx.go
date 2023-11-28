package context

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	abcitypes "github.com/tendermint/tendermint/abci/types"
)

func (c *Context) Tx(messages ...sdk.Msg) (resp *sdk.TxResponse, err error) {
	fromAddr, err := c.FromAddr()
	if err != nil {
		return nil, err
	}

	acc, err := c.QueryAccount(fromAddr)
	if err != nil {
		return nil, err
	}

	txf := c.txf.WithAccountNumber(acc.GetAccountNumber()).
		WithSequence(acc.GetSequence())

	if txf.SimulateAndExecute() {
		_, gas, err := tx.CalculateGas(c.ctx, txf, messages...)
		if err != nil {
			return nil, err
		}

		txf = txf.WithGas(gas)
	}

	txb, err := tx.BuildUnsignedTx(txf, messages...)
	if err != nil {
		return nil, err
	}

	txb.SetFeeGranter(c.ctx.GetFeeGranterAddress())
	if err := tx.Sign(txf, c.ctx.FromName, txb, true); err != nil {
		return nil, err
	}

	txBytes, err := c.ctx.TxConfig.TxEncoder()(txb.GetTx())
	if err != nil {
		return nil, err
	}

	resp, err = c.ctx.BroadcastTx(txBytes)
	if err != nil {
		return nil, err
	}

	switch resp.Code {
	case abcitypes.CodeTypeOK:
		return resp, nil
	case errors.ErrTxInMempoolCache.ABCICode():
		return resp, nil
	default:
		return nil, fmt.Errorf(resp.RawLog)
	}
}
