package governance

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain-go-sdk/module/governance/types"
	"github.com/okex/okexchain-go-sdk/types/params"
	dexutils "github.com/okex/okexchain/x/dex/client/utils"
	dextypes "github.com/okex/okexchain/x/dex/types"
	"github.com/okex/okexchain/x/gov"
)

// SubmitTextProposal submits the text proposal on OKChain
func (gc govClient) SubmitTextProposal(fromInfo keys.Info, passWd, proposalPath, memo string, accNum, seqNum uint64) (
	resp sdk.TxResponse, err error) {
	if err = params.CheckKeyParams(fromInfo, passWd); err != nil {
		return
	}

	proposal, err := parseProposalFromFile(proposalPath)
	if err != nil {
		return
	}

	deposit, err := sdk.ParseDecCoins(proposal.Deposit)
	if err != nil {
		return
	}

	msg := types.NewMsgSubmitProposal(
		types.NewTextProposal(proposal.Title, proposal.Description),
		deposit,
		fromInfo.GetAddress(),
	)

	return gc.BuildAndBroadcast(fromInfo.GetName(), passWd, memo, []sdk.Msg{msg}, accNum, seqNum)

}

// SubmitParamChangeProposal submits the proposal to change the params on OKChain
func (gc govClient) SubmitParamChangeProposal(fromInfo keys.Info, passWd, proposalPath, memo string, accNum, seqNum uint64) (
	resp sdk.TxResponse, err error) {
	if err = params.CheckKeyParams(fromInfo, passWd); err != nil {
		return
	}

	proposal, err := parseParamChangeProposalFromFile(proposalPath)
	if err != nil {
		return
	}

	msg := types.NewMsgSubmitProposal(
		types.NewParameterChangeProposal(
			proposal.Title,
			proposal.Description,
			proposal.Changes.ToParamChanges(),
			proposal.Height,
		),
		proposal.Deposit,
		fromInfo.GetAddress(),
	)

	return gc.BuildAndBroadcast(fromInfo.GetName(), passWd, memo, []sdk.Msg{msg}, accNum, seqNum)

}

// SubmitDelistProposal submits the proposal to delist a token pair from dex
func (gc govClient) SubmitDelistProposal(fromInfo keys.Info, passWd, proposalPath, memo string, accNum, seqNum uint64) (
	resp sdk.TxResponse, err error) {
	if err = params.CheckKeyParams(fromInfo, passWd); err != nil {
		return
	}

	proposal, err := dexutils.ParseDelistProposalJSON(gc.GetCodec(), proposalPath)
	if err != nil {
		return
	}

	msg := gov.NewMsgSubmitProposal(
		dextypes.NewDelistProposal(
			proposal.Title,
			proposal.Description,
			fromInfo.GetAddress(),
			proposal.BaseAsset,
			proposal.QuoteAsset,
		),
		proposal.Deposit,
		fromInfo.GetAddress(),
	)

	return gc.BuildAndBroadcast(fromInfo.GetName(), passWd, memo, []sdk.Msg{msg}, accNum, seqNum)

}

// SubmitCommunityPoolSpendProposal submits the proposal to spend the tokens from the community pool on OKChain
func (gc govClient) SubmitCommunityPoolSpendProposal(fromInfo keys.Info, passWd, proposalPath, memo string, accNum,
	seqNum uint64) (resp sdk.TxResponse, err error) {
	if err = params.CheckKeyParams(fromInfo, passWd); err != nil {
		return
	}

	proposal, err := parseCommunityPoolSpendProposalFromFile(proposalPath)
	if err != nil {
		return
	}

	msg := types.NewMsgSubmitProposal(
		types.NewCommunityPoolSpendProposal(
			proposal.Title,
			proposal.Description,
			proposal.Recipient,
			proposal.Amount,
		),
		proposal.Deposit,
		fromInfo.GetAddress(),
	)

	return gc.BuildAndBroadcast(fromInfo.GetName(), passWd, memo, []sdk.Msg{msg}, accNum, seqNum)

}

// Deposit increases the deposit amount on a specific proposal
func (gc govClient) Deposit(fromInfo keys.Info, passWd, depositCoinsStr, memo string, proposalID, accNum,
	seqNum uint64) (resp sdk.TxResponse, err error) {
	if err = params.CheckProposalOperation(fromInfo, passWd, proposalID); err != nil {
		return
	}

	deposit, err := sdk.ParseDecCoins(depositCoinsStr)
	if err != nil {
		return
	}

	msg := types.NewMsgDeposit(fromInfo.GetAddress(), proposalID, deposit)

	return gc.BuildAndBroadcast(fromInfo.GetName(), passWd, memo, []sdk.Msg{msg}, accNum, seqNum)

}

// Vote votes for an active proposal
// options: yes/no/no_with_veto/abstain
func (gc govClient) Vote(fromInfo keys.Info, passWd, voteOption, memo string, proposalID, accNum, seqNum uint64) (
	resp sdk.TxResponse, err error) {
	if err = params.CheckProposalOperation(fromInfo, passWd, proposalID); err != nil {
		return
	}

	voteOptionBytes, err := voteOptionFromString(voteOption)
	if err != nil {
		return
	}

	msg := types.NewMsgVote(fromInfo.GetAddress(), proposalID, voteOptionBytes)

	return gc.BuildAndBroadcast(fromInfo.GetName(), passWd, memo, []sdk.Msg{msg}, accNum, seqNum)

}
