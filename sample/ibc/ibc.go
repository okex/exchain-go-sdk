package main

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/crypto"
	gosdk "github.com/okex/exchain-go-sdk"
	"github.com/okex/exchain-go-sdk/exposed"
	"github.com/okex/exchain-go-sdk/utils"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	query2 "github.com/okex/exchain/libs/cosmos-sdk/types/query"
	client_types "github.com/okex/exchain/libs/ibc-go/modules/core/02-client/types"
	chantypes "github.com/okex/exchain/libs/ibc-go/modules/core/04-channel/types"
	//secp256k12 "github.com/okex/exchain/libs/cosmos-sdk/crypto/keys/ibc-key"
	secp256k12 "github.com/okex/exchain/libs/tendermint/crypto/secp256k1"
	"log"
)

const (
	rpcURL = "tcp://127.0.0.1:26657"
	// user's name
	name = "admin17"
	// user's mnemonic
	mnemonic = ""
	// user's password
	passWd = "12345678"
	// target address
	addr            = "cosmos1n064mg7jcxt2axur29mmek5ys7ghta4u4mhcjp"
	baseCoin        = "okt"
	aliceKey string = ""
)

func main() {
	//-------------------- 1. preparation --------------------//
	// NOTE: either of the both ways below to pay fees is available

	// WAY 1: create a client config with fixed fees
	//config, err := gosdk.NewClientConfig(rpcURL, "exchain-101", gosdk.BroadcastSync, "0.000001okt", 10000,
	//	0, "")
	//if err != nil {
	//	log.Fatal(err)
	//}

	// WAY 2: alternative client config with the fees by auto gas calculation
	config, err := gosdk.NewClientConfig(rpcURL, "exchain-100", gosdk.BroadcastBlock, "0.000001okt", 450000,
		1.1, "0.000000000000000012okt")
	if err != nil {
		log.Fatal(err)
	}

	cli := gosdk.NewClient(config)

	// create an account with your own mnemonic，name and password

	if err != nil {
		log.Fatal(err)
	}

	//-------------------- 2. query for the information of your address --------------------//

	accInfo, err := cli.Auth().QueryAccount("ex1hr26cyc335g7p5e948a7vkmwnx3fmxfzwdyryf")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(accInfo)

	testTransfer(cli.Ibc(), accInfo.GetAccountNumber(), accInfo.GetSequence())
	testQueryDenomTrace(cli.Ibc())
	testQueryDenomTraces(cli.Ibc())
	testQueryParmas(cli.Ibc())
	testQueryEscrowAddress(cli.Ibc())
	testQueryChannel(cli.Ibc())
	testQueryChannels(cli.Ibc())
	testQueryConnectionChannels(cli.Ibc())
	testQueryChannelClientState(cli.Ibc())
	testQueryChannelConsensusState(cli.Ibc())
	testQueryPackCommitment(cli.Ibc())
	testQueryPacketCommitments(cli.Ibc())
	testQueryPacketReceipt(cli.Ibc())
	testQueryPacketAcknowledgement(cli.Ibc())
	testQueryPacketAcknowledgements(cli.Ibc())
	testQueryUnreceivedPackets(cli.Ibc())
	testQueryUnreceivedAcks(cli.Ibc())
	testQueryNextSequenceReceive(cli.Ibc())
	testQueryTx(cli.Ibc())
	testQueryTxs(cli.Ibc())
	testQueryHeader(cli.Ibc())
}

func testTransfer(ibc exposed.Ibc, accountNum, sequenceNum uint64) {
	// sequence number of the account must be increased by 1 whenever a transaction of the account takes effect
	priStr, err := utils.GeneratePrivateKeyFromMnemo(mnemonic)
	if err != nil {
		log.Fatal(err)
	}

	key, err := crypto.HexToECDSA(priStr)
	if err != nil {
		log.Fatal(err)
	}

	d := crypto.FromECDSA(key)

	fee := sdk.NewCoinAdapter("wei", sdk.NewInt(45000000000000))
	fees := []sdk.CoinAdapter{fee}

	res, err := ibc.Transfer(secp256k12.GenPrivKeySecp256k1(d), "channel-0", addr, "1okt", fees, "memo", client_types.Height{RevisionNumber: 101, RevisionHeight: 10000})
	if err != nil {
		log.Fatal(err)
	}

	log.Println(res)
}

func testQueryDenomTrace(ibc exposed.Ibc) {
	log.Println("testQueryDenomTrace ======================================== ")
	query, err := ibc.QueryDenomTrace("CD3872E1E59BAA23BDAB04A829035D4988D6397569EC77F1DC991E4520D4092B")
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(query)
}

func testQueryDenomTraces(ibc exposed.Ibc) {
	log.Println("testQueryDenomTraces ======================================== ")
	query, err := ibc.QueryDenomTraces(&query2.PageRequest{
		Key:        nil,
		Offset:     0,
		Limit:      0,
		CountTotal: false,
	})
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(query)
}

func testQueryParmas(ibc exposed.Ibc) {
	log.Println("testQueryParmas ======================================== ")
	queryParams, err := ibc.QueryIbcParams()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(queryParams)
}

func testQueryEscrowAddress(ibc exposed.Ibc) {
	log.Println("testQueryEscrowAddress ======================================== ")
	queryEscrowAddress := ibc.QueryEscrowAddress("transfer", "channel-0")
	log.Println(queryEscrowAddress)
}

func testQueryChannel(ibc exposed.Ibc) {
	log.Println("testQueryChannel ======================================== ")
	channel, err := ibc.QueryChannel(&chantypes.QueryChannelRequest{
		PortId:    "transfer",
		ChannelId: "channel-0",
	})
	if err != nil {
		log.Fatalln(err)
		return
	}
	log.Println(channel)
}

func testQueryChannels(ibc exposed.Ibc) {
	log.Println("testQueryChannels ======================================== ")
	channels, err := ibc.QueryChannels()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(channels)
}

func testQueryConnectionChannels(ibc exposed.Ibc) {
	log.Println("testQueryConnectionChannels ======================================== ")
	res, err := ibc.ConnectionChannels(&chantypes.QueryConnectionChannelsRequest{
		Connection: "connection-0",
		Pagination: nil,
	})

	if err != nil {
		log.Println(err)
		return
	}

	log.Println(res)
}

func testQueryChannelClientState(ibc exposed.Ibc) {
	log.Println("testQueryChannelClientState ======================================== ")
	res, err := ibc.ChannelClientState(&chantypes.QueryChannelClientStateRequest{
		PortId:    "transfer",
		ChannelId: "channel-0",
	})
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(res.IdentifiedClientState.GetClientState())

}

func testQueryChannelConsensusState(ibc exposed.Ibc) {
	log.Println("testQueryChannelConsensusState ======================================== ")
	res, err := ibc.ChannelClientState(&chantypes.QueryChannelClientStateRequest{
		PortId:    "transfer",
		ChannelId: "channel-0",
	})
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(res.IdentifiedClientState.GetClientState())
}

func testQueryPackCommitment(ibc exposed.Ibc) {
	log.Println("testQueryPackCommitment ======================================== ")
	res, err := ibc.PacketCommitment(&chantypes.QueryPacketCommitmentRequest{
		PortId:    "transfer",
		ChannelId: "channel-0",
		Sequence:  1,
	})

	if err != nil {
		log.Println(err)
		return
	}

	log.Println(res)

}

func testQueryPacketCommitments(ibc exposed.Ibc) {
	log.Println("testQueryPacketCommitments ======================================== ")
	res, err := ibc.PacketCommitments(&chantypes.QueryPacketCommitmentsRequest{
		PortId:     "transfer",
		ChannelId:  "channel-0",
		Pagination: nil,
	})

	if err != nil {
		log.Println(err)
		return
	}

	log.Println(res)
}

func testQueryPacketReceipt(ibc exposed.Ibc) {
	log.Println("testQueryPacketReceipt ======================================== ")
	res, err := ibc.PacketReceipt(&chantypes.QueryPacketReceiptRequest{
		PortId:    "transfer",
		ChannelId: "channel-0",
		Sequence:  1,
	})

	if err != nil {
		log.Println(err)
		return
	}

	log.Println(res)
}

func testQueryPacketAcknowledgement(ibc exposed.Ibc) {
	log.Println("testQueryPacketAcknowledgement ======================================== ")
	res, err := ibc.PacketAcknowledgement(&chantypes.QueryPacketAcknowledgementRequest{
		PortId:    "transfer",
		ChannelId: "channel-0",
		Sequence:  1,
	})

	if err != nil {
		log.Println(err)
		return
	}

	log.Println(res)
}

func testQueryPacketAcknowledgements(ibc exposed.Ibc) {
	log.Println("testQueryPacketAcknowledgements ======================================== ")
	res, err := ibc.PacketAcknowledgements(&chantypes.QueryPacketAcknowledgementsRequest{
		PortId:                    "transfer",
		ChannelId:                 "channel-0",
		Pagination:                nil,
		PacketCommitmentSequences: []uint64{1},
	})

	if err != nil {
		log.Println(err)
		return
	}

	log.Println(res)
}

func testQueryUnreceivedPackets(ibc exposed.Ibc) {
	log.Println("testQueryUnreceivedPackets ======================================== ")
	res, err := ibc.UnreceivedPackets(&chantypes.QueryUnreceivedPacketsRequest{
		PortId:                    "transfer",
		ChannelId:                 "channel-0",
		PacketCommitmentSequences: []uint64{1},
	})

	if err != nil {
		log.Println(err)
		return
	}

	log.Println(res)
}

func testQueryUnreceivedAcks(ibc exposed.Ibc) {
	log.Println("testQueryUnreceivedAcks ======================================== ")
	res, err := ibc.UnreceivedAcks(&chantypes.QueryUnreceivedAcksRequest{
		PortId:             "transfer",
		ChannelId:          "channel-0",
		PacketAckSequences: []uint64{1},
	})

	if err != nil {
		log.Println(err)
		return
	}

	log.Println(res)
}

func testQueryNextSequenceReceive(ibc exposed.Ibc) {
	log.Println("testNextSequenceReceive ======================================== ")
	res, err := ibc.NextSequenceReceive(&chantypes.QueryNextSequenceReceiveRequest{
		PortId:    "transfer",
		ChannelId: "channel-0",
	})

	if err != nil {
		log.Println(err)
		return
	}

	log.Println(res)
}

func testQueryTx(ibc exposed.Ibc) {
	log.Println("testQueryTx ======================================== ")

	res, err := ibc.QueryTx("D7702BCC93BC3CA3C16EB0F9B1F945D33D1860B931B78FDDB6F0517B120E5E91")

	if err != nil {
		log.Println(err)
		return
	}

	resBytes, err := json.Marshal(res)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(string(resBytes))
}

func testQueryTxs(ibc exposed.Ibc) {
	log.Println("testQueryTxs ======================================== ")
	events := []string{"message.action=transfer"}
	res, err := ibc.QueryTxs(1, 5000, events)

	if err != nil {
		log.Println(err)
		return
	}

	resBytes, err := json.Marshal(res)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(string(resBytes))
}

func testQueryHeader(ibc exposed.Ibc) {
	log.Println("testQueryHeader ======================================== ")

	res, err := ibc.QueryHeaderAtHeight(1)

	if err != nil {
		log.Println(err)
		return
	}

	resBytes, err := json.Marshal(res)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(string(resBytes))
}
