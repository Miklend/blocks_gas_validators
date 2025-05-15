package alchemy

import "context"

type Storage interface {
	InsertBlocksBatch(ctx context.Context, blocks []*Block, chain string) error
	InsertBlocksCopy(ctx context.Context, blocks []*Block, chain string) error
	Create(ctx context.Context, block *Block, chain string) error
}
