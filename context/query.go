package context

import (
	"encoding/hex"
	"log"
	"strings"
	"time"

	coretypes "github.com/tendermint/tendermint/rpc/core/types"
)

func (c *Context) QueryTx(hash string) (result *coretypes.ResultTx, err error) {
	log.Println("QueryTx", hash)

	buf, err := hex.DecodeString(hash)
	if err != nil {
		return nil, err
	}

	result, err = c.ctx.Client.Tx(c, buf, false)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return nil, nil
		}

		return nil, err
	}

	return result, nil
}

func (c *Context) QueryTxWithRetry(hash string) (result *coretypes.ResultTx, err error) {
	log.Println("QueryTxWithRetry", hash)

	for tries := c.queryMaxTries; tries > 0; tries-- {
		result, err = c.QueryTx(hash)
		if err != nil {
			return nil, err
		}
		if result != nil {
			break
		}

		time.Sleep(1 * time.Second)
	}

	return result, nil
}
