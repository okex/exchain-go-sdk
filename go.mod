module github.com/okex/exchain-go-sdk

go 1.16

require (
	github.com/bartekn/go-bip39 v0.0.0-20171116152956-a05967ea095d
	github.com/ethereum/go-ethereum v1.10.8
	github.com/golang/mock v1.6.0
	github.com/okex/exchain v1.2.1-0.20220425224558-c1344e0c6085
	github.com/stretchr/testify v1.7.1
	github.com/tendermint/tendermint v0.34.14
)

replace (
	github.com/ethereum/go-ethereum => github.com/okex/go-ethereum v1.10.8-oec3
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	github.com/tendermint/go-amino => github.com/okex/go-amino v0.15.1-exchain6
//github.com/okex/exchain => ./third/exchain
)
