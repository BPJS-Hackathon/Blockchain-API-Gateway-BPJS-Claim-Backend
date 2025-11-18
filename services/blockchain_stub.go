package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"time"
)

// Simple in-memory stub client that returns fake txHash and emits confirmed events after delay
type StubBlockchainClient struct {
	events chan BlockEvent
	done   chan struct{}
}

func NewStubBlockchainClient() *StubBlockchainClient {
	return &StubBlockchainClient{
		events: make(chan BlockEvent, 10),
		done:   make(chan struct{}),
	}
}

func (s *StubBlockchainClient) SubmitTransaction(ctx context.Context, txType string, payload []byte) (string, error) {
	// create pseudo tx hash
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	txHash := hex.EncodeToString(b)

	// simulate eventual confirmation after 5 seconds in background
	go func(tx string) {
		time.Sleep(5 * time.Second)
		ev := BlockEvent{
			Height:    uint64(time.Now().Unix()),
			Hash:      "block-" + tx[:8],
			TxHashes:  []string{tx},
			Confirmed: true,
			Raw:       nil,
		}
		select {
		case s.events <- ev:
		default:
			log.Println("stub: event queue full")
		}
	}(txHash)

	return txHash, nil
}

func (s *StubBlockchainClient) SubscribeBlocks(ctx context.Context) (<-chan BlockEvent, error) {
	return s.events, nil
}

func (s *StubBlockchainClient) Close() error {
	close(s.done)
	close(s.events)
	return nil
}
