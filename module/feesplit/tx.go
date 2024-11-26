package feesplit

import (
	"fmt"
	"github.com/okex/exchain/libs/cosmos-sdk/crypto/keys"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/feesplit/types"
)

func (c feesplitClient) RegisterFeeSplit(fromInfo keys.Info, passWd string, accNum, seqNum uint64, memo string, contractAddress string, nonces []uint64, withdrawAddress string) (*sdk.TxResponse, error) {
	if err := types.ValidateNonZeroAddress(contractAddress); err != nil {
		return nil, fmt.Errorf("invalid contract hex address %w", err)
	}

	if len(nonces) == 0 {
		return nil, fmt.Errorf("invalid nonces")
	}

	if withdrawAddress == "" {
		withdrawAddress = fromInfo.GetAddress().String()
	}

	if _, err := sdk.AccAddressFromBech32(withdrawAddress); err != nil {
		return nil, fmt.Errorf("invalid withdraw bech32 address %w", err)
	}

	msg := &types.MsgRegisterFeeSplit{
		ContractAddress:   contractAddress,
		DeployerAddress:   fromInfo.GetAddress().String(),
		WithdrawerAddress: withdrawAddress,
		Nonces:            nonces,
	}

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	res, err := c.BuildAndBroadcastWithNonce(fromInfo.GetName(), passWd, memo, []sdk.Msg{msg}, accNum, seqNum)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (c feesplitClient) CancelFeeSplit(fromInfo keys.Info, passWd string, accNum, seqNum uint64, memo string, contractAddress string) (*sdk.TxResponse, error) {
	if err := types.ValidateNonZeroAddress(contractAddress); err != nil {
		return nil, fmt.Errorf("invalid contract hex address %w", err)
	}

	msg := &types.MsgCancelFeeSplit{
		ContractAddress: contractAddress,
		DeployerAddress: fromInfo.GetAddress().String(),
	}

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	res, err := c.BuildAndBroadcastWithNonce(fromInfo.GetName(), passWd, memo, []sdk.Msg{msg}, accNum, seqNum)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (c feesplitClient) UpdateFeeSplit(fromInfo keys.Info, passWd string, accNum, seqNum uint64, memo string, contractAddress string, withdrawAddress string) (*sdk.TxResponse, error) {
	if err := types.ValidateNonZeroAddress(contractAddress); err != nil {
		return nil, fmt.Errorf("invalid contract hex address %w", err)
	}

	if _, err := sdk.AccAddressFromBech32(withdrawAddress); err != nil {
		return nil, fmt.Errorf("invalid withdraw bech32 address %w", err)
	}

	msg := &types.MsgUpdateFeeSplit{
		ContractAddress:   contractAddress,
		DeployerAddress:   fromInfo.GetAddress().String(),
		WithdrawerAddress: withdrawAddress,
	}

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	res, err := c.BuildAndBroadcastWithNonce(fromInfo.GetName(), passWd, memo, []sdk.Msg{msg}, accNum, seqNum)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
