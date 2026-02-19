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

// IncrementByName increments the appropriate field based on message name.
func (mc *MessageCount) IncrementByName(name string, count int64) {
	switch name {
	case "BlockHeaders":
		atomic.AddInt64(&mc.BlockHeaders, count)
	case "BlockBodies":
		atomic.AddInt64(&mc.BlockBodies, count)
	case "NewBlock":
		atomic.AddInt64(&mc.Blocks, count)
	case "NewBlockHashes":
		atomic.AddInt64(&mc.BlockHashes, count)
	case "GetBlockHeaders":
		atomic.AddInt64(&mc.BlockHeaderRequests, count)
	case "GetBlockBodies":
		atomic.AddInt64(&mc.BlockBodiesRequests, count)
	case "Transactions", "PooledTransactions":
		atomic.AddInt64(&mc.Transactions, count)
	case "NewPooledTransactionHashes":
		atomic.AddInt64(&mc.TransactionHashes, count)
	case "GetPooledTransactions":
		atomic.AddInt64(&mc.TransactionRequests, count)
	case "Ping":
		atomic.AddInt64(&mc.Pings, count)
	case "NewWitness":
		atomic.AddInt64(&mc.NewWitness, count)
	case "NewWitnessHashes":
		atomic.AddInt64(&mc.NewWitnessHashes, count)
	case "GetWitness":
		atomic.AddInt64(&mc.GetWitnessRequest, count)
	case "Witness":
		atomic.AddInt64(&mc.Witness, count)
	}
}

// PeerMessages tracks message counts for a single peer connection.
// This is used to provide per-peer visibility via the API without
// creating high-cardinality Prometheus metrics.
type PeerMessages struct {
	Received        MessageCount
	Sent            MessageCount
	PacketsReceived MessageCount
	PacketsSent     MessageCount
}

// NewPeerMessages creates a new PeerMessages instance.
func NewPeerMessages() *PeerMessages {
	return &PeerMessages{}
}

// IncrementReceived increments the received message count.
func (pm *PeerMessages) IncrementReceived(name string, count int64) {
	pm.Received.IncrementByName(name, count)
	pm.PacketsReceived.IncrementByName(name, 1)
}

// IncrementSent increments the sent message count.
func (pm *PeerMessages) IncrementSent(name string, count int64) {
	pm.Sent.IncrementByName(name, count)
	pm.PacketsSent.IncrementByName(name, 1)
}

// Load returns a snapshot of all message counts.
func (pm *PeerMessages) Load() PeerMessages {
	return PeerMessages{
		Received:        pm.Received.Load(),
		Sent:            pm.Sent.Load(),
		PacketsReceived: pm.PacketsReceived.Load(),
		PacketsSent:     pm.PacketsSent.Load(),
	}
}
