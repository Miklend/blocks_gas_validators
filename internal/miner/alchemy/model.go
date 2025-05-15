package alchemy

import "time"

type Block struct {
	BlockNumber       uint64    `json:"block_number"`
	BlockTime         time.Time `json:"block_time"`
	BlockTimestamp    uint64    `json:"block_timestamp"`
	TransactionsCount int       `json:"transactions_count"`
	BlockSizeBytes    uint64    `json:"block_size_bytes"`
	GasLimit          uint64    `json:"gas_limit"`
	GasUsed           uint64    `json:"gas_used"`
	BlockFullness     float64   `json:"block_fullness"`
	Validator         string    `json:"validator"`
	GasStats          GasStats  `json:"gas_stats"`
}

type GasStats struct {
	Min       float64   `json:"gas_min"`
	Max       float64   `json:"gas_max"`
	Avg       float64   `json:"gas_avg"`
	Stddev    float64   `json:"gas_stddev"`
	AllPrices []float64 `json:"all_prices"`
}

type JSONBlock struct {
	Number       string            `json:"number"`
	Timestamp    string            `json:"timestamp"`
	Transactions []JSONTransaction `json:"transactions"`
	Size         string            `json:"size"`
	GasLimit     string            `json:"gasLimit"`
	GasUsed      string            `json:"gasUsed"`
	Miner        string            `json:"miner"`
}

type JSONTransaction struct {
	GasPrice string `json:"gasPrice"`
}
