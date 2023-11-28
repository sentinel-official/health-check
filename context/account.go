package context

import (
	"log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func (c *Context) QueryAccount(accAddr sdk.AccAddress) (result authtypes.AccountI, err error) {
	log.Println("QueryAccount", accAddr.String())

	qc := authtypes.NewQueryClient(c.ctx)
	resp, err := qc.Account(
		c,
		&authtypes.QueryAccountRequest{
			Address: accAddr.String(),
		},
	)
	if err != nil {
		return nil, err
	}
	if err := c.ctx.InterfaceRegistry.UnpackAny(resp.Account, &result); err != nil {
		return nil, err
	}

	return result, nil
}
