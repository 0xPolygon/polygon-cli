package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// metricInfo represents a Prometheus metric definition.
type metricInfo struct {
	Name      string
	Namespace string
	Help      string
	Type      string // Gauge, Counter, Histogram, Summary, GaugeVec, CounterVec, etc.
	Labels    []string
	File      string
}

// fullName returns the complete metric name with namespace.
func (m metricInfo) fullName() string {
	if m.Namespace != "" {
		return m.Namespace + "_" + m.Name
	}
	return m.Name
}

// genMetricsDoc generates a METRICS.md file documenting all Prometheus metrics.
func genMetricsDoc(outputPath string) error {
	metrics := []metricInfo{}

	// Parse p2p package metrics
	p2pMetrics, err := parseMetricsFile("p2p/metrics.go")
	if err != nil {
		return fmt.Errorf("failed to parse p2p/metrics.go: %w", err)
	}
	metrics = append(metrics, p2pMetrics...)

	// Parse sensor command metrics
	sensorMetrics, err := parseMetricsFile("cmd/p2p/sensor/sensor.go")
	if err != nil {
		return fmt.Errorf("failed to parse cmd/p2p/sensor/sensor.go: %w", err)
	}
	metrics = append(metrics, sensorMetrics...)

	// Sort metrics by full name
	sort.Slice(metrics, func(i, j int) bool {
		return metrics[i].fullName() < metrics[j].fullName()
	})

	// Generate markdown
	return writeMetricsMarkdown(metrics, outputPath)
}

// parseMetricsFile extracts Prometheus metrics from a Go file.
func parseMetricsFile(filePath string) ([]metricInfo, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	metrics := []metricInfo{}
	relPath, _ := filepath.Rel(".", filePath)

	ast.Inspect(node, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		// Look for promauto.New* calls
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		pkg, ok := sel.X.(*ast.Ident)
		if !ok || pkg.Name != "promauto" {
			return true
		}

		metricType := strings.TrimPrefix(sel.Sel.Name, "New")
		if metricType == "" {
			return true
		}

		// Extract the options (first argument should be *Opts struct)
		if len(call.Args) == 0 {
			return true
		}

		metric := metricInfo{
			Type: metricType,
			File: relPath,
		}

		// Parse the options struct
		if comp, ok := call.Args[0].(*ast.CompositeLit); ok {
			for _, elt := range comp.Elts {
				if kv, ok := elt.(*ast.KeyValueExpr); ok {
					key := kv.Key.(*ast.Ident).Name

					switch key {
					case "Name":
						if lit, ok := kv.Value.(*ast.BasicLit); ok {
							metric.Name = strings.Trim(lit.Value, `"`)
						}
					case "Namespace":
						if lit, ok := kv.Value.(*ast.BasicLit); ok {
							metric.Namespace = strings.Trim(lit.Value, `"`)
						}
					case "Help":
						if lit, ok := kv.Value.(*ast.BasicLit); ok {
							metric.Help = strings.Trim(lit.Value, `"`)
						}
					}
				}
			}
		}

		// For Vec types, extract labels from second argument
		if strings.HasSuffix(metricType, "Vec") && len(call.Args) > 1 {
			if comp, ok := call.Args[1].(*ast.CompositeLit); ok {
				for _, elt := range comp.Elts {
					if lit, ok := elt.(*ast.BasicLit); ok {
						label := strings.Trim(lit.Value, `"`)
						metric.Labels = append(metric.Labels, label)
					}
				}
			}
		}

		if metric.Name != "" {
			metrics = append(metrics, metric)
		}

		return true
	})

	return metrics, nil
}

// writeMetricsMarkdown generates the METRICS.md file in panoptichain format.
func writeMetricsMarkdown(metrics []metricInfo, outputPath string) error {
	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintln(f)
	fmt.Fprintln(f, "## Sensor Metrics")
	fmt.Fprintln(f)

	// Group metrics by category (BlockMetrics vs sensor.go metrics)
	// For now, just output them in order with their descriptions
	for _, m := range metrics {
		fmt.Fprintf(f, "\n### %s\n", m.fullName())
		fmt.Fprintln(f, m.Help)
		fmt.Fprintln(f)
		fmt.Fprintf(f, "Metric Type: %s\n", m.Type)

		if len(m.Labels) > 0 {
			fmt.Fprintln(f)
			fmt.Fprintln(f, "Variable Labels:")
			for _, label := range m.Labels {
				fmt.Fprintf(f, "- %s\n", label)
			}
		}

		fmt.Fprintln(f)
	}

	return nil
}
