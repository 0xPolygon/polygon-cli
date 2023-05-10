package p2p

import (
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/maticnetwork/polygon-cli/cmd/p2p/client"
	"github.com/maticnetwork/polygon-cli/cmd/p2p/crawl"
	"github.com/maticnetwork/polygon-cli/cmd/p2p/ping"
)

var verbosity int

var P2pCmd = &cobra.Command{
	Use:   "p2p",
	Short: "Commands related to devp2p",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		setMonitorLogLevel(verbosity)
	},
}

func init() {
	P2pCmd.PersistentFlags().IntVarP(&verbosity, "verbosity", "v", 400, "0 - Silent\n100 Fatal\n200 Error\n300 Warning\n400 Info\n500 Debug\n600 Trace")

	P2pCmd.AddCommand(client.ClientCmd)
	P2pCmd.AddCommand(crawl.CrawlCmd)
	P2pCmd.AddCommand(ping.PingCmd)
}

// setMonitorLogLevel sets the log level based on the flags.
func setMonitorLogLevel(verbosity int) {
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
}
