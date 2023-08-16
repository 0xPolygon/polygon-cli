package p2p

import (
	"sync/atomic"
)

// MessageCount is used to help the outer goroutine to receive summary of the
// number and type of messages that were sent. This is used for distributed
// logging. It can be used to count the different types of messages received
// across all peer connections to provide a summary.
type MessageCount struct {
	BlockHeaders        int32 `json:",omitempty"`
	BlockBodies         int32 `json:",omitempty"`
	Blocks              int32 `json:",omitempty"`
	BlockHashes         int32 `json:",omitempty"`
	BlockHeaderRequests int32 `json:",omitempty"`
	BlockBodiesRequests int32 `json:",omitempty"`
	Transactions        int32 `json:",omitempty"`
	TransactionHashes   int32 `json:",omitempty"`
	TransactionRequests int32 `json:",omitempty"`
	Pings               int32 `json:",omitempty"`
	Errors              int32 `json:",omitempty"`
	Disconnects         int32 `json:",omitempty"`
}

// Load takes a snapshot of all the counts in a thread-safe manner. Make sure
// you call this and read from the returned object.
func (count *MessageCount) Load() MessageCount {
	return MessageCount{
		BlockHeaders:        atomic.LoadInt32(&count.BlockHeaders),
		BlockBodies:         atomic.LoadInt32(&count.BlockBodies),
		Blocks:              atomic.LoadInt32(&count.Blocks),
		BlockHashes:         atomic.LoadInt32(&count.BlockHashes),
		BlockHeaderRequests: atomic.LoadInt32(&count.BlockHeaderRequests),
		BlockBodiesRequests: atomic.LoadInt32(&count.BlockBodiesRequests),
		Transactions:        atomic.LoadInt32(&count.Transactions),
		TransactionHashes:   atomic.LoadInt32(&count.TransactionHashes),
		TransactionRequests: atomic.LoadInt32(&count.TransactionRequests),
		Pings:               atomic.LoadInt32(&count.Pings),
		Errors:              atomic.LoadInt32(&count.Errors),
		Disconnects:         atomic.LoadInt32(&count.Disconnects),
	}
}

// Clear clears all of the counts from the message counter.
func (count *MessageCount) Clear() {
	atomic.StoreInt32(&count.BlockHeaders, 0)
	atomic.StoreInt32(&count.BlockBodies, 0)
	atomic.StoreInt32(&count.Blocks, 0)
	atomic.StoreInt32(&count.BlockHashes, 0)
	atomic.StoreInt32(&count.BlockHeaderRequests, 0)
	atomic.StoreInt32(&count.BlockBodiesRequests, 0)
	atomic.StoreInt32(&count.Transactions, 0)
	atomic.StoreInt32(&count.TransactionHashes, 0)
	atomic.StoreInt32(&count.TransactionRequests, 0)
	atomic.StoreInt32(&count.Pings, 0)
	atomic.StoreInt32(&count.Errors, 0)
	atomic.StoreInt32(&count.Disconnects, 0)
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

func sum(ints ...int32) int32 {
	var sum int32 = 0
	for _, i := range ints {
		sum += i
	}

	return sum
}
