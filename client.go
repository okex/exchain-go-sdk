package gosdk

import (
	"fmt"
	"github.com/okex/exchain-go-sdk/module/feesplit"
	"github.com/okex/exchain-go-sdk/module/ibc"
	"github.com/okex/exchain-go-sdk/module/wasm"
	ibcTypes "github.com/okex/exchain/libs/ibc-go/modules/apps/transfer/types"
	feesplitTypes "github.com/okex/exchain/x/feesplit/types"

	"github.com/okex/exchain-go-sdk/exposed"
	"github.com/okex/exchain-go-sdk/module"
	"github.com/okex/exchain-go-sdk/module/ammswap"
	ammswaptypes "github.com/okex/exchain-go-sdk/module/ammswap/types"
	"github.com/okex/exchain-go-sdk/module/auth"
	authtypes "github.com/okex/exchain-go-sdk/module/auth/types"
	"github.com/okex/exchain-go-sdk/module/dex"
	dextypes "github.com/okex/exchain-go-sdk/module/dex/types"
	"github.com/okex/exchain-go-sdk/module/distribution"
	distrtypes "github.com/okex/exchain-go-sdk/module/distribution/types"
	"github.com/okex/exchain-go-sdk/module/evm"
	evmtypes "github.com/okex/exchain-go-sdk/module/evm/types"
	"github.com/okex/exchain-go-sdk/module/farm"
	"github.com/okex/exchain-go-sdk/module/governance"
	govtypes "github.com/okex/exchain-go-sdk/module/governance/types"
	"github.com/okex/exchain-go-sdk/module/order"
	ordertypes "github.com/okex/exchain-go-sdk/module/order/types"
	"github.com/okex/exchain-go-sdk/module/slashing"
	slashingtypes "github.com/okex/exchain-go-sdk/module/slashing/types"
	"github.com/okex/exchain-go-sdk/module/staking"
	stakingtypes "github.com/okex/exchain-go-sdk/module/staking/types"
	"github.com/okex/exchain-go-sdk/module/tendermint"
	tmtypes "github.com/okex/exchain-go-sdk/module/tendermint/types"
	"github.com/okex/exchain-go-sdk/module/token"
	tokentypes "github.com/okex/exchain-go-sdk/module/token/types"
	gosdktypes "github.com/okex/exchain-go-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	farmtypes "github.com/okex/exchain/x/farm/types"
	wasmTypes "github.com/okex/exchain/x/wasm/types"
)

// Client - structure of the main client of ExChain GoSDK
type Client struct {
	config  gosdktypes.ClientConfig
	cdc     *codec.Codec
	modules map[string]gosdktypes.Module
}

// NewClient creates a new instance of Client
func NewClient(config gosdktypes.ClientConfig) Client {
	cdc := gosdktypes.NewCodec()
	pClient := &Client{
		config:  config,
		cdc:     cdc,
		modules: make(map[string]gosdktypes.Module),
	}
	pBaseClient := module.NewBaseClient(cdc, &pClient.config)

	pClient.registerModule(
		ammswap.NewAmmSwapClient(pBaseClient),
		auth.NewAuthClient(pBaseClient),
		dex.NewDexClient(pBaseClient),
		distribution.NewDistrClient(pBaseClient),
		evm.NewEvmClient(pBaseClient),
		farm.NewFarmClient(pBaseClient),
		governance.NewGovClient(pBaseClient),
		order.NewOrderClient(pBaseClient),
		staking.NewStakingClient(pBaseClient),
		slashing.NewSlashingClient(pBaseClient),
		token.NewTokenClient(pBaseClient),
		tendermint.NewTendermintClient(pBaseClient),
		ibc.NewIbcClient(pBaseClient),
		wasm.NewWasmClient(pBaseClient),
		feesplit.NewfeesplitClient(pBaseClient),
	)

	return *pClient
}

func (cli *Client) registerModule(mods ...gosdktypes.Module) {
	for _, mod := range mods {
		moduleName := mod.Name()
		if _, ok := cli.modules[moduleName]; ok {
			panic(fmt.Sprintf("duplicated module: %s", moduleName))
		}
		// register codec by each module
		mod.RegisterCodec(cli.cdc)
		cli.modules[moduleName] = mod
	}
	gosdktypes.RegisterBasicCodec(cli.cdc)
	cli.cdc.Seal()
}

// GetConfig returns the client config
func (cli *Client) GetConfig() gosdktypes.ClientConfig {
	return cli.config
}

// nolint
func (cli *Client) AmmSwap() exposed.AmmSwap {
	return cli.modules[ammswaptypes.ModuleName].(exposed.AmmSwap)
}
func (cli *Client) Auth() exposed.Auth {
	return cli.modules[authtypes.ModuleName].(exposed.Auth)
}
func (cli *Client) Dex() exposed.Dex {
	return cli.modules[dextypes.ModuleName].(exposed.Dex)
}
func (cli *Client) Distribution() exposed.Distribution {
	return cli.modules[distrtypes.ModuleName].(exposed.Distribution)
}
func (cli *Client) Evm() exposed.Evm {
	return cli.modules[evmtypes.ModuleName].(exposed.Evm)
}
func (cli *Client) Farm() exposed.Farm {
	return cli.modules[farmtypes.ModuleName].(exposed.Farm)
}
func (cli *Client) Governance() exposed.Governance {
	return cli.modules[govtypes.ModuleName].(exposed.Governance)
}
func (cli *Client) Order() exposed.Order {
	return cli.modules[ordertypes.ModuleName].(exposed.Order)
}
func (cli *Client) Slashing() exposed.Slashing {
	return cli.modules[slashingtypes.ModuleName].(exposed.Slashing)
}
func (cli *Client) Staking() exposed.Staking {
	return cli.modules[stakingtypes.ModuleName].(exposed.Staking)
}
func (cli *Client) Tendermint() exposed.Tendermint {
	return cli.modules[tmtypes.ModuleName].(exposed.Tendermint)
}
func (cli *Client) Token() exposed.Token {
	return cli.modules[tokentypes.ModuleName].(exposed.Token)
}

func (cli *Client) Ibc() exposed.Ibc {
	return cli.modules[ibcTypes.ModuleName].(exposed.Ibc)
}

func (cli *Client) Wasm() exposed.Wasm {
	return cli.modules[wasmTypes.ModuleName].(exposed.Wasm)
}
func (cli *Client) Feesplit() exposed.Feesplit {
	return cli.modules[feesplitTypes.ModuleName].(exposed.Feesplit)
}
