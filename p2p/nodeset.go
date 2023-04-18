// Copyright 2019 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

package p2p

import (
	"bytes"
	"encoding/json"
	"os"
	"sort"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/p2p/enode"
)

const jsonIndent = "    "

// NodeSet is the nodes.json file format. It holds a set of node records
// as a JSON object.
type NodeSet map[enode.ID]NodeJSON

type NodeJSON struct {
	Seq uint64      `json:"seq"`
	N   *enode.Node `json:"record"`

	// The score tracks how many liveness checks were performed. It is incremented by one
	// every time the node passes a check, and halved every time it doesn't.
	Score int `json:"score,omitempty"`
	// These two track the time of last successful contact.
	FirstResponse time.Time `json:"firstResponse,omitempty"`
	LastResponse  time.Time `json:"lastResponse,omitempty"`
	// This one tracks the time of our last attempt to contact the node.
	LastCheck time.Time `json:"lastCheck,omitempty"`
}

func LoadNodesJSON(file string) (NodeSet, error) {
	var nodes NodeSet
	if err := common.LoadJSON(file, &nodes); err != nil {
		return nil, err
	}
	return nodes, nil
}

func WriteNodesJSON(file string, nodes NodeSet) error {
	nodesJSON, err := json.MarshalIndent(nodes, "", jsonIndent)
	if err != nil {
		return err
	}
	if file == "-" {
		_, err = os.Stdout.Write(nodesJSON)
		return err
	}
	return os.WriteFile(file, nodesJSON, 0644)
}

// Nodes returns the node records contained in the set.
func (ns NodeSet) Nodes() []*enode.Node {
	result := make([]*enode.Node, 0, len(ns))
	for _, n := range ns {
		result = append(result, n.N)
	}
	// Sort by ID.
	sort.Slice(result, func(i, j int) bool {
		return bytes.Compare(result[i].ID().Bytes(), result[j].ID().Bytes()) < 0
	})
	return result
}
