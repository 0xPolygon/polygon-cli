package p2p

import (
	"sync/atomic"
)

// Direction represents the direction of a message (sent or received).
type Direction string

const (
	// MsgReceived represents messages received from peers.
	MsgReceived Direction = "received"
	// MsgSent represents messages sent to peers.
	MsgSent Direction = "sent"

	// PacketSuffix is appended to message names to create packet count metrics.
	PacketSuffix = "Packet"
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
}

// Load takes a snapshot of all the counts in a thread-safe manner. Make sure
// you call this and read from the returned object.
func (mc *MessageCount) Load() MessageCount {
	return MessageCount{
		BlockHeaders:        atomic.LoadInt64(&mc.BlockHeaders),
		BlockBodies:         atomic.LoadInt64(&mc.BlockBodies),
		Blocks:              atomic.LoadInt64(&mc.Blocks),
		BlockHashes:         atomic.LoadInt64(&mc.BlockHashes),
		BlockHeaderRequests: atomic.LoadInt64(&mc.BlockHeaderRequests),
		BlockBodiesRequests: atomic.LoadInt64(&mc.BlockBodiesRequests),
		Transactions:        atomic.LoadInt64(&mc.Transactions),
		TransactionHashes:   atomic.LoadInt64(&mc.TransactionHashes),
		TransactionRequests: atomic.LoadInt64(&mc.TransactionRequests),
		Pings:               atomic.LoadInt64(&mc.Pings),
		Errors:              atomic.LoadInt64(&mc.Errors),
		Disconnects:         atomic.LoadInt64(&mc.Disconnects),
		NewWitness:          atomic.LoadInt64(&mc.NewWitness),
		NewWitnessHashes:    atomic.LoadInt64(&mc.NewWitnessHashes),
		GetWitnessRequest:   atomic.LoadInt64(&mc.GetWitnessRequest),
		Witness:             atomic.LoadInt64(&mc.Witness),
	}
}

// Clear clears all of the counts from the message counter.
func (mc *MessageCount) Clear() {
	atomic.StoreInt64(&mc.BlockHeaders, 0)
	atomic.StoreInt64(&mc.BlockBodies, 0)
	atomic.StoreInt64(&mc.Blocks, 0)
	atomic.StoreInt64(&mc.BlockHashes, 0)
	atomic.StoreInt64(&mc.BlockHeaderRequests, 0)
	atomic.StoreInt64(&mc.BlockBodiesRequests, 0)
	atomic.StoreInt64(&mc.Transactions, 0)
	atomic.StoreInt64(&mc.TransactionHashes, 0)
	atomic.StoreInt64(&mc.TransactionRequests, 0)
	atomic.StoreInt64(&mc.Pings, 0)
	atomic.StoreInt64(&mc.Errors, 0)
	atomic.StoreInt64(&mc.Disconnects, 0)
	atomic.StoreInt64(&mc.NewWitness, 0)
	atomic.StoreInt64(&mc.NewWitnessHashes, 0)
	atomic.StoreInt64(&mc.GetWitnessRequest, 0)
	atomic.StoreInt64(&mc.Witness, 0)
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
