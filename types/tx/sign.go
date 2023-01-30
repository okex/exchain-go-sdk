package tx

import (
	"github.com/okex/exchain/app/crypto/hd"
	"github.com/okex/exchain/libs/cosmos-sdk/crypto/keys"
	authtypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	"github.com/okex/exchain/libs/tendermint/crypto"
)

var (
	// Kb - global keybase
	Kb keys.Keybase
)

func init() {
	Kb = keys.NewInMemory(hd.EthSecp256k1Options()...)
}

// MakeSignature completes the signature
func MakeSignature(name, passphrase string, msg authtypes.StdSignMsg) (sig authtypes.StdSignature, err error) {
	sigBytes, pubkey, err := Kb.Sign(name, passphrase, msg.Bytes())
	if err != nil {
		return
	}
	return authtypes.StdSignature{
		PubKey:    pubkey,
		Signature: sigBytes,
	}, nil
}

func ExportPrivateKeyObject(name, passphrase string) (crypto.PrivKey, error) {
	return Kb.ExportPrivateKeyObject(name, passphrase)
}
