package services

import "context"

// BlockEvent represents events emitted by developer 2 node
type BlockEvent struct {
	Height    uint64
	Hash      string
	TxHashes  []string
	Confirmed bool
	Raw       []byte
}

type BlockchainClient interface {
	SubmitTransaction(ctx context.Context, txType string, payload []byte) (string, error)
	SubscribeBlocks(ctx context.Context) (<-chan BlockEvent, error)
	Close() error
}
