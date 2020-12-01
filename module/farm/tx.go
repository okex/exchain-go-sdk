package farm

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain-go-sdk/types/params"
	farmtypes "github.com/okex/okexchain/x/farm/types"
)

// CreatePool creates a farm pool
func (fc farmClient) CreatePool(fromInfo keys.Info, passWd, poolName, minLockAmountStr, yieldToken, memo string, accNum,
	seqNum uint64) (resp sdk.TxResponse, err error) {
	if err = params.CheckCreatePoolParams(fromInfo, passWd, poolName, minLockAmountStr, yieldToken); err != nil {
		return
	}

	minLockAmount, err := sdk.ParseDecCoin(minLockAmountStr)
	if err != nil {
		return
	}

	msg := farmtypes.NewMsgCreatePool(fromInfo.GetAddress(), poolName, minLockAmount, yieldToken)
	return fc.BuildAndBroadcast(fromInfo.GetName(), passWd, memo, []sdk.Msg{msg}, accNum, seqNum)
}
