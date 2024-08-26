package crawl

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p/discover"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/maticnetwork/polygon-cli/p2p"
)

type (
	crawlParams struct {
		Bootnodes            string
		Timeout              string
		timeout              time.Duration
		Threads              int
		NetworkID            uint64
		NodesFile            string
		Database             string
		RevalidationInterval string
		OnlyURLs             bool

		revalidationInterval time.Duration
	}
)

var (
	inputCrawlParams crawlParams
)

// crawlCmd represents the crawl command. This is responsible for crawling the
// devp2p layer and generating a nodes json file with peers.
var CrawlCmd = &cobra.Command{
	Use:   "crawl [nodes file]",
	Short: "Crawl a network on the devp2p layer and generate a nodes JSON file.",
	Long:  "If no nodes.json file exists, it will be created.",
	Args:  cobra.MinimumNArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		inputCrawlParams.NodesFile = args[0]

		inputCrawlParams.timeout, err = time.ParseDuration(inputCrawlParams.Timeout)
		if err != nil {
			return err
		}

		inputCrawlParams.revalidationInterval, err = time.ParseDuration(inputCrawlParams.RevalidationInterval)
		if err != nil {
			return err
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		nodes, err := p2p.ReadNodeSet(inputCrawlParams.NodesFile)
		if err != nil {
			log.Warn().Err(err).Msgf("Creating nodes file %v because it does not exist", inputCrawlParams.NodesFile)
		}

		var cfg discover.Config
		cfg.PrivateKey, _ = crypto.GenerateKey()
		cfg.Bootnodes, err = p2p.ParseBootnodes(inputCrawlParams.Bootnodes)
		if err != nil {
			return fmt.Errorf("unable to parse bootnodes: %w", err)
		}

		db, err := enode.OpenDB(inputCrawlParams.Database)
		if err != nil {
			return err
		}

		ln := enode.NewLocalNode(db, cfg.PrivateKey)
		socket, err := p2p.Listen(ln)
		if err != nil {
			return err
		}

		disc, err := discover.ListenV4(socket, ln, cfg)
		if err != nil {
			return err
		}
		defer disc.Close()

		c := newCrawler(nodes, disc, disc.RandomNodes())
		c.revalidateInterval = inputCrawlParams.revalidationInterval

		log.Info().Msg("Starting crawl")

		output := c.run(inputCrawlParams.timeout, inputCrawlParams.Threads)

		if inputCrawlParams.OnlyURLs {
			return p2p.WriteURLs(inputCrawlParams.NodesFile, output)
		}

		return p2p.WriteNodeSet(inputCrawlParams.NodesFile, output, false)
	},
}

func init() {
	CrawlCmd.PersistentFlags().StringVarP(&inputCrawlParams.Bootnodes, "bootnodes", "b", "",
		`Comma separated nodes used for bootstrapping. At least one bootnode is
required, so other nodes in the network can discover each other.`)
	if err := CrawlCmd.MarkPersistentFlagRequired("bootnodes"); err != nil {
		log.Error().Err(err).Msg("Failed to mark bootnodes as required persistent flag")
	}
	CrawlCmd.PersistentFlags().StringVarP(&inputCrawlParams.Timeout, "timeout", "t", "30m0s", "Time limit for the crawl")
	CrawlCmd.PersistentFlags().IntVarP(&inputCrawlParams.Threads, "parallel", "p", 16, "How many parallel discoveries to attempt")
	CrawlCmd.PersistentFlags().Uint64VarP(&inputCrawlParams.NetworkID, "network-id", "n", 0, "Filter discovered nodes by this network id")
	CrawlCmd.PersistentFlags().StringVarP(&inputCrawlParams.Database, "database", "d", "", "Node database for updating and storing client information")
	CrawlCmd.PersistentFlags().StringVarP(&inputCrawlParams.RevalidationInterval, "revalidation-interval", "r", "10m", "Time before retrying to connect to a failed peer")
	CrawlCmd.PersistentFlags().BoolVarP(&inputCrawlParams.OnlyURLs, "only-urls", "u", true, "Only writes the enode URLs to the output")
}
