package dex

import (
	"fmt"

	"github.com/okex/exchain/libs/cosmos-sdk/crypto/keys"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain-go-sdk/types/params"
	dextypes "github.com/okex/exchain/x/dex/types"
)

// List lists a trading pair on dex
func (dc dexClient) List(fromInfo keys.Info, passWd, baseAsset, quoteAsset, initPriceStr, memo string, accNum,
	seqNum uint64) (resp sdk.TxResponse, err error) {
	if err = params.CheckDexAssetsParams(fromInfo, passWd, baseAsset, quoteAsset); err != nil {
		return
	}

	initPrice := sdk.MustNewDecFromStr(initPriceStr)
	msg := dextypes.NewMsgList(fromInfo.GetAddress(), baseAsset, quoteAsset, initPrice)
	return dc.BuildAndBroadcastWithNonce(fromInfo.GetName(), passWd, memo, []sdk.Msg{msg}, accNum, seqNum)
}

// Deposit deposits some tokens to a specific product
func (dc dexClient) Deposit(fromInfo keys.Info, passWd, product, amountStr, memo string, accNum, seqNum uint64) (
	resp sdk.TxResponse, err error) {
	if err = params.CheckProductParams(fromInfo, passWd, product); err != nil {
		return
	}

	amount, err := sdk.ParseDecCoin(amountStr)
	if err != nil {
		return
	}

	msg := dextypes.NewMsgDeposit(product, amount, fromInfo.GetAddress())
	return dc.BuildAndBroadcastWithNonce(fromInfo.GetName(), passWd, memo, []sdk.Msg{msg}, accNum, seqNum)
}

// Withdraw withdraws some tokens from a specific product
func (dc dexClient) Withdraw(fromInfo keys.Info, passWd, product, amountStr, memo string, accNum, seqNum uint64) (
	resp sdk.TxResponse, err error) {
	if err = params.CheckProductParams(fromInfo, passWd, product); err != nil {
		return
	}

	amount, err := sdk.ParseDecCoin(amountStr)
	if err != nil {
		return
	}

	msg := dextypes.NewMsgWithdraw(product, amount, fromInfo.GetAddress())
	return dc.BuildAndBroadcastWithNonce(fromInfo.GetName(), passWd, memo, []sdk.Msg{msg}, accNum, seqNum)
}

// TransferOwnership changes the owner of a product
func (dc dexClient) TransferOwnership(fromInfo keys.Info, passWd, product, toAddrStr, memo string, accNum, seqNum uint64) (
	resp sdk.TxResponse, err error) {
	if err = params.CheckProductParams(fromInfo, passWd, product); err != nil {
		return
	}

	toAddr, err := sdk.AccAddressFromBech32(toAddrStr)
	if err != nil {
		return resp, fmt.Errorf("failed. parse Address [%s] error: %s", toAddrStr, err)
	}

	msg := dextypes.NewMsgTransferOwnership(fromInfo.GetAddress(), toAddr, product)
	return dc.BuildAndBroadcastWithNonce(fromInfo.GetName(), passWd, memo, []sdk.Msg{msg}, accNum, seqNum)
}

// ConfirmOwnership confirms the transfer-ownership of a product
func (dc dexClient) ConfirmOwnership(fromInfo keys.Info, passWd, product, memo string, accNum, seqNum uint64) (resp sdk.TxResponse, err error) {
	if err = params.CheckProductParams(fromInfo, passWd, product); err != nil {
		return
	}

	msg := dextypes.NewMsgConfirmOwnership(fromInfo.GetAddress(), product)
	return dc.BuildAndBroadcastWithNonce(fromInfo.GetName(), passWd, memo, []sdk.Msg{msg}, accNum, seqNum)
}

func (dc dexClient) RegisterDexOperator(fromInfo keys.Info, passWd, handleFeeAddrStr, website, memo string, accNum, seqNum uint64) (
	resp sdk.TxResponse, err error) {
	if err = params.CheckKeyParams(fromInfo, passWd); err != nil {
		return
	}

	handleFeeAddr, err := sdk.AccAddressFromBech32(handleFeeAddrStr)
	if err != nil {
		return
	}

	msg := dextypes.NewMsgCreateOperator(website, fromInfo.GetAddress(), handleFeeAddr)
	return dc.BuildAndBroadcastWithNonce(fromInfo.GetName(), passWd, memo, []sdk.Msg{msg}, accNum, seqNum)
}

func (dc dexClient) EditDexOperator(fromInfo keys.Info, passWd, handleFeeAddrStr, website, memo string, accNum, seqNum uint64) (
	resp sdk.TxResponse, err error) {
	if err = params.CheckKeyParams(fromInfo, passWd); err != nil {
		return
	}

	handleFeeAddr, err := sdk.AccAddressFromBech32(handleFeeAddrStr)
	if err != nil {
		return
	}

	msg := dextypes.NewMsgUpdateOperator(website, fromInfo.GetAddress(), handleFeeAddr)
	return dc.BuildAndBroadcastWithNonce(fromInfo.GetName(), passWd, memo, []sdk.Msg{msg}, accNum, seqNum)
}
