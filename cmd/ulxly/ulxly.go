package ulxly

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"io"
	"os"
	"strings"

	// note - this won't deal with the complexity of handling deposits prior to the ulxly
	"github.com/maticnetwork/polygon-cli/bindings/ulxly"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const (
	TreeDepth = 32
)

type uLxLyArgs struct {
	FromBlock     *uint64
	ToBlock       *uint64
	RPCURL        *string
	BridgeAddress *string
	FilterSize    *uint64

	InputFileName *string
	DepositNum    *uint32
}

type SMT struct {
	Data       map[uint32][]common.Hash
	Height     uint8
	Count      uint64
	Branches   map[uint32][][TreeDepth]byte
	Root       [TreeDepth]byte
	ZeroHashes [][TreeDepth]byte
	Proofs     map[uint32]Proof
}
type Proof struct {
	Siblings     [TreeDepth]common.Hash
	Root         common.Hash
	DepositCount uint32
}

var ulxlyInputArgs uLxLyArgs

var ULxLyCmd = &cobra.Command{
	Use:   "ulxly",
	Short: "Utilities for interacting with the lxly bridge",
	Long:  "TODO",
	Args:  cobra.NoArgs,
}

// polycli ulxly deposits --bridge-address 0x528e26b25a34a4A5d0dbDa1d57D318153d2ED582 --rpc-url https://sepolia-rpc.invalid --from-block 4880876 --to-block 6015235 --filter-size 9999
// polycli ulxly deposits --bridge-address 0x528e26b25a34a4A5d0dbDa1d57D318153d2ED582 --rpc-url https://sepolia-rpc.invalid --from-block 4880876 --to-block 6025854 --filter-size 999 > cardona-4880876-to-6025854.ndjson
var DepositsCmd = &cobra.Command{
	Use:     "deposits",
	Short:   "get a range of deposits",
	Long:    "TODO",
	Args:    cobra.NoArgs,
	PreRunE: checkDepositArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		// Dial the Ethereum RPC server.
		rpc, err := ethrpc.DialContext(ctx, *ulxlyInputArgs.RPCURL)
		if err != nil {
			log.Error().Err(err).Msg("Unable to dial rpc")
			return err
		}
		defer rpc.Close()
		ec := ethclient.NewClient(rpc)

		bridgeV2, err := ulxly.NewUlxly(common.HexToAddress(*ulxlyInputArgs.BridgeAddress), ec)
		if err != nil {
			return err
		}
		fromBlock := *ulxlyInputArgs.FromBlock
		toBlock := *ulxlyInputArgs.ToBlock
		currentBlock := fromBlock
		for currentBlock < toBlock {
			endBlock := currentBlock + *ulxlyInputArgs.FilterSize
			if endBlock > toBlock {
				endBlock = toBlock
			}

			opts := bind.FilterOpts{
				Start:   currentBlock,
				End:     &endBlock,
				Context: ctx,
			}
			evtV2Iterator, err := bridgeV2.FilterBridgeEvent(&opts)
			if err != nil {
				return err
			}

			for evtV2Iterator.Next() {
				evt := evtV2Iterator.Event
				log.Info().Uint32("deposit", evt.DepositCount).Uint64("block-number", evt.Raw.BlockNumber).Msg("Found ulxly Deposit")
				jBytes, err := json.Marshal(evt)
				if err != nil {
					return err
				}
				fmt.Println(string(jBytes))
			}
			evtV2Iterator.Close()
			currentBlock = endBlock
		}

		return nil
	},
}
var ProofCmd = &cobra.Command{
	Use:     "proof",
	Short:   "generate a merkle proof",
	Long:    "TODO",
	Args:    cobra.NoArgs,
	PreRunE: checkProofArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		rawDepositData, err := getInputData(cmd, args)
		if err != nil {
			return err
		}
		readDeposits(rawDepositData)
		return nil
	},
}

func checkProofArgs(cmd *cobra.Command, args []string) error {
	return nil
}
func getInputData(cmd *cobra.Command, args []string) ([]byte, error) {
	if ulxlyInputArgs.InputFileName != nil && *ulxlyInputArgs.InputFileName != "" {
		return os.ReadFile(*ulxlyInputArgs.InputFileName)
	}

	if len(args) > 1 {
		concat := strings.Join(args[1:], " ")
		return []byte(concat), nil
	}

	return io.ReadAll(os.Stdin)
}
func readDeposits(rawDeposits []byte) error {
	buf := bytes.NewBuffer(rawDeposits)
	scanner := bufio.NewScanner(buf)
	smt := new(SMT)
	smt.Init()
	seenDeposit := make(map[uint32]common.Hash, 0)
	lastDeposit := uint32(0)
	for scanner.Scan() {
		evt := new(ulxly.UlxlyBridgeEvent)
		err := json.Unmarshal(scanner.Bytes(), evt)
		if err != nil {
			return err
		}
		if _, hasBeenSeen := seenDeposit[evt.DepositCount]; hasBeenSeen {
			log.Warn().Uint32("deposit", evt.DepositCount).Str("tx-hash", evt.Raw.TxHash.String()).Msg("Skipping duplicate deposit")
			continue
		}
		seenDeposit[evt.DepositCount] = evt.Raw.TxHash
		if lastDeposit+1 != evt.DepositCount && lastDeposit != 0 {
			log.Error().Uint32("missing-deposit", lastDeposit+1).Uint32("current-deposit", evt.DepositCount).Msg("Missing deposit")
			return fmt.Errorf("missing deposit: %d", lastDeposit+1)
		}
		lastDeposit = evt.DepositCount
		smt.AddLeaf(evt)
		log.Info().
			Uint64("block-number", evt.Raw.BlockNumber).
			Uint32("deposit-count", evt.DepositCount).
			Str("tx-hash", evt.Raw.TxHash.String()).
			Str("root", common.Hash(smt.Root).String()).
			Msg("adding event to tree")

	}

	p := smt.Proofs[*ulxlyInputArgs.DepositNum]

	fmt.Println(p.String())
	return nil
}

func (p *Proof) String() string {
	jsonBytes, err := json.Marshal(p)
	if err != nil {
		log.Error().Err(err).Msg("error marshalling proof to json")
		return ""
	}
	return string(jsonBytes)

}

// This implementation looks good. We get this hash
// 0xf8c64768317c96c6c3c0f72b5a2cd3d03e831c200bf6bf0ab4d181877d59ddab
// for this deposit
// https://sepolia.etherscan.io/tx/0xf2003cf43a205bc777bc2d22fcb05b69ebb23464b39250d164cf9b09150b7833#eventlog
// And that seems to match a call to `getLeafValue`
func hashDeposit(deposit *ulxly.UlxlyBridgeEvent) common.Hash {
	var res [TreeDepth]byte
	origNet := make([]byte, 4) //nolint:gomnd
	binary.BigEndian.PutUint32(origNet, deposit.OriginNetwork)
	destNet := make([]byte, 4) //nolint:gomnd
	binary.BigEndian.PutUint32(destNet, deposit.DestinationNetwork)
	var buf [TreeDepth]byte
	metaHash := crypto.Keccak256Hash(deposit.Metadata)
	copy(res[:], crypto.Keccak256Hash([]byte{deposit.LeafType}, origNet, deposit.OriginAddress.Bytes(), destNet, deposit.DestinationAddress[:], deposit.Amount.FillBytes(buf[:]), metaHash.Bytes()).Bytes())
	return res
}

func (s *SMT) Init() {
	s.Branches = make(map[uint32][][TreeDepth]byte)
	s.Height = TreeDepth
	s.Data = make(map[uint32][]common.Hash, 0)
	s.ZeroHashes = generateZeroHashes(TreeDepth)
	s.Proofs = make(map[uint32]Proof)
}

// cast call --rpc-url https://eth-sepolia.g.alchemy.com/v2/demo --block 4880875 0xad1490c248c5d3cbae399fd529b79b42984277df 'lastMainnetExitRoot()(bytes32)'
// cast call --rpc-url https://eth-sepolia.g.alchemy.com/v2/demo --block 4880876 0xad1490c248c5d3cbae399fd529b79b42984277df 'lastMainnetExitRoot()(bytes32)'
// The first mainnet exit root for cardona is
// 0x112b077c64ed4a22dfb0ab3c2622d6ddbf3a5423afeb05878c2c21c4cb5e65da
func (s *SMT) AddLeaf(deposit *ulxly.UlxlyBridgeEvent) {
	leaf := hashDeposit(deposit)
	log.Info().Str("leaf-hash", common.Bytes2Hex(leaf[:])).Msg("Leaf hash calculated")

	node := leaf
	s.Count = uint64(deposit.DepositCount) + 1
	size := s.Count
	branches := make([][TreeDepth]byte, TreeDepth)
	if deposit.DepositCount == 0 {
		branches = generateZeroHashes(TreeDepth)
	} else {
		copy(branches, s.Branches[deposit.DepositCount-1])
	}

	for height := uint64(0); height < TreeDepth; height += 1 {
		if isBitSet(size, height) {
			copy(branches[height][:], node[:])
			break
		}
		node = crypto.Keccak256Hash(branches[height][:], node[:])
	}
	s.Branches[deposit.DepositCount] = branches
	s.Root = s.GetRoot(deposit.DepositCount)
}

func (s *SMT) GetRoot(depositNum uint32) common.Hash {
	var node common.Hash
	size := s.Count
	var zeroHashes = s.ZeroHashes

	siblings := [TreeDepth]common.Hash{}
	for height := 0; height < TreeDepth; height++ {
		currentZeroHashHeight := zeroHashes[height]
		left := crypto.Keccak256Hash(s.Branches[depositNum][height][:], node.Bytes())
		right := crypto.Keccak256Hash(node.Bytes(), currentZeroHashHeight[:])
		if depositNum == 24391 {
			log.Debug().
				Int("height", height).
				Str("sib-1", node.String()).
				Str("sib-2", common.Hash(s.Branches[depositNum][height]).String()).
				Str("sib-3", common.Hash(currentZeroHashHeight).String()).
				Str("left", left.String()).
				Str("right", right.String()).
				Msg("tree values")
		}

		if ((size >> height) & 1) == 1 {
			copy(siblings[height][:], currentZeroHashHeight[:])
			node = left
		} else {
			copy(siblings[height][:], s.Branches[depositNum][height][:])
			node = right
		}
	}
	if depositNum^1 == 1 {
		copy(siblings[0][:], zeroHashes[0][:])
	} else {
		copy(siblings[0][:], s.Branches[depositNum-1][0][:])
	}

	s.Proofs[depositNum] = Proof{
		Siblings:     siblings,
		Root:         node,
		DepositCount: depositNum,
	}
	log.Info().Str("inner-root", s.Proofs[depositNum].Root.String()).Str("root", node.String()).Msg("root Check")
	return node
}

// GetProof will try to generate the paired hashes for each level of the tree for a given deposit
func (s *SMT) GetProof(depositNum uint32) ([][TreeDepth]byte, error) {
	// it's possible to call this with a deposit that doesn't exist. In theory this should be fine, but it doesn't really fit the intent of the code here
	if uint64(depositNum) > s.Count {
		return nil, fmt.Errorf("deposit number %d is greater than the count %d of leaves", depositNum, s.Count)
	}

	// At the bottom of the path, we need to find the deposit that is paired with the despot that we're checking. We should be able to flip the last bit and get the deposit number
	siblingDepositNum := depositNum ^ 1

	var siblingLeaf [TreeDepth]byte
	if _, hasKey := s.Branches[siblingDepositNum]; !hasKey || siblingDepositNum > depositNum {
		siblingLeaf = [TreeDepth]byte{}
	} else {
		siblingLeaf = s.Branches[siblingDepositNum][0]
	}
	siblings := [TreeDepth][TreeDepth]byte{}
	siblings[0] = siblingLeaf
	siblingMask := uint32(1)
	// Iterate through the levels of the tree
	for height := 1; height < TreeDepth; height += 1 {
		// At each height, we're going to flip a bit in the deposit number in order to find the complimentary node for the given height
		siblingMask = uint32(siblingMask<<1) + 1
		flipMask := uint32(1 << height)
		currentSiblingDeposit := (depositNum | siblingMask) ^ flipMask

		// we'll get the branches from the particular deposit that we're interested in
		b, hasKey := s.Branches[currentSiblingDeposit]
		if !hasKey || isEmpty(b[height]) || currentSiblingDeposit > depositNum {
			b = generateZeroHashes(TreeDepth)
		}
		siblings[height] = b[height]
		log.Info().
			Int("height", height).
			Uint32("current-sibling-num", currentSiblingDeposit).
			Str("current-sibling", fmt.Sprintf("%b", currentSiblingDeposit)).
			Str("current-mask", fmt.Sprintf("%b", siblingMask)).
			Str("sib", common.Hash(b[height]).String()).
			Msg("Getting Deposit")

	}
	return siblings[:], nil
}

func isEmpty(h [TreeDepth]byte) bool {
	for i := 0; i < TreeDepth; i = i + 1 {
		if h[i] != 0 {
			return false
		}
	}
	return true
}

// https://eth2book.info/capella/part2/deposits-withdrawals/contract/
func generateZeroHashes(height uint8) [][TreeDepth]byte {
	var zeroHashes = [][TreeDepth]byte{}
	zeroHashes = append(zeroHashes, common.Hash{})
	for i := 1; i <= int(height); i++ {
		zeroHashes = append(zeroHashes, crypto.Keccak256Hash(zeroHashes[i-1][:], zeroHashes[i-1][:]))
	}
	return zeroHashes
}

// isBitSet checks if the bit at position h in number dc is set to 1
func isBitSet(dc uint64, h uint64) bool {
	mask := uint64(1)
	for i := uint64(0); i < h; i++ {
		mask *= 2
	}
	return dc&mask > 0
}
func checkDepositArgs(cmd *cobra.Command, args []string) error {
	if *ulxlyInputArgs.BridgeAddress == "" {
		return fmt.Errorf("please provide the bridge address")
	}
	if *ulxlyInputArgs.FromBlock > *ulxlyInputArgs.ToBlock {
		return fmt.Errorf("the from block should be less than the to block")
	}
	return nil
}

func init() {
	ULxLyCmd.AddCommand(DepositsCmd)
	ULxLyCmd.AddCommand(ProofCmd)
	ulxlyInputArgs.FromBlock = DepositsCmd.PersistentFlags().Uint64("from-block", 0, "The block height to start query at.")
	ulxlyInputArgs.ToBlock = DepositsCmd.PersistentFlags().Uint64("to-block", 0, "The block height to start query at.")
	ulxlyInputArgs.RPCURL = DepositsCmd.PersistentFlags().String("rpc-url", "http://127.0.0.1:8545", "The RPC to query for events")
	ulxlyInputArgs.FilterSize = DepositsCmd.PersistentFlags().Uint64("filter-size", 1000, "The batch size for individual filter queries")

	ulxlyInputArgs.BridgeAddress = DepositsCmd.Flags().String("bridge-address", "", "The address of the lxly bridge")
	ulxlyInputArgs.InputFileName = ProofCmd.PersistentFlags().String("file-name", "", "The filename with ndjson data of deposits")
	ulxlyInputArgs.DepositNum = ProofCmd.PersistentFlags().Uint32("deposit-number", 0, "The deposit that we would like to prove")
}
