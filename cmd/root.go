package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/0xPolygon/polygon-cli/cmd/abi"
	"github.com/0xPolygon/polygon-cli/cmd/cdk"
	"github.com/0xPolygon/polygon-cli/cmd/contract"
	"github.com/0xPolygon/polygon-cli/cmd/dbbench"
	"github.com/0xPolygon/polygon-cli/cmd/dockerlogger"
	"github.com/0xPolygon/polygon-cli/cmd/dumpblocks"
	"github.com/0xPolygon/polygon-cli/cmd/ecrecover"
	"github.com/0xPolygon/polygon-cli/cmd/enr"
	"github.com/0xPolygon/polygon-cli/cmd/fixnoncegap"
	"github.com/0xPolygon/polygon-cli/cmd/foldtrace"
	"github.com/0xPolygon/polygon-cli/cmd/fork"
	"github.com/0xPolygon/polygon-cli/cmd/fund"
	"github.com/0xPolygon/polygon-cli/cmd/hash"
	"github.com/0xPolygon/polygon-cli/cmd/loadtest"
	"github.com/0xPolygon/polygon-cli/cmd/metricstodash"
	"github.com/0xPolygon/polygon-cli/cmd/mnemonic"
	"github.com/0xPolygon/polygon-cli/cmd/monitor"
	"github.com/0xPolygon/polygon-cli/cmd/monitorv2"
	"github.com/0xPolygon/polygon-cli/cmd/nodekey"
	"github.com/0xPolygon/polygon-cli/cmd/p2p"
	"github.com/0xPolygon/polygon-cli/cmd/parsebatchl2data"
	"github.com/0xPolygon/polygon-cli/cmd/parseethwallet"
	"github.com/0xPolygon/polygon-cli/cmd/publish"
	"github.com/0xPolygon/polygon-cli/cmd/report"
	"github.com/0xPolygon/polygon-cli/cmd/retest"
	"github.com/0xPolygon/polygon-cli/cmd/rpcfuzz"
	"github.com/0xPolygon/polygon-cli/cmd/signer"
	"github.com/0xPolygon/polygon-cli/cmd/ulxly"
	"github.com/0xPolygon/polygon-cli/cmd/version"
	"github.com/0xPolygon/polygon-cli/cmd/wallet"
	"github.com/0xPolygon/polygon-cli/cmd/wrapcontract"
	"github.com/0xPolygon/polygon-cli/util"
)

var (
	cfgFile        string
	verbosityInput string
	pretty         bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd *cobra.Command

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.EnableTraverseRunHooks = true
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
		Long:  "Polycli is a collection of tools that are meant to be useful while building, testing, and running blockchain applications.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			verbosity, err := util.ParseVerbosity(verbosityInput)
			if err != nil {
				return err
			}
			util.SetLogLevel(verbosity)
			logMode := util.JSON
			if pretty {
				logMode = util.Console
			}
			return util.SetLogMode(logMode)
		},
	}

	// Define flags and configuration settings.
	f := cmd.PersistentFlags()
	f.StringVar(&cfgFile, "config", "", "config file (default is $HOME/.polygon-cli.yaml)")
	f.StringVarP(&verbosityInput, "verbosity", "v", "info", `log level (string or int):
  0   - silent
  100 - panic
  200 - fatal
  300 - error
  400 - warn
  500 - info (default)
  600 - debug
  700 - trace`)
	f.BoolVar(&pretty, "pretty-logs", true, "output logs in pretty format instead of JSON")

	// Define local flags which will only run when this action is called directly.
	cmd.Flags().BoolP("toggle", "t", false, "help message for toggle")
	cmd.SetOut(os.Stdout)

	// Define commands.
	cmd.AddCommand(
		abi.ABICmd,
		cdk.CDKCmd,
		dbbench.DBBenchCmd,
		dumpblocks.DumpblocksCmd,
		ecrecover.EcRecoverCmd,
		enr.ENRCmd,
		fixnoncegap.FixNonceGapCmd,
		fork.ForkCmd,
		fund.FundCmd,
		hash.HashCmd,
		loadtest.LoadtestCmd,
		metricstodash.MetricsToDashCmd,
		mnemonic.MnemonicCmd,
		monitor.MonitorCmd,
		monitorv2.MonitorV2Cmd,
		nodekey.NodekeyCmd,
		p2p.P2pCmd,
		parseethwallet.ParseETHWalletCmd,
		report.ReportCmd,
		retest.RetestCmd,
		rpcfuzz.RPCFuzzCmd,
		signer.SignerCmd,
		ulxly.ULxLyCmd,
		version.VersionCmd,
		wallet.WalletCmd,
		wrapcontract.WrapContractCmd,
		foldtrace.FoldTraceCmd,
		parsebatchl2data.ParseBatchL2Data,
		publish.Cmd,
		dockerlogger.Cmd,
		contract.Cmd,
	)

	return cmd
}
