package p2p

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/rs/zerolog/log"
)

const jsonIndent = "    "

// NodeSet is the mapping of the node ID to the URL.
type NodeSet map[enode.ID]string

// ReadNodeSet parses a list of discovery node URLs loaded from a JSON file.
func ReadNodeSet(file string) ([]*enode.Node, error) {
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

// WriteNodeSet writes the node set as a JSON list of URLs to a file.
func WriteNodeSet(file string, nodes NodeSet) error {
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
