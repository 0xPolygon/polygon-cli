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

type NodeJSON struct {
	URL    string  `json:"url"`
	Hello  *Hello  `json:"hello,omitempty"`
	Status *Status `json:"status,omitempty"`
	Error  string  `json:"error,omitempty"`
	Time   int64   `json:"time,omitempty"`
	Nodes  int     `json:"nodes,omitempty"`
}

type NodeSet map[enode.ID][]NodeJSON

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

func WriteNodeSet(file string, ns NodeSet, writeErrors bool) error {
	if !writeErrors {
		keys := []enode.ID{}

		for id, events := range ns {
			onlyErrors := true

			for _, event := range events {
				if len(event.Error) == 0 {
					onlyErrors = false
				}
			}

			if onlyErrors {
				keys = append(keys, id)
			}
		}

		for _, key := range keys {
			delete(ns, key)
		}
	}

	bytes, err := json.MarshalIndent(ns, "", jsonIndent)
	if err != nil {
		return err
	}

	if len(file) == 0 {
		_, err = os.Stdout.Write(bytes)
		return err
	}

	return os.WriteFile(file, bytes, 0644)
}

func WriteURLs(file string, ns NodeSet) error {
	m := make(map[string]struct{})
	for _, events := range ns {
		for _, event := range events {
			if len(event.Error) > 0 {
				continue
			}

			m[event.URL] = struct{}{}
		}
	}

	urls := []string{}
	for url := range m {
		urls = append(urls, url)
	}

	bytes, err := json.MarshalIndent(urls, "", jsonIndent)
	if err != nil {
		return err
	}

	if len(file) == 0 {
		_, err = os.Stdout.Write(bytes)
		return err
	}

	return os.WriteFile(file, bytes, 0644)
}

func WritePeers(file string, nodes map[enode.ID]string) error {
	urls := []string{}
	for _, node := range nodes {
		urls = append(urls, node)
	}

	bytes, err := json.MarshalIndent(urls, "", jsonIndent)
	if err != nil {
		return err
	}

	if len(file) == 0 {
		_, err = os.Stdout.Write(bytes)
		return err
	}

	return os.WriteFile(file, bytes, 0644)
}
