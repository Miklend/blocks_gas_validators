package collect

import (
	"blocks_gas_validators/internal/miner/alchemy"
	"fmt"
	"math"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

func NewBlockMetrics(block *types.Block) alchemy.Block {
	txns := block.Transactions()
	gasStats := CalculateGasStats(txns)
	loc, _ := time.LoadLocation("Europe/Moscow")

	t := time.Unix(int64(block.Time()), 0).In(loc)
	return alchemy.Block{
		BlockNumber:       block.NumberU64(),
		BlockTime:         t,
		BlockTimestamp:    block.Time(),
		TransactionsCount: len(txns),
		BlockSizeBytes:    block.Size(),
		GasLimit:          block.GasLimit(),
		GasUsed:           block.GasUsed(),
		BlockFullness:     float64(block.GasUsed()) / float64(block.GasLimit()) * 100,
		Validator:         block.Coinbase().Hex(),
		GasStats:          gasStats,
	}
}

func CalculateGasStats(transactions types.Transactions) alchemy.GasStats {
	const weiToGwei = 1e9 // 1 gwei = 10^9 wei

	if len(transactions) == 0 {
		return alchemy.GasStats{
			Min:       0.0,
			Max:       0.0,
			Avg:       0.0,
			Stddev:    0.0,
			AllPrices: []float64{},
		}
	}

	var gasPrices []float64
	var min, max float64
	var sum float64

	for _, tx := range transactions {
		price := tx.GasPrice().Int64()
		priceGwei := float64(price) / weiToGwei

		gasPrices = append(gasPrices, priceGwei)

		if len(gasPrices) == 1 || priceGwei < min {
			min = priceGwei
		}
		if len(gasPrices) == 1 || priceGwei > max {
			max = priceGwei
		}

		sum += priceGwei
	}

	avg := sum / float64(len(transactions))

	var sumSq float64
	for _, g := range gasPrices {
		diff := g - avg
		sumSq += diff * diff
	}
	stddev := math.Sqrt(sumSq / float64(len(gasPrices)))

	return alchemy.GasStats{
		Min:       min,
		Max:       max,
		Avg:       avg,
		Stddev:    stddev,
		AllPrices: gasPrices,
	}
}

func NewBlockMetricsFromJSON(jsonBlock alchemy.JSONBlock) (alchemy.Block, error) {
	// Конвертация hex строк в числа
	blockNumber, err := hexutil.DecodeUint64(jsonBlock.Number)
	if err != nil {
		return alchemy.Block{}, fmt.Errorf("failed to parse block number: %w", err)
	}

	timestamp, err := hexutil.DecodeUint64(jsonBlock.Timestamp)
	if err != nil {
		return alchemy.Block{}, fmt.Errorf("failed to parse timestamp: %w", err)
	}

	blockSize, err := hexutil.DecodeUint64(jsonBlock.Size)
	if err != nil {
		return alchemy.Block{}, fmt.Errorf("failed to parse block size: %w", err)
	}

	gasLimit, err := hexutil.DecodeUint64(jsonBlock.GasLimit)
	if err != nil {
		return alchemy.Block{}, fmt.Errorf("failed to parse gas limit: %w", err)
	}

	gasUsed, err := hexutil.DecodeUint64(jsonBlock.GasUsed)
	if err != nil {
		return alchemy.Block{}, fmt.Errorf("failed to parse gas used: %w", err)
	}

	// Расчет статистики по gas
	gasStats := CalculateGasStatsFromJSON(jsonBlock.Transactions)

	// Конвертация времени
	loc, _ := time.LoadLocation("Europe/Moscow")
	t := time.Unix(int64(timestamp), 0).In(loc)

	return alchemy.Block{
		BlockNumber:       blockNumber,
		BlockTime:         t,
		BlockTimestamp:    timestamp,
		TransactionsCount: len(jsonBlock.Transactions),
		BlockSizeBytes:    blockSize,
		GasLimit:          gasLimit,
		GasUsed:           gasUsed,
		BlockFullness:     float64(gasUsed) / float64(gasLimit) * 100,
		Validator:         jsonBlock.Miner,
		GasStats:          gasStats,
	}, nil
}

func CalculateGasStatsFromJSON(transactions []alchemy.JSONTransaction) alchemy.GasStats {
	const weiToGwei = 1e9 // 1 gwei = 10^9 wei

	if len(transactions) == 0 {
		return alchemy.GasStats{
			Min:       0.0,
			Max:       0.0,
			Avg:       0.0,
			Stddev:    0.0,
			AllPrices: []float64{},
		}
	}

	var gasPrices []float64
	var min, max float64
	var sum float64

	for _, tx := range transactions {
		price, err := hexutil.DecodeBig(tx.GasPrice)
		if err != nil {
			continue // Пропускаем транзакции с невалидными ценами
		}

		priceGwei := float64(price.Int64()) / weiToGwei
		gasPrices = append(gasPrices, priceGwei)

		if len(gasPrices) == 1 || priceGwei < min {
			min = priceGwei
		}
		if len(gasPrices) == 1 || priceGwei > max {
			max = priceGwei
		}

		sum += priceGwei
	}

	if len(gasPrices) == 0 {
		return alchemy.GasStats{
			Min:       0.0,
			Max:       0.0,
			Avg:       0.0,
			Stddev:    0.0,
			AllPrices: []float64{},
		}
	}

	avg := sum / float64(len(gasPrices))

	var sumSq float64
	for _, g := range gasPrices {
		diff := g - avg
		sumSq += diff * diff
	}
	stddev := math.Sqrt(sumSq / float64(len(gasPrices)))

	return alchemy.GasStats{
		Min:       min,
		Max:       max,
		Avg:       avg,
		Stddev:    stddev,
		AllPrices: gasPrices,
	}
}
