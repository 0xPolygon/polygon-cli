package metricstodash

import (
	"fmt"

	_ "embed"

	"github.com/0xPolygon/polygon-cli/dashboard"
	"github.com/spf13/cobra"
)

var (
	//go:embed usage.md
	usage                                 string
	inputMetricsToDashFile                string
	inputMetricsToDashPrefix              string
	inputMetricsToDashTitle               string
	inputMetricsToDashDesc                string
	inputMetricsToDashHeight              int
	inputMetricsToDashWidth               int
	inputMetricsToDashTemplateVars        []string
	inputMetricsToDashTemplateVarDefaults []string
	inputMetricsToDashStripPrefixes       []string
	inputMetricsToDashPretty              bool
	inputMetricsToDashShowHelp            bool
)

// MetricsToDashCmd represents the metricsToDash command
var MetricsToDashCmd = &cobra.Command{
	Use:     "metrics-to-dash",
	Aliases: []string{"metricstodash", "metricsToDash"},
	Short:   "Create a dashboard from an Openmetrics / Prometheus response.",
	Long:    usage,
	RunE: func(cmd *cobra.Command, args []string) error {
		do := dashboard.DashboardOptions{
			File:                inputMetricsToDashFile,
			Title:               inputMetricsToDashTitle,
			Prefix:              inputMetricsToDashPrefix,
			Description:         inputMetricsToDashDesc,
			WidgetHeight:        inputMetricsToDashHeight,
			WidgetWidth:         inputMetricsToDashWidth,
			TemplateVars:        inputMetricsToDashTemplateVars,
			TemplateVarDefaults: inputMetricsToDashTemplateVarDefaults,
			StripPrefixes:       inputMetricsToDashStripPrefixes,
			Pretty:              inputMetricsToDashPretty,
			ShowHelp:            inputMetricsToDashShowHelp,
		}
		data, err := dashboard.ConvertMetricsToDashboard(&do)
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	},
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	f := MetricsToDashCmd.Flags()
	f.StringVarP(&inputMetricsToDashFile, "input-file", "i", "", "the metrics file to be used")
	f.StringVarP(&inputMetricsToDashPrefix, "prefix", "p", "", "prefix to use before all metrics")
	f.StringVarP(&inputMetricsToDashTitle, "title", "t", "Polycli Dashboard", "title for the dashboard")
	f.StringVarP(&inputMetricsToDashDesc, "desc", "d", "Polycli Dashboard", "description for the dashboard")
	f.IntVarP(&inputMetricsToDashWidth, "width", "W", 4, "widget width")
	f.IntVarP(&inputMetricsToDashHeight, "height", "H", 3, "widget height")

	f.StringArrayVarP(&inputMetricsToDashTemplateVars, "template-vars", "T", []string{}, "the template variables to use for the dashboard")
	f.StringArrayVarP(&inputMetricsToDashTemplateVarDefaults, "template-var-defaults", "D", []string{}, "the defaults to use for the template variables")

	f.StringArrayVarP(&inputMetricsToDashStripPrefixes, "strip-prefix", "s", []string{}, "a prefix that can be removed from the metrics")
	f.BoolVarP(&inputMetricsToDashPretty, "pretty-name", "P", true, "prettify metric names")

	f.BoolVarP(&inputMetricsToDashShowHelp, "show-help", "S", false, "show help text for each metric")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// MetricsToDashCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
