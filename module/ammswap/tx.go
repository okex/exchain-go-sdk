package ammswap

import (
	"time"

	"github.com/okex/exchain/libs/cosmos-sdk/crypto/keys"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain-go-sdk/types/params"
	ammswaptypes "github.com/okex/exchain/x/ammswap/types"
)

// AddLiquidity adds the number of liquidity of a token pair
func (ac ammswapClient) AddLiquidity(fromInfo keys.Info, passWd, minLiquidity, maxBaseAmount, quoteAmount, deadlineDuration,
	memo string, accNum, seqNum uint64) (resp sdk.TxResponse, err error) {
	if err = params.CheckKeyParams(fromInfo, passWd); err != nil {
		return
	}

	minLiquidityDec, err := sdk.NewDecFromStr(minLiquidity)
	if err != nil {
		return
	}

	maxBaseAmountDecCoin, err := sdk.ParseDecCoin(maxBaseAmount)
	if err != nil {
		return
	}

	quoteAmountDecCoin, err := sdk.ParseDecCoin(quoteAmount)
	if err != nil {
		return
	}

	duration, err := time.ParseDuration(deadlineDuration)
	if err != nil {
		return
	}
	deadline := time.Now().Add(duration).Unix()

	msg := ammswaptypes.NewMsgAddLiquidity(minLiquidityDec, maxBaseAmountDecCoin, quoteAmountDecCoin, deadline, fromInfo.GetAddress())
	return ac.BuildAndBroadcastWithNonce(fromInfo.GetName(), passWd, memo, []sdk.Msg{msg}, accNum, seqNum)
}

// RemoveLiquidity removes the number of liquidity of a token pair
func (ac ammswapClient) RemoveLiquidity(fromInfo keys.Info, passWd, liquidity, minBaseAmount, minQuoteAmount, deadlineDuration,
	memo string, accNum, seqNum uint64) (resp sdk.TxResponse, err error) {
	if err = params.CheckKeyParams(fromInfo, passWd); err != nil {
		return
	}

	liquidityDec, err := sdk.NewDecFromStr(liquidity)
	if err != nil {
		return
	}

	minBaseAmountDecCoin, err := sdk.ParseDecCoin(minBaseAmount)
	if err != nil {
		return
	}

	minQuoteAmountDecCoin, err := sdk.ParseDecCoin(minQuoteAmount)
	if err != nil {
		return
	}

	duration, err := time.ParseDuration(deadlineDuration)
	if err != nil {
		return
	}
	deadline := time.Now().Add(duration).Unix()

	msg := ammswaptypes.NewMsgRemoveLiquidity(liquidityDec, minBaseAmountDecCoin, minQuoteAmountDecCoin, deadline, fromInfo.GetAddress())
	return ac.BuildAndBroadcastWithNonce(fromInfo.GetName(), passWd, memo, []sdk.Msg{msg}, accNum, seqNum)
}

// CreateExchange creates a token pair in swap module
func (ac ammswapClient) CreateExchange(fromInfo keys.Info, passWd, baseToken, quoteToken, memo string, accNum, seqNum uint64) (
	resp sdk.TxResponse, err error) {
	if err = params.CheckKeyParams(fromInfo, passWd); err != nil {
		return
	}

	msg := ammswaptypes.NewMsgCreateExchange(baseToken, quoteToken, fromInfo.GetAddress())
	return ac.BuildAndBroadcastWithNonce(fromInfo.GetName(), passWd, memo, []sdk.Msg{msg}, accNum, seqNum)
}

// TokenSwap swaps the number of specific token with another type token
func (ac ammswapClient) TokenSwap(fromInfo keys.Info, passWd, soldTokenAmount, minBoughtTokenAmount, recipient, deadlineDuration,
	memo string, accNum, seqNum uint64) (resp sdk.TxResponse, err error) {
	if err = params.CheckKeyParams(fromInfo, passWd); err != nil {
		return
	}

	soldTokenDecCoin, err := sdk.ParseDecCoin(soldTokenAmount)
	if err != nil {
		return
	}

	minBoughtTokenDecCoin, err := sdk.ParseDecCoin(minBoughtTokenAmount)
	if err != nil {
		return
	}

	duration, err := time.ParseDuration(deadlineDuration)
	if err != nil {
		return
	}
	deadline := time.Now().Add(duration).Unix()

	var recip sdk.AccAddress
	if recipient == "" {
		recip = fromInfo.GetAddress()
	} else {
		if recip, err = sdk.AccAddressFromBech32(recipient); err != nil {
			return
		}
	}

	msg := ammswaptypes.NewMsgTokenToToken(soldTokenDecCoin, minBoughtTokenDecCoin, deadline, recip, fromInfo.GetAddress())
	return ac.BuildAndBroadcastWithNonce(fromInfo.GetName(), passWd, memo, []sdk.Msg{msg}, accNum, seqNum)
}
