package context

import (
	"log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	sessiontypes "github.com/sentinel-official/hub/x/session/types"
)

func (c *Context) QuerySession(id uint64) (result *sessiontypes.Session, err error) {
	log.Println("QuerySession", id)

	qsc := sessiontypes.NewQueryServiceClient(c.ctx)
	resp, err := qsc.QuerySession(
		c,
		&sessiontypes.QuerySessionRequest{
			Id: id,
		},
	)
	if err != nil {
		return nil, err
	}

	return &resp.Session, nil
}

func (c *Context) QuerySessionsForAccount(accAddr sdk.AccAddress, pagination *query.PageRequest) (result sessiontypes.Sessions, err error) {
	log.Println("QuerySessionsForAccount", accAddr.String(), pagination.String())

	qsc := sessiontypes.NewQueryServiceClient(c.ctx)
	resp, err := qsc.QuerySessionsForAccount(
		c,
		&sessiontypes.QuerySessionsForAccountRequest{
			Address:    accAddr.String(),
			Pagination: pagination,
		},
	)
	if err != nil {
		return nil, err
	}

	return resp.Sessions, nil
}
