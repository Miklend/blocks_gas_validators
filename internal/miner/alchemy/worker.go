package alchemy

import (
	"context"
	"sync"
)

type Worker interface {
	LastRun(ctx context.Context, in <-chan *Block)
	HistoryBatch(ctx context.Context, in <-chan []*Block, wg *sync.WaitGroup)
}
