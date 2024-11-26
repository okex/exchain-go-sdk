package wasm

import (
	"fmt"
	"github.com/okex/exchain/libs/cosmos-sdk/crypto/keys"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/okex/exchain/x/wasm/ioutils"
	"github.com/okex/exchain/x/wasm/types"
	"io/ioutil"
	"strconv"
	"strings"
)

func (c wasmClient) StoreCode(fromInfo keys.Info, passWd string, accNum, seqNum uint64, memo string, wasmFilePath string, onlyAddr string, everybody, nobody bool) (int, error) {
	msg, err := parseStoreCodeMsg(wasmFilePath, fromInfo.GetAddress(), onlyAddr, everybody, nobody)
	if err != nil {
		return 0, errors.Wrapf(err, "parse StoreCodeMsg failed")
	}

	if err = msg.ValidateBasic(); err != nil {
		return 0, err
	}

	res, err := c.BuildAndBroadcastWithNonce(fromInfo.GetName(), passWd, memo, []sdk.Msg{msg}, accNum, seqNum)
	if err != nil {
		return 0, err
	}

	if res.Code != 0 {
		return 0, errors.New("wasm", res.Code, "send tx failed"+res.String())
	}

	return parseCodeID(res.RawLog), nil
}

func (c wasmClient) InstantiateContract(fromInfo keys.Info, passWd string, accNum, seqNum uint64, memo string, codeID uint64, initMsg string, amount string, label string, adminAddr string, noAdmin bool) (string, error) {
	msg, err := parseInstantiateMsg(codeID, initMsg, fromInfo.GetAddress(), amount, label, adminAddr, noAdmin)
	if err != nil {
		return "", errors.Wrapf(err, "parse InstantiateContractMsg failed")
	}

	if err = msg.ValidateBasic(); err != nil {
		return "", err
	}

	res, err := c.BuildAndBroadcastWithNonce(fromInfo.GetName(), passWd, memo, []sdk.Msg{msg}, accNum, seqNum)
	if err != nil {
		return "", err
	}

	if res.Code != 0 {
		return "", errors.New("wasm", res.Code, "send tx failed"+res.String())
	}
	return parseContractAddress(res.RawLog), nil
}

func (c wasmClient) ExecuteContract(fromInfo keys.Info, passWd string, accNum, seqNum uint64, memo string, contractAddr string, execMsg string, amount string) (*sdk.TxResponse, error) {
	msg, err := parseExecuteMsg(contractAddr, execMsg, fromInfo.GetAddress(), amount)
	if err != nil {
		return nil, errors.Wrapf(err, "parse ExecuteContractMsg failed")
	}

	if err = msg.ValidateBasic(); err != nil {
		return nil, err
	}

	res, err := c.BuildAndBroadcastWithNonce(fromInfo.GetName(), passWd, memo, []sdk.Msg{msg}, accNum, seqNum)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (c wasmClient) MigrateContract(fromInfo keys.Info, passWd string, accNum, seqNum uint64, memo string, codeID uint64, contractAddr string, migrateMsg string) (*sdk.TxResponse, error) {
	msg, err := parseMigrateContractMsg(codeID, contractAddr, fromInfo.GetAddress(), migrateMsg)
	if err != nil {
		return nil, errors.Wrapf(err, "parse MigrateContractMsg failed")
	}

	if err = msg.ValidateBasic(); err != nil {
		return nil, err
	}

	res, err := c.BuildAndBroadcastWithNonce(fromInfo.GetName(), passWd, memo, []sdk.Msg{msg}, accNum, seqNum)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (c wasmClient) UpdateContractAdmin(fromInfo keys.Info, passWd string, accNum, seqNum uint64, memo string, contractAddr string, adminAddr string) (*sdk.TxResponse, error) {
	_, err := sdk.AccAddressFromBech32(contractAddr)
	if err != nil {
		return nil, err
	}

	msg := types.MsgUpdateAdmin{
		Sender:   fromInfo.GetAddress().String(),
		Contract: contractAddr,
		NewAdmin: adminAddr,
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

func (c wasmClient) ClearContractAdmin(fromInfo keys.Info, passWd string, accNum, seqNum uint64, memo string, contractAddr string) (*sdk.TxResponse, error) {
	_, err := sdk.AccAddressFromBech32(contractAddr)
	if err != nil {
		return nil, err
	}

	msg := types.MsgClearAdmin{
		Sender:   fromInfo.GetAddress().String(),
		Contract: contractAddr,
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

func parseStoreCodeMsg(wasmFilePath string, sender sdk.AccAddress, onlyAddrStr string, everybody, nobody bool) (types.MsgStoreCode, error) {
	wasm, err := ioutil.ReadFile(wasmFilePath)
	if err != nil {
		return types.MsgStoreCode{}, err
	}

	// gzip the wasm file
	if ioutils.IsWasm(wasm) {
		wasm, err = ioutils.GzipIt(wasm)

		if err != nil {
			return types.MsgStoreCode{}, err
		}
	} else if !ioutils.IsGzip(wasm) {
		return types.MsgStoreCode{}, fmt.Errorf("invalid input file. Use wasm binary or gzip")
	}

	var perm *types.AccessConfig

	if onlyAddrStr != "" {
		allowedAddr, err := sdk.AccAddressFromBech32(onlyAddrStr)
		if err != nil {
			return types.MsgStoreCode{}, err
		}
		x := types.AccessTypeOnlyAddress.With(allowedAddr)
		perm = &x
	} else if everybody {
		perm = &types.AllowEverybody
	} else if nobody {
		perm = &types.AllowNobody
	}

	msg := types.MsgStoreCode{
		Sender:                sender.String(),
		WASMByteCode:          wasm,
		InstantiatePermission: perm,
	}
	return msg, nil
}

func parseInstantiateMsg(codeID uint64, initMsg string, sender sdk.AccAddress, amountStr, label, adminStr string, noAdmin bool) (types.MsgInstantiateContract, error) {
	amount, err := sdk.ParseCoinsNormalized(amountStr)
	if err != nil {
		return types.MsgInstantiateContract{}, fmt.Errorf("amount: %s", err)
	}

	if label == "" {
		return types.MsgInstantiateContract{}, fmt.Errorf("label is required on all contracts")
	}

	// ensure sensible admin is set (or explicitly immutable)
	if adminStr == "" && !noAdmin {
		return types.MsgInstantiateContract{}, fmt.Errorf("you must set an admin or explicitly pass no-admin to make it immutible (wasmd issue #719)")
	}
	if adminStr != "" && noAdmin {
		return types.MsgInstantiateContract{}, fmt.Errorf("you set an admin and passed no-admin, those cannot both be true")
	}

	// build and sign the transaction, then broadcast to Tendermint
	msg := types.MsgInstantiateContract{
		Sender: sender.String(),
		CodeID: codeID,
		Label:  label,
		Funds:  sdk.CoinsToCoinAdapters(amount),
		Msg:    []byte(initMsg),
		Admin:  adminStr,
	}
	return msg, nil
}

func parseExecuteMsg(contractAddr string, execMsg string, sender sdk.AccAddress, amountStr string) (types.MsgExecuteContract, error) {
	amount, err := sdk.ParseCoinsNormalized(amountStr)
	if err != nil {
		return types.MsgExecuteContract{}, err
	}

	return types.MsgExecuteContract{
		Sender:   sender.String(),
		Contract: contractAddr,
		Funds:    sdk.CoinsToCoinAdapters(amount),
		Msg:      []byte(execMsg),
	}, nil
}

func parseMigrateContractMsg(codeID uint64, contractAddr string, sender sdk.AccAddress, migrateMsg string) (types.MsgMigrateContract, error) {
	_, err := sdk.AccAddressFromBech32(contractAddr)
	if err != nil {
		return types.MsgMigrateContract{}, err
	}

	msg := types.MsgMigrateContract{
		Sender:   sender.String(),
		Contract: contractAddr,
		CodeID:   codeID,
		Msg:      []byte(migrateMsg),
	}
	return msg, nil
}

func parseCodeID(str string) int {
	index := strings.LastIndex(str, ":")
	codeIDStr := str[index:]
	codeIDStr = codeIDStr[2 : strings.Index(codeIDStr, "}")-1]
	codeID, _ := strconv.Atoi(codeIDStr)

	return codeID
}

func parseContractAddress(str string) string {
	index := strings.Index(str, "address")
	contractAddr := str[index+18 : index+18+42]
	return contractAddr
}
