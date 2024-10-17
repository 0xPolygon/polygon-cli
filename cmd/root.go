package cmd

import (
	"fmt"
	"github.com/0xPolygon/polygon-cli/cmd/retest"
	"github.com/0xPolygon/polygon-cli/cmd/ulxly"
	"os"

	"github.com/0xPolygon/polygon-cli/cmd/fork"
	"github.com/0xPolygon/polygon-cli/cmd/p2p"
	"github.com/0xPolygon/polygon-cli/cmd/parseethwallet"
	"github.com/0xPolygon/polygon-cli/util"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/0xPolygon/polygon-cli/cmd/abi"
	"github.com/0xPolygon/polygon-cli/cmd/dbbench"
	"github.com/0xPolygon/polygon-cli/cmd/dumpblocks"
	"github.com/0xPolygon/polygon-cli/cmd/ecrecover"
	"github.com/0xPolygon/polygon-cli/cmd/enr"
	"github.com/0xPolygon/polygon-cli/cmd/fund"
	"github.com/0xPolygon/polygon-cli/cmd/hash"
	"github.com/0xPolygon/polygon-cli/cmd/loadtest"
	"github.com/0xPolygon/polygon-cli/cmd/metricsToDash"
	"github.com/0xPolygon/polygon-cli/cmd/mnemonic"
	"github.com/0xPolygon/polygon-cli/cmd/monitor"
	"github.com/0xPolygon/polygon-cli/cmd/nodekey"
	"github.com/0xPolygon/polygon-cli/cmd/rpcfuzz"
	"github.com/0xPolygon/polygon-cli/cmd/signer"
	"github.com/0xPolygon/polygon-cli/cmd/version"
	"github.com/0xPolygon/polygon-cli/cmd/wallet"
)

var (
	cfgFile   string
	verbosity int
	pretty    bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd *cobra.Command

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd = NewPolycliCommand()
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".polygon-cli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".polygon-cli")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

// NewPolycliCommand creates the `polycli` command.
func NewPolycliCommand() *cobra.Command {
	// Parent command to which all subcommands are added.
	cmd := &cobra.Command{
		Use:   "polycli",
		Short: "A Swiss Army knife of blockchain tools.",
		Long:  "Polycli is a collection of tools that are meant to be useful while building, testing, and running block chain applications.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			util.SetLogLevel(verbosity)
			logMode := util.JSON
			if pretty {
				logMode = util.Console
			}
			return util.SetLogMode(logMode)
		},
	}

	// Define flags and configuration settings.
	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.polygon-cli.yaml)")
	cmd.PersistentFlags().IntVarP(&verbosity, "verbosity", "v", 500, "0 - Silent\n100 Panic\n200 Fatal\n300 Error\n400 Warning\n500 Info\n600 Debug\n700 Trace")
	cmd.PersistentFlags().BoolVar(&pretty, "pretty-logs", true, "Should logs be in pretty format or JSON")

	// Define local flags which will only run when this action is called directly.
	cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	cmd.SetOut(os.Stdout)

	// Define commands.
	cmd.AddCommand(
		abi.ABICmd,
		dumpblocks.DumpblocksCmd,
		ecrecover.EcRecoverCmd,
		fork.ForkCmd,
		fund.FundCmd,
		hash.HashCmd,
		enr.ENRCmd,
		dbbench.DBBenchCmd,
		loadtest.LoadtestCmd,
		metricsToDash.MetricsToDashCmd,
		mnemonic.MnemonicCmd,
		monitor.MonitorCmd,
		nodekey.NodekeyCmd,
		p2p.P2pCmd,
		parseethwallet.ParseETHWalletCmd,
		retest.RetestCmd,
		rpcfuzz.RPCFuzzCmd,
		signer.SignerCmd,
		ulxly.ULxLyCmd,
		version.VersionCmd,
		wallet.WalletCmd,
	)
	return cmd
}
