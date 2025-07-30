package p2p

import (
	"sync/atomic"
)

// MessageCount is used to help the outer goroutine to receive summary of the
// number and type of messages that were sent. This is used for distributed
// logging. It can be used to count the different types of messages received
// across all peer connections to provide a summary.
type MessageCount struct {
	BlockHeaders        int64 `json:"block_headers,omitempty"`
	BlockBodies         int64 `json:"block_bodies,omitempty"`
	Blocks              int64 `json:"blocks,omitempty"`
	BlockHashes         int64 `json:"block_hashes,omitempty"`
	BlockHeaderRequests int64 `json:"block_header_requests,omitempty"`
	BlockBodiesRequests int64 `json:"block_bodies_requests,omitempty"`
	Transactions        int64 `json:"transactions,omitempty"`
	TransactionHashes   int64 `json:"transaction_hashes,omitempty"`
	TransactionRequests int64 `json:"transaction_requests,omitempty"`
	Pings               int64 `json:"pings,omitempty"`
	Errors              int64 `json:"errors,omitempty"`
	Disconnects         int64 `json:"disconnects,omitempty"`
	NewWitness          int64 `json:"new_witness,omitempty"`
	NewWitnessHashes    int64 `json:"new_witness_hashes,omitempty"`
	GetWitnessRequest   int64 `json:"get_witness_request,omitempty"`
	Witness             int64 `json:"witness,omitempty"`
	WitnessDataBytes    int64 `json:"witness_data_bytes,omitempty"`
	WitnessLatencyMs    int64 `json:"witness_latency_ms,omitempty"`
	WitnessPagesFetched int64 `json:"witness_pages_fetched,omitempty"`
}

// Load takes a snapshot of all the counts in a thread-safe manner. Make sure
// you call this and read from the returned object.
func (count *MessageCount) Load() MessageCount {
	return MessageCount{
		BlockHeaders:        atomic.LoadInt64(&count.BlockHeaders),
		BlockBodies:         atomic.LoadInt64(&count.BlockBodies),
		Blocks:              atomic.LoadInt64(&count.Blocks),
		BlockHashes:         atomic.LoadInt64(&count.BlockHashes),
		BlockHeaderRequests: atomic.LoadInt64(&count.BlockHeaderRequests),
		BlockBodiesRequests: atomic.LoadInt64(&count.BlockBodiesRequests),
		Transactions:        atomic.LoadInt64(&count.Transactions),
		TransactionHashes:   atomic.LoadInt64(&count.TransactionHashes),
		TransactionRequests: atomic.LoadInt64(&count.TransactionRequests),
		Pings:               atomic.LoadInt64(&count.Pings),
		Errors:              atomic.LoadInt64(&count.Errors),
		Disconnects:         atomic.LoadInt64(&count.Disconnects),
		NewWitness:          atomic.LoadInt64(&count.NewWitness),
		NewWitnessHashes:    atomic.LoadInt64(&count.NewWitnessHashes),
		GetWitnessRequest:   atomic.LoadInt64(&count.GetWitnessRequest),
		Witness:             atomic.LoadInt64(&count.Witness),
		WitnessDataBytes:    atomic.LoadInt64(&count.WitnessDataBytes),
		WitnessLatencyMs:    atomic.LoadInt64(&count.WitnessLatencyMs),
		WitnessPagesFetched: atomic.LoadInt64(&count.WitnessPagesFetched),
	}
}

// Clear clears all of the counts from the message counter.
func (count *MessageCount) Clear() {
	atomic.StoreInt64(&count.BlockHeaders, 0)
	atomic.StoreInt64(&count.BlockBodies, 0)
	atomic.StoreInt64(&count.Blocks, 0)
	atomic.StoreInt64(&count.BlockHashes, 0)
	atomic.StoreInt64(&count.BlockHeaderRequests, 0)
	atomic.StoreInt64(&count.BlockBodiesRequests, 0)
	atomic.StoreInt64(&count.Transactions, 0)
	atomic.StoreInt64(&count.TransactionHashes, 0)
	atomic.StoreInt64(&count.TransactionRequests, 0)
	atomic.StoreInt64(&count.Pings, 0)
	atomic.StoreInt64(&count.Errors, 0)
	atomic.StoreInt64(&count.Disconnects, 0)
	atomic.StoreInt64(&count.NewWitness, 0)
	atomic.StoreInt64(&count.NewWitnessHashes, 0)
	atomic.StoreInt64(&count.GetWitnessRequest, 0)
	atomic.StoreInt64(&count.Witness, 0)
	atomic.StoreInt64(&count.WitnessDataBytes, 0)
	atomic.StoreInt64(&count.WitnessLatencyMs, 0)
	atomic.StoreInt64(&count.WitnessPagesFetched, 0)
}

// IsEmpty checks whether the sum of all the counts is empty. Make sure to call
// Load before this method to get an accurate count.
func (c *MessageCount) IsEmpty() bool {
	return sum(
		c.BlockHeaders,
		c.BlockBodies,
		c.BlockHashes,
		c.BlockHeaderRequests,
		c.BlockBodiesRequests,
		c.Transactions,
		c.TransactionHashes,
		c.TransactionRequests,
		c.Pings,
		c.Errors,
		c.Disconnects,
	) == 0
}

func sum(ints ...int64) int64 {
	var sum int64 = 0
	for _, i := range ints {
		sum += i
	}

	return sum
}
