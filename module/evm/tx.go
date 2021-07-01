package evm

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/okex/exchain-go-sdk/module/evm/types"
	"github.com/okex/exchain-go-sdk/utils"
	apptypes "github.com/okex/exchain/app/types"
	evmtypes "github.com/okex/exchain/x/evm/types"
)

// SendTx sends tx to transfer or call contract function
func (ec evmClient) SendTx(fromInfo keys.Info, passWd, toAddrStr, amountStr, payloadStr, memo string, accNum, seqNum uint64) (
	resp sdk.TxResponse, err error) {
	toAddr, err := utils.ToCosmosAddress(toAddrStr)
	if err != nil {
		return
	}

	amount := sdk.ZeroDec()
	if len(amountStr) != 0 {
		amount, err = sdk.NewDecFromStr(amountStr)
		if err != nil {
			return
		}
	}

	var data []byte
	if len(payloadStr) != 0 {
		if !strings.HasPrefix(payloadStr, "0x") {
			payloadStr = fmt.Sprintf("0x%s", payloadStr)
		}

		data, err = hexutil.Decode(payloadStr)
		if err != nil {
			return
		}
	}

	msg := evmtypes.NewMsgEthermint(
		seqNum,
		&toAddr,
		sdk.NewIntFromBigInt(amount.Int),
		ec.bc.GetConfig().Gas,
		sdk.NewInt(apptypes.DefaultGasPrice),
		data,
		fromInfo.GetAddress(),
	)
	return ec.bc.BuildAndBroadcast(fromInfo.GetName(), passWd, memo, []sdk.Msg{msg}, accNum, seqNum)
}

// CreateContract generates a transaction to deploy a smart contract
func (ec evmClient) CreateContract(fromInfo keys.Info, passWd, amountStr, payloadStr, memo string, accNum, seqNum uint64) (
	resp sdk.TxResponse, contractAddrStr string, err error) {
	if !strings.HasPrefix(payloadStr, "0x") {
		payloadStr = fmt.Sprintf("0x%s", payloadStr)
	}

	data, err := hexutil.Decode(payloadStr)
	if err != nil {
		return
	}

	amount := sdk.ZeroDec()
	if len(amountStr) != 0 {
		amount, err = sdk.NewDecFromStr(amountStr)
		if err != nil {
			return
		}
	}

	msg := evmtypes.NewMsgEthermint(
		seqNum,
		nil,
		sdk.NewIntFromBigInt(amount.Int),
		ec.bc.GetConfig().Gas,
		sdk.NewInt(apptypes.DefaultGasPrice),
		data,
		fromInfo.GetAddress(),
	)
	if resp, err = ec.bc.BuildAndBroadcast(fromInfo.GetName(), passWd, memo, []sdk.Msg{msg}, accNum, seqNum); err != nil {
		return
	}

	contractAddrStr = ethcrypto.CreateAddress(common.BytesToAddress(fromInfo.GetAddress().Bytes()), seqNum).Hex()
	return
}

// SendTxEthereum sends an ethereum tx
func (ec evmClient) SendTxEthereum(privHex, toAddrStr, amountStr, payloadStr string, gasLimit, seqNum uint64) (
	resp sdk.TxResponse, err error) {
	priv, err := ethcrypto.HexToECDSA(privHex)
	if err != nil {
		return
	}

	toAddr, err := utils.ToHexAddress(toAddrStr)
	if err != nil {
		return
	}

	amount := sdk.ZeroDec()
	if len(amountStr) != 0 {
		amount, err = sdk.NewDecFromStr(amountStr)
		if err != nil {
			return
		}
	}

	var data []byte
	if len(payloadStr) != 0 {
		if !strings.HasPrefix(payloadStr, "0x") {
			payloadStr = fmt.Sprintf("0x%s", payloadStr)
		}

		data, err = hexutil.Decode(payloadStr)
		if err != nil {
			return
		}
	}

	ethMsg := evmtypes.NewMsgEthereumTx(
		seqNum,
		&toAddr,
		amount.Int,
		gasLimit,
		types.DefaultGasPrice,
		data,
	)

	config := ec.bc.GetConfig()
	if err = ethMsg.Sign(config.ChainIDBigInt, priv); err != nil {
		return
	}

	bytes, err := ec.bc.GetCodec().MarshalBinaryLengthPrefixed(ethMsg)
	if err != nil {
		return resp, fmt.Errorf("failed. encoded MsgEthereumTx error: %s", err)
	}

	return ec.bc.Broadcast(bytes, ec.bc.GetConfig().BroadcastMode)
}

// CreateContractEthereum generates an ethereum tx to deploy a smart contract
func (ec evmClient) CreateContractEthereum(privHex, amountStr, payloadStr string, gasLimit, seqNum uint64) (
	resp sdk.TxResponse, err error) {
	priv, err := ethcrypto.HexToECDSA(privHex)
	if err != nil {
		return
	}

	amount := sdk.ZeroDec()
	if len(amountStr) != 0 {
		amount, err = sdk.NewDecFromStr(amountStr)
		if err != nil {
			return
		}
	}

	var data []byte
	if len(payloadStr) != 0 {
		if !strings.HasPrefix(payloadStr, "0x") {
			payloadStr = fmt.Sprintf("0x%s", payloadStr)
		}

		data, err = hexutil.Decode(payloadStr)
		if err != nil {
			return
		}
	}

	ethMsg := evmtypes.NewMsgEthereumTxContract(
		seqNum,
		amount.Int,
		gasLimit,
		types.DefaultGasPrice,
		data,
	)

	config := ec.bc.GetConfig()
	if err = ethMsg.Sign(config.ChainIDBigInt, priv); err != nil {
		return
	}

	bytes, err := ec.bc.GetCodec().MarshalBinaryLengthPrefixed(ethMsg)
	if err != nil {
		return resp, fmt.Errorf("failed. encoded MsgEthereumTx error: %s", err)
	}

	return ec.bc.Broadcast(bytes, ec.bc.GetConfig().BroadcastMode)
}
