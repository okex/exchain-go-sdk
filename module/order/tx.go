package order

import (
	"errors"
	"strings"

	"github.com/okex/exchain/libs/cosmos-sdk/crypto/keys"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain-go-sdk/types/params"
	"github.com/okex/exchain-go-sdk/utils"
	ordertypes "github.com/okex/exchain/x/order/types"
)

// NewOrders places orders with some detail info
func (oc orderClient) NewOrders(fromInfo keys.Info, passWd, products, sides, prices, quantities, memo string, accNum,
	seqNum uint64) (resp sdk.TxResponse, err error) {
	if len(products) == 0 || len(sides) == 0 || len(prices) == 0 || len(quantities) == 0 {
		return resp, errors.New("failed. empty param input")
	}

	productStrs := strings.Split(products, ",")
	sideStrs := strings.Split(sides, ",")
	priceStrs := strings.Split(prices, ",")
	quantityStrs := strings.Split(quantities, ",")
	if err = params.CheckNewOrderParams(fromInfo, passWd, productStrs, sideStrs, priceStrs, quantityStrs); err != nil {
		return
	}

	orderItems, err := utils.BuildOrderItems(productStrs, sideStrs, priceStrs, quantityStrs)
	if err != nil {
		return
	}

	msg := ordertypes.NewMsgNewOrders(fromInfo.GetAddress(), orderItems)
	return oc.BuildAndBroadcastWithNonce(fromInfo.GetName(), passWd, memo, []sdk.Msg{msg}, accNum, seqNum)
}

// CancelOrders cancels orders by orderIDs
func (oc orderClient) CancelOrders(fromInfo keys.Info, passWd, orderIDs, memo string, accNum, seqNum uint64) (
	resp sdk.TxResponse, err error) {
	if len(orderIDs) == 0 {
		return resp, errors.New("failed. empty orderIDs input")
	}

	orderIDStrs := strings.Split(orderIDs, ",")
	if err = params.CheckCancelOrderParams(fromInfo, passWd, orderIDStrs); err != nil {
		return
	}

	msg := ordertypes.NewMsgCancelOrders(fromInfo.GetAddress(), orderIDStrs)
	return oc.BuildAndBroadcastWithNonce(fromInfo.GetName(), passWd, memo, []sdk.Msg{msg}, accNum, seqNum)
}
