package gosdk

import (
	"github.com/okex/okchain-go-sdk/module/auth"
	"github.com/okex/okchain-go-sdk/module/backend"
	"github.com/okex/okchain-go-sdk/module/dex"
	"github.com/okex/okchain-go-sdk/module/governance"
	"github.com/okex/okchain-go-sdk/module/order"
	"github.com/okex/okchain-go-sdk/module/staking"
	"github.com/okex/okchain-go-sdk/module/tendermint"
	"github.com/okex/okchain-go-sdk/module/token"
	sdk "github.com/okex/okchain-go-sdk/types"
)

// const
const (
	BroadcastSync  = sdk.BroadcastSync
	BroadcastAsync = sdk.BroadcastAsync
	BroadcastBlock = sdk.BroadcastBlock

	// vote for the proposal
	VoteYes        = "yes"
	VoteAbstain    = "abstain"
	VoteNo         = "no"
	VoteNoWithVeto = "no_with_veto"
)

var (
	// NewClientConfig gives an easy way for the callers to set client config
	NewClientConfig = sdk.NewClientConfig
)

// nolint
type (
	TxResponse = sdk.TxResponse
	// auth
	Account = auth.Account
	// staking
	Validator     = staking.Validator
	DelegatorResp = staking.DelegatorResp
	// token
	Token             = token.Token
	AccountTokensInfo = token.AccountTokensInfo
	// dex
	TokenPair = dex.TokenPair
	// order
	BookRes     = order.BookRes
	OrderDetail = order.OrderDetail
	// backend
	Ticker      = backend.Ticker
	MatchResult = backend.MatchResult
	Order       = backend.Order
	Deal        = backend.Deal
	// tendermint
	Block            = tendermint.Block
	BlockResults     = tendermint.BlockResults
	ResultCommit     = tendermint.ResultCommit
	ResultValidators = tendermint.ResultValidators
	ResultTx         = tendermint.ResultTx
	ResultTxs        = tendermint.ResultTxs
	// governance
	Proposal = governance.Proposal
)
