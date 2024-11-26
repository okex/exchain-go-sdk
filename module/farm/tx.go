package farm

import (
	"github.com/okex/exchain/libs/cosmos-sdk/crypto/keys"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain-go-sdk/types/params"
	farmtypes "github.com/okex/exchain/x/farm/types"
)

// CreatePool creates a farm pool
func (fc farmClient) CreatePool(fromInfo keys.Info, passWd, poolName, minLockAmountStr, yieldToken, memo string, accNum,
	seqNum uint64) (resp sdk.TxResponse, err error) {
	if err = params.CheckCreatePoolParams(fromInfo, passWd, poolName, minLockAmountStr, yieldToken); err != nil {
		return
	}

	minLockAmount, err := sdk.ParseDecCoin(minLockAmountStr)
	if err != nil {
		return
	}

	msg := farmtypes.NewMsgCreatePool(fromInfo.GetAddress(), poolName, minLockAmount, yieldToken)
	return fc.BuildAndBroadcastWithNonce(fromInfo.GetName(), passWd, memo, []sdk.Msg{msg}, accNum, seqNum)
}

// DestroyPool destroys a farm pool
func (fc farmClient) DestroyPool(fromInfo keys.Info, passWd, poolName, memo string, accNum, seqNum uint64) (
	resp sdk.TxResponse, err error) {
	if err = params.CheckPoolNameParams(fromInfo, passWd, poolName); err != nil {
		return
	}

	msg := farmtypes.NewMsgDestroyPool(fromInfo.GetAddress(), poolName)
	return fc.BuildAndBroadcastWithNonce(fromInfo.GetName(), passWd, memo, []sdk.Msg{msg}, accNum, seqNum)
}

// Provide provides a number of yield tokens into a pool
func (fc farmClient) Provide(fromInfo keys.Info, passWd, poolName, amountStr, yieldPerBlockStr string, startHeightToYield int64,
	memo string, accNum, seqNum uint64) (resp sdk.TxResponse, err error) {
	if err = params.CheckPoolNameParams(fromInfo, passWd, poolName); err != nil {
		return
	}

	amount, err := sdk.ParseDecCoin(amountStr)
	if err != nil {
		return
	}

	amountYieldPerBlock, err := sdk.NewDecFromStr(yieldPerBlockStr)
	if err != nil {
		return
	}

	msg := farmtypes.NewMsgProvide(poolName, fromInfo.GetAddress(), amount, amountYieldPerBlock, startHeightToYield)
	return fc.BuildAndBroadcastWithNonce(fromInfo.GetName(), passWd, memo, []sdk.Msg{msg}, accNum, seqNum)
}

// Lock locks a number of tokens for yield farming
func (fc farmClient) Lock(fromInfo keys.Info, passWd, poolName, amountStr, memo string, accNum, seqNum uint64) (
	resp sdk.TxResponse, err error) {
	if err = params.CheckPoolNameParams(fromInfo, passWd, poolName); err != nil {
		return
	}

	amount, err := sdk.ParseDecCoin(amountStr)
	if err != nil {
		return
	}

	msg := farmtypes.NewMsgLock(poolName, fromInfo.GetAddress(), amount)
	return fc.BuildAndBroadcastWithNonce(fromInfo.GetName(), passWd, memo, []sdk.Msg{msg}, accNum, seqNum)
}

// Unlock unlocks a number of tokens from the farm pool and claims the current yield
func (fc farmClient) Unlock(fromInfo keys.Info, passWd, poolName, amountStr, memo string, accNum, seqNum uint64) (
	resp sdk.TxResponse, err error) {
	if err = params.CheckPoolNameParams(fromInfo, passWd, poolName); err != nil {
		return
	}

	amount, err := sdk.ParseDecCoin(amountStr)
	if err != nil {
		return
	}

	msg := farmtypes.NewMsgUnlock(poolName, fromInfo.GetAddress(), amount)
	return fc.BuildAndBroadcastWithNonce(fromInfo.GetName(), passWd, memo, []sdk.Msg{msg}, accNum, seqNum)
}

// Claim claims yield farming rewards
func (fc farmClient) Claim(fromInfo keys.Info, passWd, poolName, memo string, accNum, seqNum uint64) (
	resp sdk.TxResponse, err error) {
	if err = params.CheckPoolNameParams(fromInfo, passWd, poolName); err != nil {
		return
	}

	msg := farmtypes.NewMsgClaim(poolName, fromInfo.GetAddress())
	return fc.BuildAndBroadcastWithNonce(fromInfo.GetName(), passWd, memo, []sdk.Msg{msg}, accNum, seqNum)
}
