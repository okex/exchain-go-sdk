package module

import (
	"errors"
	"fmt"

	"github.com/okex/exchain-go-sdk/types"
	"github.com/okex/exchain-go-sdk/types/tx"
	"github.com/okex/exchain-go-sdk/utils"

	extypes "github.com/okex/exchain/libs/tendermint/types"

	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	authtypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	tmbytes "github.com/okex/exchain/libs/tendermint/libs/bytes"
	rpcclient "github.com/okex/exchain/libs/tendermint/rpc/client"
	rpchttp "github.com/okex/exchain/libs/tendermint/rpc/client/http"
)

const (
	simulationPath = "/app/simulate"
)

type baseClient struct {
	rpcclient.Client
	config *types.ClientConfig
	cdc    *codec.Codec
}

// NewBaseClient creates a new instance of baseClient
func NewBaseClient(cdc *codec.Codec, pConfig *types.ClientConfig) *baseClient {
	rpc, err := rpchttp.New(pConfig.NodeURI, "/websocket")
	if err != nil {
		panic(fmt.Sprintf("failed to get client: %s", err))
	}
	return &baseClient{
		Client: rpc,
		config: pConfig,
		cdc:    cdc,
	}
}

// Query executes the basic query
func (bc *baseClient) Query(path string, key tmbytes.HexBytes) (res []byte, height int64, err error) {
	opts := rpcclient.ABCIQueryOptions{
		Height: 0,
		Prove:  false,
	}

	result, err := bc.ABCIQueryWithOptions(path, key, opts)
	if err != nil {
		return
	}

	resp := result.Response
	if !resp.IsOK() {
		return res, height, errors.New(resp.Log)
	}

	return resp.Value, resp.Height, err
}

// QueryStore executes the direct query to the store
func (bc *baseClient) QueryStore(key tmbytes.HexBytes, storeName, endPath string) ([]byte, int64, error) {
	path := fmt.Sprintf("/store/%s/%s", storeName, endPath)
	return bc.Query(path, key)
}

// Broadcast broadcasts by different modes
func (bc *baseClient) Broadcast(txBytes []byte, broadcastMode string) (res sdk.TxResponse, err error) {
	switch broadcastMode {
	case types.BroadcastSync:
		retBroadcastTx, err := bc.BroadcastTxSync(txBytes)
		return sdk.NewResponseFormatBroadcastTx(retBroadcastTx), err

	case types.BroadcastAsync:
		retBroadcastTx, err := bc.BroadcastTxAsync(txBytes)
		return sdk.NewResponseFormatBroadcastTx(retBroadcastTx), err

	case types.BroadcastBlock:
		retBroadcastTxCommit, err := bc.BroadcastTxCommit(txBytes)
		if err != nil {
			return sdk.NewResponseFormatBroadcastTxCommit(retBroadcastTxCommit), err
		}
		if !retBroadcastTxCommit.CheckTx.IsOK() {
			return sdk.NewResponseFormatBroadcastTxCommit(retBroadcastTxCommit), fmt.Errorf(retBroadcastTxCommit.CheckTx.Log)
		}
		if !retBroadcastTxCommit.DeliverTx.IsOK() {
			return sdk.NewResponseFormatBroadcastTxCommit(retBroadcastTxCommit), fmt.Errorf(retBroadcastTxCommit.DeliverTx.Log)
		}
		return sdk.NewResponseFormatBroadcastTxCommit(retBroadcastTxCommit), err

	default:
		err = fmt.Errorf("failed. unsupported broadcast mode %s; supported types: sync, async, block", broadcastMode)
	}
	return
}

// GetCodec gets the codec of the base client
func (bc *baseClient) GetCodec() *codec.Codec {
	return bc.cdc
}

// GetConfig gets the client config
func (bc *baseClient) GetConfig() types.ClientConfig {
	return *bc.config
}

// BuildAndBroadcast implements the TxHandler interface ; abandonded
func (bc *baseClient) BuildAndBroadcast(fromName, passphrase, memo string, msgs []sdk.Msg, accNumber,
	seqNumber uint64) (resp sdk.TxResponse, err error) {
	stdTx, err := bc.BuildStdTx(fromName, passphrase, memo, msgs, accNumber, seqNumber)
	if err != nil {
		return resp, fmt.Errorf("failed. build stdTx error: %s", err)
	}

	bytes, err := bc.cdc.MarshalBinaryLengthPrefixed(stdTx)
	if err != nil {
		return resp, fmt.Errorf("failed. encoded stdTx error: %s", err)
	}

	return bc.Broadcast(bytes, bc.GetConfig().BroadcastMode)
}

// BuildAndBroadcast implements the TxHandler interface
func (bc *baseClient) BuildAndBroadcastWithNonce(fromName, passphrase, memo string, msgs []sdk.Msg, accNumber,
	seqNumber uint64) (resp sdk.TxResponse, err error) {
	stdTx, err := bc.BuildStdTx(fromName, passphrase, memo, msgs, accNumber, seqNumber)
	if err != nil {
		return resp, fmt.Errorf("failed. build stdTx error: %s", err)
	}

	bytes, err := bc.cdc.MarshalBinaryLengthPrefixed(stdTx)
	if err != nil {
		return resp, fmt.Errorf("failed. encoded stdTx error: %s", err)
	}

	wrapedTx := &extypes.WrapCMTx{
		Tx:    bytes,
		Nonce: seqNumber,
	}
	txBytes, err := bc.cdc.MarshalJSON(wrapedTx)
	if err != nil {
		panic(fmt.Sprintln("MarshalJSON fail", err))
	}

	return bc.Broadcast(txBytes, bc.GetConfig().BroadcastMode)
	// return bc.Broadcast(bytes, bc.GetConfig().BroadcastMode)
}

// BuildAndSign builds std sign context and sign it
func (bc *baseClient) BuildStdTx(fromName, passphrase, memo string, msgs []sdk.Msg, accNumber, seqNumber uint64) (
	stdTx *authtypes.StdTx, err error) {
	config := bc.GetConfig()
	if len(config.ChainID) == 0 {
		return stdTx, errors.New("failed. empty chain ID")
	}

	var stdFee authtypes.StdFee
	if config.GasPrices.IsZero() {
		// fixed fees
		stdFee = authtypes.NewStdFee(config.Gas, config.Fees)
	} else {
		// auto gas calculation
		var txBytes []byte
		txBytes, err = bc.BuildTxForSim(msgs, memo, accNumber, seqNumber)
		if err != nil {
			return stdTx, fmt.Errorf("failed. build tx for simulation error: %s", err)
		}

		stdFee, err = bc.CalculateGas(txBytes)
		if err != nil {
			return
		}
	}

	signMsg := authtypes.StdSignMsg{
		ChainID:       config.ChainID,
		AccountNumber: accNumber,
		Sequence:      seqNumber,
		Memo:          memo,
		Msgs:          msgs,
		Fee:           stdFee,
	}

	sigBytes, err := tx.MakeSignature(fromName, passphrase, signMsg)
	if err != nil {
		return
	}

	return authtypes.NewStdTx(signMsg.Msgs, signMsg.Fee, []authtypes.StdSignature{sigBytes}, signMsg.Memo), err
}

// BuildUnsignedStdTxOffline builds a stdTx without signature
func (bc *baseClient) BuildUnsignedStdTxOffline(msgs []sdk.Msg, memo string) *authtypes.StdTx {
	config := bc.GetConfig()
	fee := authtypes.NewStdFee(config.Gas, bc.GetConfig().Fees)
	return authtypes.NewStdTx(msgs, fee, nil, memo)
}

// CalculateGas is designed for auto gas calculation and builds an available stdFee
func (bc *baseClient) CalculateGas(txBytes []byte) (stdFee authtypes.StdFee, err error) {
	config := bc.GetConfig()
	// estimate the gas by a simulation query
	rawRes, _, err := bc.Query(simulationPath, txBytes)
	if err != nil {
		return stdFee, utils.ErrClientQuery(err.Error())
	}

	// get simulation response
	var simRes sdk.SimulationResponse
	if err = bc.GetCodec().UnmarshalBinaryBare(rawRes, &simRes); err != nil {
		return
	}

	// enlarge the simulation result by gas adjustment in config
	adjustedGasLimt := uint64(config.GasAdjustment * float64(simRes.GasUsed))
	return calculateStdFee(config.GasPrices, adjustedGasLimt), err
}

// BuildTxForSim creates a StdSignMsg and encodes a transaction with the StdSignMsg for tx simulation
func (bc *baseClient) BuildTxForSim(msgs []sdk.Msg, memo string, accNumber, seqNumber uint64) ([]byte, error) {
	config := bc.GetConfig()

	// build std tx for simulation
	simStdTx := authtypes.NewStdTx(msgs, calculateStdFee(config.GasPrices, config.Gas), []authtypes.StdSignature{{}}, memo)
	return bc.GetCodec().MarshalBinaryLengthPrefixed(simStdTx)
}

func calculateStdFee(gasPrices sdk.DecCoins, gas uint64) authtypes.StdFee {
	gasLimitDec := sdk.NewDec(int64(gas))
	gasPricesLen := len(gasPrices)
	fees := make(sdk.Coins, gasPricesLen)
	for i := 0; i < gasPricesLen; i++ {
		// Derive the fees based on the provided gas prices, where fee = ceil(gasPrice * gasLimit)
		fees[i] = sdk.NewCoin(gasPrices[i].Denom, gasPrices[i].Amount.Mul(gasLimitDec).Ceil().RoundInt())
	}

	return authtypes.NewStdFee(gas, fees)
}
