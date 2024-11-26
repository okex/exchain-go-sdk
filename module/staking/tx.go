package staking

import (
	"fmt"

	"github.com/okex/exchain/libs/cosmos-sdk/crypto/keys"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain-go-sdk/types/params"
	"github.com/okex/exchain-go-sdk/utils"
	"github.com/okex/exchain/x/common"
	stakingtypes "github.com/okex/exchain/x/staking/types"
)

// Deposit deposits an amount of okt to delegator account
func (sc stakingClient) Deposit(fromInfo keys.Info, passWd, coinsStr, memo string, accNum, seqNum uint64) (
	resp sdk.TxResponse, err error) {
	if err = params.CheckKeyParams(fromInfo, passWd); err != nil {
		return
	}

	coin, err := sdk.ParseDecCoin(coinsStr)
	if err != nil {
		return resp, fmt.Errorf("failed : parse Coins [%s] error: %s", coinsStr, err)
	}

	msg := stakingtypes.NewMsgDeposit(fromInfo.GetAddress(), coin)
	return sc.BuildAndBroadcastWithNonce(fromInfo.GetName(), passWd, memo, []sdk.Msg{msg}, accNum, seqNum)
}

// Withdraw withdraws an amount of okt and the corresponding shares from all validators
func (sc stakingClient) Withdraw(fromInfo keys.Info, passWd, coinsStr, memo string, accNum, seqNum uint64) (
	resp sdk.TxResponse, err error) {
	if err = params.CheckKeyParams(fromInfo, passWd); err != nil {
		return
	}

	coin, err := sdk.ParseDecCoin(coinsStr)
	if err != nil {
		return resp, fmt.Errorf("failed : parse Coins [%s] error: %s", coinsStr, err)
	}

	msg := stakingtypes.NewMsgWithdraw(fromInfo.GetAddress(), coin)
	return sc.BuildAndBroadcastWithNonce(fromInfo.GetName(), passWd, memo, []sdk.Msg{msg}, accNum, seqNum)
}

// Vote votes to the some specific validators
func (sc stakingClient) AddShares(fromInfo keys.Info, passWd string, valAddrsStr []string, memo string, accNum, seqNum uint64) (
	resp sdk.TxResponse, err error) {
	if err = params.CheckAddSharesParams(fromInfo, passWd, valAddrsStr); err != nil {
		return
	}

	valAddrs, err := utils.ParseValAddresses(valAddrsStr)
	if err != nil {
		return resp, fmt.Errorf("failed. validator address parsed error: %s", err.Error())
	}

	msg := stakingtypes.NewMsgAddShares(fromInfo.GetAddress(), valAddrs)
	return sc.BuildAndBroadcastWithNonce(fromInfo.GetName(), passWd, memo, []sdk.Msg{msg}, accNum, seqNum)
}

// DestroyValidator deregisters the validator and unbond the min-self-delegation
func (sc stakingClient) DestroyValidator(fromInfo keys.Info, passWd string, memo string, accNum, seqNum uint64) (
	resp sdk.TxResponse, err error) {
	if err = params.CheckKeyParams(fromInfo, passWd); err != nil {
		return
	}

	msg := stakingtypes.NewMsgDestroyValidator(fromInfo.GetAddress())
	return sc.BuildAndBroadcastWithNonce(fromInfo.GetName(), passWd, memo, []sdk.Msg{msg}, accNum, seqNum)
}

// CreateValidator creates a new validator
func (sc stakingClient) CreateValidator(fromInfo keys.Info, passWd, pubkeyStr, moniker, identity, website, details,
	memo string, accNum, seqNum uint64) (resp sdk.TxResponse, err error) {
	if err = params.CheckKeyParams(fromInfo, passWd); err != nil {
		return
	}

	pubkey, err := stakingtypes.GetConsPubKeyBech32(pubkeyStr)
	if err != nil {
		return
	}

	description := stakingtypes.NewDescription(moniker, identity, website, details)
	minSelfDelegation := sdk.NewDecCoinFromDec(common.NativeToken, stakingtypes.DefaultMinSelfDelegation)
	msg := stakingtypes.NewMsgCreateValidator(sdk.ValAddress(fromInfo.GetAddress()), pubkey, description, minSelfDelegation)
	return sc.BuildAndBroadcastWithNonce(fromInfo.GetName(), passWd, memo, []sdk.Msg{msg}, accNum, seqNum)
}

// EditValidator edits the description on a validator by the owner
func (sc stakingClient) EditValidator(fromInfo keys.Info, passWd, moniker, identity, website, details, memo string, accNum,
	seqNum uint64) (resp sdk.TxResponse, err error) {
	if err = params.CheckKeyParams(fromInfo, passWd); err != nil {
		return
	}

	description := stakingtypes.NewDescription(moniker, identity, website, details)
	msg := stakingtypes.NewMsgEditValidator(sdk.ValAddress(fromInfo.GetAddress()), description)
	return sc.BuildAndBroadcastWithNonce(fromInfo.GetName(), passWd, memo, []sdk.Msg{msg}, accNum, seqNum)
}

// RegisterProxy registers the identity of proxy
func (sc stakingClient) RegisterProxy(fromInfo keys.Info, passWd, memo string, accNum, seqNum uint64) (
	resp sdk.TxResponse, err error) {
	if err = params.CheckKeyParams(fromInfo, passWd); err != nil {
		return
	}

	msg := stakingtypes.NewMsgRegProxy(fromInfo.GetAddress(), true)
	return sc.BuildAndBroadcastWithNonce(fromInfo.GetName(), passWd, memo, []sdk.Msg{msg}, accNum, seqNum)
}

// UnregisterProxy registers the identity of proxy
func (sc stakingClient) UnregisterProxy(fromInfo keys.Info, passWd, memo string, accNum, seqNum uint64) (
	resp sdk.TxResponse, err error) {
	if err = params.CheckKeyParams(fromInfo, passWd); err != nil {
		return
	}

	msg := stakingtypes.NewMsgRegProxy(fromInfo.GetAddress(), false)
	return sc.BuildAndBroadcastWithNonce(fromInfo.GetName(), passWd, memo, []sdk.Msg{msg}, accNum, seqNum)
}

// BindProxy binds the staking tokens to a proxy
func (sc stakingClient) BindProxy(fromInfo keys.Info, passWd, proxyAddrStr, memo string, accNum, seqNum uint64) (
	resp sdk.TxResponse, err error) {
	if err = params.CheckSendParams(fromInfo, passWd, proxyAddrStr); err != nil {
		return
	}

	proxyAddr, err := sdk.AccAddressFromBech32(proxyAddrStr)
	if err != nil {
		return resp, fmt.Errorf("failed. parse Address [%s] error: %s", proxyAddrStr, err)
	}

	msg := stakingtypes.NewMsgBindProxy(fromInfo.GetAddress(), proxyAddr)
	return sc.BuildAndBroadcastWithNonce(fromInfo.GetName(), passWd, memo, []sdk.Msg{msg}, accNum, seqNum)
}

// UnbindProxy unbinds the staking tokens from a proxy
func (sc stakingClient) UnbindProxy(fromInfo keys.Info, passWd, memo string, accNum, seqNum uint64) (
	resp sdk.TxResponse, err error) {
	if err = params.CheckKeyParams(fromInfo, passWd); err != nil {
		return
	}

	msg := stakingtypes.NewMsgUnbindProxy(fromInfo.GetAddress())
	return sc.BuildAndBroadcastWithNonce(fromInfo.GetName(), passWd, memo, []sdk.Msg{msg}, accNum, seqNum)
}
