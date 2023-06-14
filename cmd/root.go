/*
Copyright © 2022 Polygon <engineering@polygon.technology>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Lesser General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Lesser General Public License for more details.

You should have received a copy of the GNU Lesser General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/maticnetwork/polygon-cli/cmd/fork"
	"github.com/maticnetwork/polygon-cli/cmd/p2p"
	"github.com/maticnetwork/polygon-cli/cmd/parseethwallet"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/maticnetwork/polygon-cli/cmd/abi"
	"github.com/maticnetwork/polygon-cli/cmd/dumpblocks"
	"github.com/maticnetwork/polygon-cli/cmd/forge"
	"github.com/maticnetwork/polygon-cli/cmd/hash"
	"github.com/maticnetwork/polygon-cli/cmd/loadtest"
	"github.com/maticnetwork/polygon-cli/cmd/metricsToDash"
	"github.com/maticnetwork/polygon-cli/cmd/mnemonic"
	"github.com/maticnetwork/polygon-cli/cmd/monitor"
	"github.com/maticnetwork/polygon-cli/cmd/nodekey"
	"github.com/maticnetwork/polygon-cli/cmd/rpc"
	"github.com/maticnetwork/polygon-cli/cmd/version"
	"github.com/maticnetwork/polygon-cli/cmd/wallet"
)

var (
	cfgFile   string
	verbosity int
	pretty    bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "polycli",
	Short: "A Swiss Army knife of blockchain tools",
	Long: `Polycli is a collection of tools that are meant to be useful while
building, testing, and running block chain applications.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		setLogLevel(verbosity, pretty)
	},
}

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

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.polygon-cli.yaml)")
	rootCmd.PersistentFlags().IntVarP(&verbosity, "verbosity", "v", 200, "0 - Silent\n100 Fatals\n200 Errors\n300 Warnings\n400 INFO\n500 Debug\n600 Trace")
	rootCmd.PersistentFlags().BoolVar(&pretty, "pretty-logs", true, "Should we log in pretty format or JSON")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.SetOut(os.Stdout)

	rootCmd.AddCommand(dumpblocks.DumpblocksCmd)
	rootCmd.AddCommand(forge.ForgeCmd)
	rootCmd.AddCommand(hash.HashCmd)
	rootCmd.AddCommand(loadtest.LoadtestCmd)
	rootCmd.AddCommand(metricsToDash.MetricsToDashCmd)
	rootCmd.AddCommand(monitor.MonitorCmd)
	rootCmd.AddCommand(mnemonic.MnemonicCmd)
	rootCmd.AddCommand(nodekey.NodekeyCmd)
	rootCmd.AddCommand(rpc.RpcCmd)
	rootCmd.AddCommand(abi.ABICmd)
	rootCmd.AddCommand(version.VersionCmd)
	rootCmd.AddCommand(wallet.WalletCmd)
	rootCmd.AddCommand(fork.ForkCmd)
	rootCmd.AddCommand(parseethwallet.ParseETHWalletCmd)
	rootCmd.AddCommand(p2p.P2pCmd)
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

// setLogLevel sets the log level based on the flags.
func setLogLevel(verbosity int, pretty bool) {
	if verbosity < 100 {
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	} else if verbosity < 200 {
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	} else if verbosity < 300 {
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	} else if verbosity < 400 {
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	} else if verbosity < 500 {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else if verbosity < 600 {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	}

	if pretty {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		log.Debug().Msg("Starting logger in console mode")
	} else {
		log.Debug().Msg("Starting logger in json mode")
	}
}
