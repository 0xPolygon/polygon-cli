package rpctypes

import (
	"encoding/json"
	"testing"
)

func TestPrettyMarshaling(t *testing.T) {
	// Create a sample block response
	rawBlock := &RawBlockResponse{
		Number:           "0x1234567",
		Hash:             "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
		ParentHash:       "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		Nonce:            "0x42",
		SHA3Uncles:       "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
		LogsBloom:        "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		TransactionsRoot: "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
		StateRoot:        "0xd7f8974fb5ac78d9ac099b9ad5018bedc2ce0a72dad1827a1709da30580f0544",
		ReceiptsRoot:     "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
		Miner:            "0x0000000000000000000000000000000000000000",
		Difficulty:       "0x400000000",
		ExtraData:        "0x",
		Size:             "0x220",
		GasLimit:         "0x1c9c380",
		GasUsed:          "0x5208",
		Timestamp:        "0x55ba467c",
		Transactions:     []RawTransactionResponse{},
		Uncles:           []RawData32Response{},
		BaseFeePerGas:    "0x7",
		MixHash:          "0x0000000000000000000000000000000000000000000000000000000000000000",
	}

	// Create a PolyBlock
	block := NewPolyBlock(rawBlock)

	// Test regular JSON marshaling
	regularJSON, err := block.MarshalJSON()
	if err != nil {
		t.Fatalf("Failed to marshal regular JSON: %v", err)
	}

	// Test pretty JSON marshaling
	prettyJSON, err := PolyBlockToPrettyJSON(block)
	if err != nil {
		t.Fatalf("Failed to marshal pretty JSON: %v", err)
	}

	// Verify both are valid JSON
	var regularData map[string]any
	if err := json.Unmarshal(regularJSON, &regularData); err != nil {
		t.Fatalf("Regular JSON is invalid: %v", err)
	}

	var prettyData map[string]any
	if err := json.Unmarshal(prettyJSON, &prettyData); err != nil {
		t.Fatalf("Pretty JSON is invalid: %v", err)
	}

	// Verify that timestamp is different format
	regularTimestamp := regularData["timestamp"].(string)
	prettyTimestamp := prettyData["timestamp"].(float64)

	t.Logf("Regular timestamp (hex): %s", regularTimestamp)
	t.Logf("Pretty timestamp (uint64): %f", prettyTimestamp)

	// Verify that the hex timestamp matches the converted value
	if regularTimestamp != "0x55ba467c" {
		t.Errorf("Expected regular timestamp to be 0x55ba467c, got %s", regularTimestamp)
	}

	expectedTimestamp := float64(0x55ba467c)
	if prettyTimestamp != expectedTimestamp {
		t.Errorf("Expected pretty timestamp to be %f, got %f", expectedTimestamp, prettyTimestamp)
	}
}

func TestPrettyTransactionMarshaling(t *testing.T) {
	// Create a sample transaction response
	rawTx := &RawTransactionResponse{
		BlockHash:        "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
		BlockNumber:      "0x1234567",
		From:             "0x1234567890123456789012345678901234567890",
		Gas:              "0x5208",
		GasPrice:         "0x4a817c800",
		Hash:             "0xfedcba0987654321fedcba0987654321fedcba0987654321fedcba0987654321",
		Input:            "0x",
		Nonce:            "0x42",
		To:               "0x0987654321098765432109876543210987654321",
		TransactionIndex: "0x0",
		Value:            "0xde0b6b3a7640000",
		V:                "0x1c",
		R:                "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		S:                "0xfedcba0987654321fedcba0987654321fedcba0987654321fedcba0987654321",
		Type:             "0x0",
		ChainID:          "0x1",
		AccessList:       []any{},
	}

	// Create a PolyTransaction
	tx := NewPolyTransaction(rawTx)

	// Test pretty JSON marshaling
	prettyJSON, err := PolyTransactionToPrettyJSON(tx)
	if err != nil {
		t.Fatalf("Failed to marshal pretty JSON: %v", err)
	}

	// Verify it's valid JSON
	var prettyData map[string]any
	if err := json.Unmarshal(prettyJSON, &prettyData); err != nil {
		t.Fatalf("Pretty JSON is invalid: %v", err)
	}

	// Verify that gas is now a number instead of hex string
	prettyGas := prettyData["gas"].(float64)
	expectedGas := float64(0x5208)
	if prettyGas != expectedGas {
		t.Errorf("Expected pretty gas to be %f, got %f", expectedGas, prettyGas)
	}

	t.Logf("Pretty gas (uint64): %f", prettyGas)
}

func TestPrettyReceiptLogMarshaling(t *testing.T) {
	// Create a sample receipt with logs
	rawReceipt := &RawTxReceipt{
		TransactionHash:   "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
		TransactionIndex:  "0x1",
		BlockHash:         "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
		BlockNumber:       "0x12345",
		From:              "0x1234567890123456789012345678901234567890",
		To:                "0x0987654321098765432109876543210987654321",
		CumulativeGasUsed: "0x5208",
		EffectiveGasPrice: "0x4a817c800",
		GasUsed:           "0x5208",
		ContractAddress:   "0x0000000000000000000000000000000000000000",
		Logs: []RawTxLogs{
			{
				BlockHash:        "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
				BlockNumber:      "0x12345",
				TransactionIndex: "0x1",
				Address:          "0x0987654321098765432109876543210987654321",
				LogIndex:         "0x0",
				Data:             "0x0000000000000000000000000000000000000000000000000de0b6b3a7640000",
				Removed:          false,
				Topics: []RawData32Response{
					"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef", // Transfer event signature
					"0x0000000000000000000000001234567890123456789012345678901234567890", // from address
					"0x0000000000000000000000000987654321098765432109876543210987654321", // to address
				},
				TransactionHash: "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
			},
		},
		LogsBloom:    "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		Root:         "0x0000000000000000000000000000000000000000000000000000000000000000",
		Status:       "0x1",
		BlobGasPrice: "0x0",
		BlobGasUsed:  "0x0",
	}

	// Create a PolyReceipt
	receipt := NewPolyReceipt(rawReceipt)

	// Test regular JSON marshaling
	regularJSON, err := receipt.MarshalJSON()
	if err != nil {
		t.Fatalf("Failed to marshal regular JSON: %v", err)
	}

	// Test pretty JSON marshaling
	prettyJSON, err := PolyReceiptToPrettyJSON(receipt)
	if err != nil {
		t.Fatalf("Failed to marshal pretty JSON: %v", err)
	}

	// Verify both are valid JSON
	var regularData map[string]any
	if err := json.Unmarshal(regularJSON, &regularData); err != nil {
		t.Fatalf("Regular JSON is invalid: %v", err)
	}

	var prettyData map[string]any
	if err := json.Unmarshal(prettyJSON, &prettyData); err != nil {
		t.Fatalf("Pretty JSON is invalid: %v", err)
	}

	// Verify logs structure in regular vs pretty JSON
	regularLogs := regularData["logs"].([]any)
	prettyLogs := prettyData["logs"].([]any)

	if len(regularLogs) != len(prettyLogs) {
		t.Errorf("Expected same number of logs, got regular: %d, pretty: %d", len(regularLogs), len(prettyLogs))
	}

	if len(regularLogs) > 0 {
		regularLog := regularLogs[0].(map[string]any)
		prettyLog := prettyLogs[0].(map[string]any)

		// Check that blockNumber is converted from hex string to uint64
		regularBlockNumber := regularLog["blockNumber"].(string)
		prettyBlockNumber := prettyLog["blockNumber"].(float64)

		t.Logf("Regular log blockNumber (hex): %s", regularBlockNumber)
		t.Logf("Pretty log blockNumber (uint64): %f", prettyBlockNumber)

		if regularBlockNumber != "0x12345" {
			t.Errorf("Expected regular blockNumber to be 0x12345, got %s", regularBlockNumber)
		}

		expectedBlockNumber := float64(0x12345)
		if prettyBlockNumber != expectedBlockNumber {
			t.Errorf("Expected pretty blockNumber to be %f, got %f", expectedBlockNumber, prettyBlockNumber)
		}

		// Check that topics are converted properly
		regularTopics := regularLog["topics"].([]any)
		prettyTopics := prettyLog["topics"].([]any)

		t.Logf("Regular log topics[0] (hex): %s", regularTopics[0].(string))
		t.Logf("Pretty log topics[0] (hash): %s", prettyTopics[0].(string))

		// Both should be the same hash format but verify they're valid
		if len(regularTopics) != len(prettyTopics) {
			t.Errorf("Expected same number of topics, got regular: %d, pretty: %d", len(regularTopics), len(prettyTopics))
		}
	}
}
