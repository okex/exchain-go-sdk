package token

import (
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	"github.com/okex/okexchain-go-sdk/mocks"
	gosdktypes "github.com/okex/okexchain-go-sdk/types"
	"github.com/okex/okexchain/x/token"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	addr      = "okexchain1ntvyep3suq5z7789g7d5dejwzameu08m6gh7yl"
	name      = "alice"
	passWd    = "12345678"
	accPubkey = "okexchainpub17weu6qepq0ph2t3u697qar7rmdtdtqp4744jcprjd2h356zr0yh5vmw38a3my4vqjx5"
	mnemonic  = "giggle sibling fun arrow elevator spoon blood grocery laugh tortoise culture tool"
	memo      = "my memo"
	recAddr   = "okexchain1qeh2fz0a4t78ylesd4cyd2mwt5wcfnfj98ev0u"

	tokenSymbol           = "btc-000"
	defaultDesc           = "default description"
	defaultOriginalSymbol = "default original symbol"
	defaultWholeName      = "default whole name"
)

func TestTokenClient_QueryTokenInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	config, err := gosdktypes.NewClientConfig("testURL", "testChain", gosdktypes.BroadcastBlock, "", 200000,
		1.1, "0.00000001okt")
	require.NoError(t, err)
	mockCli := mocks.NewMockClient(t, ctrl, config)
	mockCli.RegisterModule(NewTokenClient(mockCli.MockBaseClient))

	originalTotalSupply, err := sdk.NewDecFromStr("10000000000")
	require.NoError(t, err)
	totalSupply, err := sdk.NewDecFromStr("20000000000")
	require.NoError(t, err)
	ownerAddr, err := sdk.AccAddressFromBech32(addr)
	require.NoError(t, err)

	expectedRet := mockCli.BuildTokenInfoBytes(defaultDesc, tokenSymbol, defaultOriginalSymbol, defaultWholeName,
		originalTotalSupply, totalSupply, ownerAddr, true, false, 0)
	expectedCdc := mockCli.GetCodec()

	mockCli.EXPECT().GetCodec().Return(expectedCdc).Times(2)
	expectedPath := fmt.Sprintf("custom/%s/info/%s", token.QuerierRoute, tokenSymbol)
	mockCli.EXPECT().Query(expectedPath, nil).Return(expectedRet, int64(1024), nil)

	tokensInfo, err := mockCli.Token().QueryTokenInfo("", tokenSymbol)
	require.NoError(t, err)

	require.Equal(t, 1, len(tokensInfo))
	require.Equal(t, defaultDesc, tokensInfo[0].Description)
	require.Equal(t, tokenSymbol, tokensInfo[0].Symbol)
	require.Equal(t, defaultOriginalSymbol, tokensInfo[0].OriginalSymbol)
	require.Equal(t, defaultWholeName, tokensInfo[0].WholeName)
	require.True(t, originalTotalSupply.Equal(tokensInfo[0].OriginalTotalSupply))
	require.True(t, totalSupply.Equal(tokensInfo[0].TotalSupply))
	require.True(t, tokensInfo[0].Mintable)
	require.True(t, ownerAddr.Equals(tokensInfo[0].Owner))
	require.Equal(t, 0, tokensInfo[0].Type)

	mockCli.EXPECT().Query(expectedPath, nil).Return(nil, int64(0), errors.New("default error"))
	_, err = mockCli.Token().QueryTokenInfo("", tokenSymbol)
	require.Error(t, err)

	expectedRet = mockCli.BuildTokenInfoBytes(defaultDesc, tokenSymbol, defaultOriginalSymbol, defaultWholeName,
		originalTotalSupply, totalSupply, ownerAddr, true, true, 0)

	expectedPath = fmt.Sprintf("custom/%s/tokens/%s", token.QuerierRoute, addr)
	mockCli.EXPECT().Query(expectedPath, nil).Return(expectedRet, int64(1024), nil)

	tokensInfo, err = mockCli.Token().QueryTokenInfo(addr, "")
	require.NoError(t, err)

	require.Equal(t, 1, len(tokensInfo))
	require.Equal(t, defaultDesc, tokensInfo[0].Description)
	require.Equal(t, tokenSymbol, tokensInfo[0].Symbol)
	require.Equal(t, defaultOriginalSymbol, tokensInfo[0].OriginalSymbol)
	require.Equal(t, defaultWholeName, tokensInfo[0].WholeName)
	require.True(t, originalTotalSupply.Equal(tokensInfo[0].OriginalTotalSupply))
	require.True(t, totalSupply.Equal(tokensInfo[0].TotalSupply))
	require.True(t, tokensInfo[0].Mintable)
	require.True(t, ownerAddr.Equals(tokensInfo[0].Owner))
	require.Equal(t, 0, tokensInfo[0].Type)

	_, err = mockCli.Token().QueryTokenInfo("", "")
	require.Error(t, err)

	mockCli.EXPECT().Query(expectedPath, nil).Return(nil, int64(1024), errors.New("default error"))
	_, err = mockCli.Token().QueryTokenInfo(addr, "")
	require.Error(t, err)
}
