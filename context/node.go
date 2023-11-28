package context

import (
	"log"

	"github.com/cosmos/cosmos-sdk/types/query"
	hubtypes "github.com/sentinel-official/hub/types"
	nodetypes "github.com/sentinel-official/hub/x/node/types"
)

func (c *Context) QueryNodes(status hubtypes.Status, pagination *query.PageRequest) (result nodetypes.Nodes, err error) {
	log.Println("QueryNodes", status.String(), pagination.String())

	qsc := nodetypes.NewQueryServiceClient(c.ctx)
	resp, err := qsc.QueryNodes(
		c,
		nodetypes.NewQueryNodesRequest(
			status,
			pagination,
		),
	)
	if err != nil {
		return nil, err
	}

	return resp.Nodes, nil
}
