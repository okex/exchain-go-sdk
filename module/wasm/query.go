package wasm

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/query"
	"github.com/okex/exchain/x/wasm/types"
)

func (c wasmClient) QueryListCode(pageReq *query.PageRequest) (*types.QueryCodesResponse, error) {
	return c.QueryClient.Codes(
		context.Background(),
		&types.QueryCodesRequest{
			Pagination: pageReq,
		},
	)
}

func (c wasmClient) QueryListContract(codeID uint64, pageReq *query.PageRequest) (*types.QueryContractsByCodeResponse, error) {
	return c.QueryClient.ContractsByCode(
		context.Background(),
		&types.QueryContractsByCodeRequest{
			CodeId:     codeID,
			Pagination: pageReq,
		},
	)
}

func (c wasmClient) QueryCode(codeID uint64) (*types.QueryCodeResponse, error) {
	return c.QueryClient.Code(
		context.Background(),
		&types.QueryCodeRequest{
			CodeId: codeID,
		},
	)
}

func (c wasmClient) QueryCodeInfo(codeID uint64) (*types.CodeInfoResponse, error) {
	res, err := c.QueryClient.Code(
		context.Background(),
		&types.QueryCodeRequest{
			CodeId: codeID,
		},
	)

	if err != nil {
		return nil, err
	}

	return res.CodeInfoResponse, nil
}

func (c wasmClient) QueryContractInfo(address string) (*types.QueryContractInfoResponse, error) {
	_, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return nil, err
	}

	return c.QueryClient.ContractInfo(
		context.Background(),
		&types.QueryContractInfoRequest{
			Address: address,
		},
	)
}

func (c wasmClient) QueryContractHistory(address string, pageReq *query.PageRequest) (*types.QueryContractHistoryResponse, error) {
	_, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return nil, err
	}

	return c.QueryClient.ContractHistory(
		context.Background(),
		&types.QueryContractHistoryRequest{
			Address:    address,
			Pagination: pageReq,
		},
	)
}

func (c wasmClient) QueryContractStateAll(address string, pageReq *query.PageRequest) (*types.QueryAllContractStateResponse, error) {
	_, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return nil, err
	}

	return c.QueryClient.AllContractState(
		context.Background(),
		&types.QueryAllContractStateRequest{
			Address:    address,
			Pagination: pageReq,
		},
	)
}

func (c wasmClient) QueryContractStateRaw(address string, queryData string) (*types.QueryRawContractStateResponse, error) {
	_, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return nil, err
	}

	decoder := newArgDecoder(hex.DecodeString)
	qData, err := decoder.DecodeString(queryData)
	if err != nil {
		return nil, err
	}

	return c.QueryClient.RawContractState(
		context.Background(),
		&types.QueryRawContractStateRequest{
			Address:   address,
			QueryData: qData,
		},
	)
}

func (c wasmClient) QueryContractStateSmart(address string, queryData string) (*types.QuerySmartContractStateResponse, error) {
	_, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return nil, err
	}
	if queryData == "" {
		return nil, errors.New("query data must not be empty")
	}

	decoder := newArgDecoder(asciiDecodeString)
	qData, err := decoder.DecodeString(queryData)
	if err != nil {
		return nil, fmt.Errorf("decode query: %s", err)
	}
	if !json.Valid(qData) {
		return nil, errors.New("query data must be json")
	}

	return c.QueryClient.SmartContractState(
		context.Background(),
		&types.QuerySmartContractStateRequest{
			Address:   address,
			QueryData: qData,
		},
	)
}

func (c wasmClient) QueryListPinnedCode(pageReq *query.PageRequest) (*types.QueryPinnedCodesResponse, error) {
	return c.QueryClient.PinnedCodes(
		context.Background(),
		&types.QueryPinnedCodesRequest{
			Pagination: pageReq,
		},
	)
}

type argumentDecoder struct {
	// dec is the default decoder
	dec                func(string) ([]byte, error)
	asciiF, hexF, b64F bool
}

func newArgDecoder(def func(string) ([]byte, error)) *argumentDecoder {
	return &argumentDecoder{dec: def}
}

func (a *argumentDecoder) DecodeString(s string) ([]byte, error) {
	found := -1
	for i, v := range []*bool{&a.asciiF, &a.hexF, &a.b64F} {
		if !*v {
			continue
		}
		if found != -1 {
			return nil, errors.New("multiple decoding flags used")
		}
		found = i
	}
	switch found {
	case 0:
		return asciiDecodeString(s)
	case 1:
		return hex.DecodeString(s)
	case 2:
		return base64.StdEncoding.DecodeString(s)
	default:
		return a.dec(s)
	}
}

func asciiDecodeString(s string) ([]byte, error) {
	return []byte(s), nil
}
