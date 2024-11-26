package ibc

import (
	"fmt"
	"github.com/okex/exchain-go-sdk/module/auth"
	"github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	ibcmsg "github.com/okex/exchain/libs/cosmos-sdk/types/ibc-adapter"
	cmwraptx "github.com/okex/exchain/libs/tendermint/types"
	"math/big"
	"strings"

	"github.com/okex/exchain-go-sdk/module"
	gosdktypes "github.com/okex/exchain-go-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/client"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	codec_types "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	signing "github.com/okex/exchain/libs/cosmos-sdk/types/tx/signing"
	ibc_tx "github.com/okex/exchain/libs/cosmos-sdk/x/auth/ibc-tx"
	signing2 "github.com/okex/exchain/libs/cosmos-sdk/x/auth/ibcsigning"
	ibc_type "github.com/okex/exchain/libs/ibc-go/modules/apps/transfer/types"
	client_types "github.com/okex/exchain/libs/ibc-go/modules/core/02-client/types"
	tmcrypto "github.com/okex/exchain/libs/tendermint/crypto"
)

const (
	src_port = "transfer"
)

var (
	txConfig client.TxConfig
)

func init() {
	txConfig = newTxConfig()
}

func newTxConfig() client.TxConfig {
	interfaceRegistry := codec_types.NewInterfaceRegistry()
	marshaler := codec.NewProtoCodec(interfaceRegistry)
	return ibc_tx.NewTxConfig(marshaler, ibc_tx.DefaultSignModes)
}

func (ibc ibcClient) GetLatestHeight() uint64 {
	status, err := ibc.Status()
	if err != nil {

		fmt.Println(err)
		return 0
	}

	return uint64(status.SyncInfo.LatestBlockHeight)
}

func (ibc ibcClient) Transfer(priKey tmcrypto.PrivKey, srcChannel string, receiver string, amount string, fee sdk.CoinAdapters, memo string, timeoutHeight client_types.Height) (resp sdk.TxResponse, err error) {

	pubKey := ibc_tx.LagacyKey2PbKey(priKey.PubKey())

	// get account info
	accountInfo, err := auth.NewAuthClient(ibc.BaseClient).QueryAccount(pubKey.Address().String())
	if err != nil {

		return sdk.TxResponse{}, err
	}
	coins, err := sdk.ParseDecCoins(amount)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	coin := coins[0]

	if coin.Denom == sdk.DefaultIbcWei {
		coin.Amount = sdk.Dec{
			coin.Amount.Div(coin.Amount.BigInt(), big.NewInt(sdk.DefaultDecInt)),
		}
	}
	if coin.Denom == "okt" || coin.Denom == "OKT" {
		coin.Denom = "wei"
	}

	//
	if !strings.HasPrefix(coin.Denom, "ibc/") {
		denomTrace := ibc_type.ParseDenomTrace(coin.Denom)
		coin.Denom = denomTrace.IBCDenom()
	}

	// generate msg
	msg := &ibc_type.MsgTransfer{
		SourcePort:       src_port,
		SourceChannel:    srcChannel,
		Token:            sdk.NewCoinAdapter(coin.Denom, sdk.NewIntFromBigInt(coin.Amount.BigInt())),
		Sender:           sdk.AccAddress(pubKey.Address().Bytes()).String(),
		Receiver:         receiver,
		TimeoutHeight:    timeoutHeight,
		TimeoutTimestamp: 0,
	}

	// build unsignedTx
	txb, err := ibc.buildUnsignedTx(msg, memo, fee)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	// sign Tx
	tx, err := ibc.signedTx(txb, priKey, accountInfo.GetAccountNumber(), accountInfo.GetSequence())
	if err != nil {

		return sdk.TxResponse{}, err
	}

	// get tx bytes
	txBytes, err := txConfig.TxEncoder()(tx)
	if err != nil {

		return sdk.TxResponse{}, err
	}

	wrapedTx := &cmwraptx.WrapCMTx{
		Tx:    txBytes,
		Nonce: accountInfo.GetSequence(),
	}

	txBytes, err = ibc.GetCodec().MarshalJSON(wrapedTx)
	if err != nil {
		return sdk.TxResponse{}, errors.Wrap(err, "MarshalJSON fail")
	}

	return ibc.Broadcast(txBytes, ibc.GetConfig().BroadcastMode)
}

// get Client Height from destination chain
func getTimeoutHeight(dstRpc string) (client_types.Height, error) {
	cdc := gosdktypes.NewCodec()

	dstClient := module.NewBaseClient(cdc, &gosdktypes.ClientConfig{
		NodeURI: dstRpc,
	})

	status, err := dstClient.Status()
	if err != nil {

		return client_types.Height{}, err
	}

	latestHeight := status.SyncInfo.LatestBlockHeight

	return client_types.Height{
		RevisionNumber: client_types.ParseChainID(status.NodeInfo.Network),
		RevisionHeight: uint64(latestHeight + 1000),
	}, nil
}

func (ibc ibcClient) buildUnsignedTx(msg ibcmsg.Msg, memo string, fee sdk.CoinAdapters) (client.TxBuilder, error) {
	txb := txConfig.NewTxBuilder()

	// config txb
	txb.SetMemo(memo)
	txb.SetFeeAmount(fee)
	txb.SetGasLimit(ibc.GetConfig().Gas)
	txb.SetTimeoutHeight(0)

	// set ibc msgs
	if err := txb.SetMsgs(msg); err != nil {

		return nil, err
	}

	return txb, nil
}

func (ibc ibcClient) signedTx(txb client.TxBuilder, priKey tmcrypto.PrivKey, accNum uint64, seqNum uint64) (signing2.Tx, error) {
	signMode := txConfig.SignModeHandler().DefaultMode()

	pubKey := ibc_tx.LagacyKey2PbKey(priKey.PubKey())
	// init signature
	signature := signing.SignatureV2{
		PubKey: pubKey,
		Data: &signing.SingleSignatureData{
			SignMode: signMode,
		},
		Sequence: seqNum,
	}

	err := txb.SetSignatures(signature)
	if err != nil {

		return nil, err
	}

	signerData := signing2.SignerData{
		ChainID:       ibc.GetConfig().ChainID,
		AccountNumber: accNum,
		Sequence:      seqNum,
	}

	// bytes to sign
	signBytes, err := txConfig.SignModeHandler().GetSignBytes(signMode, signerData, txb.GetTx())
	if err != nil {
		panic(err)
	}

	// signed bytes
	sigBytes, err := priKey.Sign(signBytes)
	if err != nil {
		panic(err)
	}

	sigData := signing.SingleSignatureData{
		SignMode:  signMode,
		Signature: sigBytes,
	}

	sig := signing.SignatureV2{
		PubKey:   pubKey,
		Data:     &sigData,
		Sequence: seqNum,
	}

	// set signature
	err = txb.SetSignatures(sig)
	if err != nil {

		return nil, err
	}

	return txb.GetTx(), nil
}
