module github.com/okex/exchain-go-sdk

go 1.16

require (
	github.com/bartekn/go-bip39 v0.0.0-20171116152956-a05967ea095d
	github.com/ethereum/go-ethereum v1.10.8
	github.com/golang/mock v1.6.0
	github.com/kr/pretty v0.2.1 // indirect
	github.com/minio/highwayhash v1.0.1 // indirect
	github.com/okex/exchain v1.6.4
	github.com/prometheus/client_golang v1.11.1 // indirect
	github.com/stretchr/testify v1.8.0
	golang.org/x/net v0.0.0-20210903162142-ad29c8ab022f // indirect
)

replace (
	github.com/ethereum/go-ethereum => github.com/okex/go-ethereum v1.10.8-oec3
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	github.com/tendermint/go-amino => github.com/okex/go-amino v0.15.1-okc4
)
