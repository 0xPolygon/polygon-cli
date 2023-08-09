package p2p

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/rs/zerolog/log"
)

const jsonIndent = "    "

// NodeSet is the nodes.json file format. It holds a set of node records
// as a JSON object.
type NodeSet map[enode.ID]NodeJSON
type StaticNodes map[enode.ID]string

type NodeJSON struct {
	Seq uint64      `json:"seq"`
	N   *enode.Node `json:"record"`
	URL string      `json:"url"`

	// The score tracks how many liveness checks were performed. It is incremented by one
	// every time the node passes a check, and halved every time it doesn't.
	Score int `json:"score,omitempty"`
	// These two track the time of last successful contact.
	FirstResponse time.Time `json:"firstResponse,omitempty"`
	LastResponse  time.Time `json:"lastResponse,omitempty"`
	// This one tracks the time of our last attempt to contact the node.
	LastCheck time.Time `json:"lastCheck,omitempty"`
}

func ReadNodeSet(file string) (NodeSet, error) {
	var nodes NodeSet
	if err := common.LoadJSON(file, &nodes); err != nil {
		return nil, err
	}
	return nodes, nil
}

func WriteNodeSet(file string, nodes NodeSet) error {
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

// ReadStaticNodes parses a list of discovery node URLs loaded from a JSON file
// from within the data directory.
func ReadStaticNodes(file string) ([]*enode.Node, error) {
	// Load the nodes from the config file.
	var nodelist []string
	if err := common.LoadJSON(file, &nodelist); err != nil {
		return nil, fmt.Errorf("failed to load node list file: %w", err)
	}

	// Interpret the list as a discovery node array
	var nodes []*enode.Node
	for _, url := range nodelist {
		if url == "" {
			continue
		}
		node, err := enode.Parse(enode.ValidSchemes, url)
		if err != nil {
			log.Warn().Err(err).Str("url", url).Msg("Failed to parse enode")
			continue
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}

func WriteStaticNodes(file string, nodes StaticNodes) error {
	urls := make([]string, 0, len(nodes))
	for _, url := range nodes {
		urls = append(urls, url)
	}

	bytes, err := json.MarshalIndent(urls, "", jsonIndent)
	if err != nil {
		return err
	}
	if file == "-" {
		_, err = os.Stdout.Write(bytes)
		return err
	}
	return os.WriteFile(file, bytes, 0644)
}
