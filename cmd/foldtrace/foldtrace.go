package foldtrace

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	_ "embed"
)

var (
	//go:embed usage.md
	usage                string
	inputFileName        string
	inputRootContextName string
	inputMetric          string
)

type TraceOperation struct {
	PC      uint64   `json:"pc"`
	OP      string   `json:"op"`
	Gas     uint64   `json:"gas"`
	GasCost uint64   `json:"gasCost"`
	Depth   uint64   `json:"depth"`
	Stack   []string `json:"stack"`
}
type TraceData struct {
	StructLogs []TraceOperation `json:"structLogs"`
}

var FoldTraceCmd = &cobra.Command{
	Use:   "fold-trace",
	Short: "Trace an execution trace and fold it for visualization.",
	Long:  usage,
	Args:  cobra.ArbitraryArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateInputMetric()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := getInputData(cmd, args)
		if err != nil {
			return fmt.Errorf("unable to read input data: %w", err)
		}
		td := &TraceData{}
		err = json.Unmarshal(data, td)
		if err != nil {
			return fmt.Errorf("unable to unmarshal input json trace: %w", err)
		}

		folded := make(map[string]uint64)

		contexts := []string{inputRootContextName}
		currentDepth := uint64(1)
		lastLabel := ""

		metricType := inputMetric
		for idx, op := range td.StructLogs {
			if op.Depth > currentDepth {
				contexts = append(contexts, lastLabel)
			} else if op.Depth < currentDepth {
				contexts = contexts[:len(contexts)-1]
			}
			currentDepth = op.Depth

			currentLabel := op.OP

			if isCall(op.OP) {
				if len(op.Stack) < 6 {
					log.Warn().Int("stackLength", len(op.Stack)).Msg("detected a call with a stack that's too short")
				} else {
					currentLabel = op.OP + " to " + op.Stack[len(op.Stack)-2]
				}
			}
			lastLabel = currentLabel

			var metricValue uint64
			if metricType == "gas" {
				metricValue = op.GasCost
			} else if metricType == "actualgas" {
				metricValue = getActualUsedGas(idx, td)
			} else {
				metricValue = 1
			}
			currentMetricPath := strings.Join(contexts, ";") + ";" + currentLabel

			if !isCall(op.OP) {
				folded[currentMetricPath] += metricValue
			}
			if isCall(op.OP) && len(op.Stack[len(op.Stack)-2]) < 6 {
				folded[currentMetricPath] += metricValue
			}

			log.Trace().Strs("context", contexts).Uint64("depth", currentDepth).Str("currentLabel", currentLabel).Uint64("pc", op.PC).Msg("trace operation")

		}
		for context, metric := range folded {
			fmt.Printf("%s %d\n", context, metric)
		}

		return nil
	},
}

func getActualUsedGas(index int, td *TraceData) uint64 {
	op := td.StructLogs[index]
	if op.OP == "RETURN" {
		return op.GasCost
	}
	if op.OP == "STOP" {
		return op.GasCost
	}

	for i := index + 1; i < len(td.StructLogs); i++ {
		if op.Depth < td.StructLogs[i].Depth {
			continue
		}
		if op.Depth > td.StructLogs[i].Depth {
			break
		}

		return op.Gas - td.StructLogs[i].Gas
	}
	log.Warn().Uint64("pc", op.PC).Uint64("depth", op.Depth).Msg("unable to look ahead for gas use")
	return op.GasCost
}

func isCall(op string) bool {
	return strings.HasSuffix(op, "CALL")
}

func validateInputMetric() error {
	switch inputMetric {
	case "gas", "count", "actualgas":
		return nil
	}
	return fmt.Errorf("invalid input metric: %s", inputMetric)
}

func init() {
	f := FoldTraceCmd.Flags()
	f.StringVar(&inputFileName, "file", "", "filename to read and hash")
	f.StringVar(&inputRootContextName, "root-context", "root context", "name for top most initial context")
	f.StringVar(&inputMetric, "metric", "gas", "metric name for analysis: gas, count, actualgas")

}

func getInputData(cmd *cobra.Command, args []string) ([]byte, error) {
	// first check and see if we have an input file
	if inputFileName != "" {
		// If we get here, we're going to assume the user
		// wants to hash a file and we're not going to look
		// for other input sources
		return os.ReadFile(inputFileName)
	}

	// This is a little tricky. If a user provides multiple args that aren't quoted, it could be confusing
	if len(args) > 1 {
		concat := strings.Join(args[1:], " ")
		return []byte(concat), nil
	}

	return io.ReadAll(os.Stdin)
}
