package context

import (
	"log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	subscriptiontypes "github.com/sentinel-official/hub/x/subscription/types"
)

func (c *Context) QuerySubscription(id uint64) (result subscriptiontypes.Subscription, err error) {
	log.Println("QuerySubscription", id)

	qsc := subscriptiontypes.NewQueryServiceClient(c.ctx)
	resp, err := qsc.QuerySubscription(
		c,
		&subscriptiontypes.QuerySubscriptionRequest{
			Id: id,
		},
	)
	if err != nil {
		return nil, err
	}
	if err := c.ctx.InterfaceRegistry.UnpackAny(resp.Subscription, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Context) QuerySubscriptionsForAccount(accAddr sdk.AccAddress, pagination *query.PageRequest) (result subscriptiontypes.Subscriptions, err error) {
	log.Println("QuerySubscriptionsForAccount", accAddr.String(), pagination.String())

	qsc := subscriptiontypes.NewQueryServiceClient(c.ctx)
	resp, err := qsc.QuerySubscriptionsForAccount(
		c,
		&subscriptiontypes.QuerySubscriptionsForAccountRequest{
			Address:    accAddr.String(),
			Pagination: pagination,
		},
	)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(resp.Subscriptions); i++ {
		var item subscriptiontypes.Subscription
		if err = c.ctx.InterfaceRegistry.UnpackAny(resp.Subscriptions[i], &item); err != nil {
			return nil, err
		}

		result = append(result, item)
	}

	return result, nil
}
