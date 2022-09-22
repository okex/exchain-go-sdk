package main

import (
	gosdk "github.com/okex/exchain-go-sdk"
	"github.com/okex/exchain-go-sdk/utils"
	"github.com/okex/exchain/libs/cosmos-sdk/types/query"
	"log"
	"strconv"
	"strings"
)

const (
	// TODO: link to mainnet of ExChain later
	rpcURL = "tcp://52.199.88.250:26657"
	// user's name
	name = "alice"
	// user's mnemonic
	mnemonic  = "giggle sibling fun arrow elevator spoon blood grocery laugh tortoise culture tool"
	mnemonic2 = "antique onion adult slot sad dizzy sure among cement demise submit scare"
	// user's password
	passWd = "12345678"
	// target address
	addr     = "ex1qj5c07sm6jetjz8f509qtrxgh4psxkv3ddyq7u"
	baseCoin = "okt"

	addr2 = "ex1fsfwwvl93qv6r56jpu084hxxzn9zphnyxhske5"
)

func main() {
	//-------------------- 1. preparation --------------------//
	// NOTE: either of the both ways below to pay fees is available

	// WAY 1: create a client config with fixed fees
	config, err := gosdk.NewClientConfig(rpcURL, "exchain-64", gosdk.BroadcastBlock, "0.01okt", 100000000,
		0, "")
	if err != nil {
		log.Fatal(err)
	}

	// WAY 2: alternative client config with the fees by auto gas calculation
	//config, err = gosdk.NewClientConfig(rpcURL, "exchain-64", gosdk.BroadcastBlock, "", 200000,
	//	1.1, "0.00000000000001okt")
	//if err != nil {
	//	log.Fatal(err)
	//}

	cli := gosdk.NewClient(config)

	// create an account with your own mnemonic，name and password
	fromInfo, _, err := utils.CreateAccountWithMnemo(mnemonic, name, passWd)
	if err != nil {
		log.Fatal(err)
	}

	fromInfo2, _, err := utils.CreateAccountWithMnemo(mnemonic2, "bob", passWd)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("account 2 address: ", fromInfo2.GetAddress().String())

	//-------------------- 2. query for the information of your address --------------------//

	accInfo, err := cli.Auth().QueryAccount(fromInfo.GetAddress().String())
	if err != nil {
		log.Fatal(err)
	}

	log.Println(accInfo)

	//-------------------- 3. transfer to other address --------------------//

	//sequence number of the account must be increased by 1 whenever a transaction of the account takes effect
	accountNum, sequenceNum := accInfo.GetAccountNumber(), accInfo.GetSequence()

	wasmFile := "/Users/finefine/workspace/cosmos/cw-contracts/contracts/nameservice/target/wasm32-unknown-unknown/release/cw_nameservice.wasm"
	res, err := cli.Wasm().StoreCode(fromInfo, passWd, accountNum, sequenceNum, "memo", wasmFile, "", false, false)
	if err != nil {
		log.Fatal(err)
	}

	index := strings.LastIndex(res.RawLog, ":")
	codeIDStr := res.RawLog[index:]
	codeIDStr = codeIDStr[2 : strings.Index(codeIDStr, "}")-1]
	codeID, _ := strconv.Atoi(codeIDStr)

	log.Println("=============================================================StoreCode===============================================================")
	log.Println(codeID)
	log.Println(res)

	// instantiate a wasm contract
	sequenceNum++
	initMsg := `{"purchase_price":{"amount":"1","denom":"okt"},"transfer_price":{"amount":"1","denom":"okt"}}`
	instantiateRes, err := cli.Wasm().InstantiateContract(fromInfo, passWd, accountNum, sequenceNum, "memo", uint64(codeID), initMsg, "1okt", "name service", "ex1qj5c07sm6jetjz8f509qtrxgh4psxkv3ddyq7u", false)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("=========================================================InstantiateContract==========================================================")
	log.Println(instantiateRes)
	index = strings.Index(instantiateRes.RawLog, "address")
	contractAddr := instantiateRes.RawLog[index+18 : index+18+61]
	log.Println("contract address: ", contractAddr)

	// execute a wasm contract
	sequenceNum++
	execMsg := `{"register":{"name":"fred"}}`
	executeRes, err := cli.Wasm().ExecuteContract(fromInfo, passWd, accountNum, sequenceNum, "memo", contractAddr, execMsg, "2okt")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("=========================================================ExecuteContract==========================================================")
	log.Println(executeRes)

	// set new admin for the contract
	sequenceNum++
	updateAdminRes, err := cli.Wasm().UpdateContractAdmin(fromInfo, passWd, accountNum, sequenceNum, "memo", contractAddr, addr2)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("=========================================================UpdateContractAdmin==========================================================")
	log.Println(updateAdminRes)

	// migrate contract to the new code
	accInfo2, err := cli.Auth().QueryAccount(fromInfo2.GetAddress().String())
	if err != nil {
		log.Fatal(err)
	}

	log.Println(accInfo2)
	//migrateMsg := initMsg
	//migrateRes, err := cli.Wasm().MigrateContract(fromInfo2, passWd, accInfo2.GetAccountNumber(), accInfo2.GetSequence(), "memo", 2, contractAddr, migrateMsg)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Println("=========================================================MigrateContract==========================================================")
	//log.Println(migrateRes)

	// clear admin for the contract
	clearAdminRes, err := cli.Wasm().ClearContractAdmin(fromInfo2, passWd, accInfo2.GetAccountNumber(), accInfo2.GetSequence(), "memo", contractAddr)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("=========================================================ClearContractAdmin==========================================================")
	log.Println(clearAdminRes)

	// query code
	queryCodeRes, err := cli.Wasm().QueryCode(uint64(codeID))
	log.Println("=========================================================QueryCode==========================================================")
	if err != nil {
		log.Println(err)
	} else {
		log.Println(queryCodeRes.DataHash)
	}

	// 	query listCode
	listCodeRes, err := cli.Wasm().QueryListCode(&query.PageRequest{
		Key:        nil,
		Offset:     0,
		Limit:      50,
		CountTotal: false,
	})

	log.Println("=========================================================QueryListCode==========================================================")
	if err != nil {
		log.Println(err)
	} else {
		log.Println(len(listCodeRes.CodeInfos))
	}

	// query ListContractByCode
	listContract, err := cli.Wasm().QueryListContract(uint64(codeID), &query.PageRequest{
		Key:        nil,
		Offset:     0,
		Limit:      50,
		CountTotal: false,
	})

	log.Println("=========================================================QueryListContract==========================================================")
	if err != nil {
		log.Println(err)
	} else {
		log.Println(listContract)
	}

	// query code info
	codeInfo, err := cli.Wasm().QueryCodeInfo(uint64(codeID))
	log.Println("=========================================================QueryCodeInfo==========================================================")
	if err != nil {
		log.Println(err)
	} else {
		log.Println(codeInfo)
	}

	// query contract info
	contractInfo, err := cli.Wasm().QueryContractInfo(contractAddr)
	log.Println("=========================================================QueryContractInfo==========================================================")
	if err != nil {
		log.Println(err)
	} else {
		log.Println(contractInfo)
	}

	// query contract history
	contractHistory, err := cli.Wasm().QueryContractHistory(contractAddr, &query.PageRequest{
		Key:        nil,
		Offset:     0,
		Limit:      50,
		CountTotal: false,
	})
	log.Println("=========================================================QueryContractHistory==========================================================")
	if err != nil {
		log.Println(err)
	} else {
		log.Println(contractHistory)
	}

	// query contract state all
	contractStateAll, err := cli.Wasm().QueryContractStateAll(contractAddr, &query.PageRequest{
		Key:        nil,
		Offset:     0,
		Limit:      50,
		CountTotal: false,
	})
	log.Println("=========================================================QueryContractStateAll==========================================================")
	if err != nil {
		log.Println(err)
	} else {
		log.Println(contractStateAll)
	}

	// query contract state raw
	contractStateRaw, err := cli.Wasm().QueryContractStateRaw(contractAddr, "")
	log.Println("=========================================================QueryContractStateRaw==========================================================")
	if err != nil {
		log.Println(err)
	} else {
		log.Println(contractStateRaw)
	}

	// query contract state smart
	contractStateSmart, err := cli.Wasm().QueryContractStateSmart(contractAddr, "")
	log.Println("=========================================================QueryContractStateSmart==========================================================")
	if err != nil {
		log.Println(err)
	} else {
		log.Println(contractStateSmart)
	}

	// query contract ListPinnedCode
	pinnedCode, err := cli.Wasm().QueryListPinnedCode(&query.PageRequest{
		Key:        nil,
		Offset:     0,
		Limit:      50,
		CountTotal: false,
	})
	log.Println("=========================================================QueryListPinnedCode==========================================================")
	if err != nil {
		log.Println(err)
	} else {
		log.Println(pinnedCode)
	}
}
