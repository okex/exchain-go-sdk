package types

import (
	"errors"
	"math/big"

	apptypes "github.com/okex/exchain/app/types"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	authtypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	tmbytes "github.com/okex/exchain/libs/tendermint/libs/bytes"
	rpcclient "github.com/okex/exchain/libs/tendermint/rpc/client"
)

// BaseClient shows the expected behavior for a base client
type BaseClient interface {
	ClientQuery
	ClientTx
	TxHandler
	SimulationHandler
	GetCodec() *codec.Codec
	GetConfig() ClientConfig
}

// TxHandler shows the expected behavior to handle tx
type TxHandler interface {
	BuildAndBroadcastWithNonce(fromName, passphrase, memo string, msgs []sdk.Msg, accNumber, seqNumber uint64) (sdk.TxResponse, error)
	BuildAndBroadcast(fromName, passphrase, memo string, msgs []sdk.Msg, accNumber, seqNumber uint64) (sdk.TxResponse, error)
	BuildStdTx(fromName, passphrase, memo string, msgs []sdk.Msg, accNumber, seqNumber uint64) (*authtypes.StdTx, error)
	BuildUnsignedStdTxOffline(msgs []sdk.Msg, memo string) *authtypes.StdTx
}

// SimulationHandler shows the expected behavior to handle simulation
type SimulationHandler interface {
	CalculateGas(txBytes []byte) (authtypes.StdFee, error)
	BuildTxForSim(msgs []sdk.Msg, memo string, accNumber, seqNumber uint64) ([]byte, error)
}

// ClientQuery shows the expected query behavior
type ClientQuery interface {
	rpcclient.SignClient
	rpcclient.HistoryClient
	rpcclient.StatusClient
	Query(path string, key tmbytes.HexBytes) ([]byte, int64, error)
	QueryStore(key tmbytes.HexBytes, storeName, endPath string) ([]byte, int64, error)
}

// ClientTx shows the expected tx behavior
type ClientTx interface {
	Broadcast(txBytes []byte, broadcastMode string) (res sdk.TxResponse, err error)
}

// ClientConfig records the base config of gosdk client
type ClientConfig struct {
	NodeURI       string
	ChainID       string
	ChainIDBigInt *big.Int
	BroadcastMode string
	Gas           uint64
	GasAdjustment float64
	Fees          sdk.DecCoins
	GasPrices     sdk.DecCoins
}

// NewClientConfig creates a new instance of ClientConfig
func NewClientConfig(nodeURI, chainID string, broadcastMode string, feesStr string, gas uint64, gasAdjustment float64,
	gasPricesStr string) (cliConfig ClientConfig, err error) {
	var fees, gasPrices sdk.DecCoins
	if len(feesStr) != 0 {
		fees, err = sdk.ParseDecCoins(feesStr)
		if err != nil {
			return
		}
	}

	if len(gasPricesStr) != 0 {
		if gasAdjustment <= 1 {
			return cliConfig, errors.New("failed. gasAdjustment must be greater than 1 with the auto gas calculating")
		}

		gasPrices, err = sdk.ParseDecCoins(gasPricesStr)
		if err != nil {
			return
		}
	}

	chainIDBigInt, err := apptypes.ParseChainID(chainID)
	if err != nil {
		return
	}

	return ClientConfig{
		NodeURI:       nodeURI,
		ChainID:       chainID,
		ChainIDBigInt: chainIDBigInt,
		BroadcastMode: broadcastMode,
		Gas:           gas,
		GasAdjustment: gasAdjustment,
		Fees:          fees,
		GasPrices:     gasPrices,
	}, err
}
