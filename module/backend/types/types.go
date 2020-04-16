package types

// const
const (
	ModuleName = "backend"

	CandlesPath        = "custom/backend/candles"
	TickersPath        = "custom/backend/tickers"
	RecentTxRecordPath = "custom/backend/matches"
)

// Ticker - structure of ticker's detail data
type Ticker struct {
	Symbol           string `json:"symbol"`
	Product          string `json:"product"`
	Timestamp        string `json:"timestamp"`
	Open             string `json:"open"`
	Close            string `json:"close"`
	High             string `json:"high"`
	Low              string `json:"low"`
	Price            string `json:"price"`
	Volume           string `json:"volume"`
	Change           string `json:"change"`
	ChangePercentage string `json:"change_percentage"`
}

// MatchResult - structure for recent tx record
type MatchResult struct {
	Timestamp   int64   `json:"timestamp"`
	BlockHeight int64   `json:"block_height"`
	Product     string  `json:"product"`
	Price       float64 `json:"price"`
	Quantity    float64 `json:"volume"`
}

// BaseResponse - structure for base response of data
type BaseResponse struct {
	Code      int         `json:"code"`
	Msg       string      `json:"msg"`
	DetailMsg string      `json:"detail_msg"`
	Data      interface{} `json:"data"`
}

// ParamPage - structure of page params
type ParamPage struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
	Total   int `json:"total"`
}

// ListDataRes - structure of list data in the list response
type ListDataRes struct {
	Data      interface{} `json:"data"`
	ParamPage ParamPage   `json:"param_page"`
}

// ListResponse - structure for list response of data
type ListResponse struct {
	Code      int         `json:"code"`
	Msg       string      `json:"msg"`
	DetailMsg string      `json:"detail_msg"`
	Data      ListDataRes `json:"data"`
}