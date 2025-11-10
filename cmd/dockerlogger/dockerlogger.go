package dockerlogger

import (
	"bufio"
	"context"
	_ "embed"
	"fmt"
	"io"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

//go:embed dockerlogger.md
var cmdUsage string

type inputArgs struct {
	network string
	filter  string
	levels  string
	service string
}

var dockerloggerInputArgs = inputArgs{}

var (
	// Colors for log output components
	timestampColor   = color.New(color.FgCyan)
	serviceNameColor = color.New(color.FgBlue)
	errorLevelColor  = color.New(color.FgRed, color.Bold)
	warnLevelColor   = color.New(color.FgYellow, color.Bold)
	infoLevelColor   = color.New(color.FgGreen)
	debugLevelColor  = color.New(color.FgMagenta)
	messageColor     = color.New(color.Reset) // Normal text color
)

// Types
type LogConfig struct {
	customWords  string
	logLevels    string
	serviceNames []string
}

// Main function
func dockerlogger(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	if dockerloggerInputArgs.network == "" {
		return fmt.Errorf("--network flag is required")
	}

	config := LogConfig{
		customWords: dockerloggerInputArgs.filter,
		logLevels:   dockerloggerInputArgs.levels,
	}

	if dockerloggerInputArgs.service != "" {
		config.serviceNames = strings.Split(dockerloggerInputArgs.service, ",")
	}

	// Set up Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("error initializing Docker client: %v", err)
	}
	defer cli.Close()

	// Monitor logs
	return monitorLogs(ctx, cli, dockerloggerInputArgs.network, &config)
}

// Cobra command for Docker logger
var Cmd = &cobra.Command{
	Use:   "dockerlogger",
	Short: "Monitor and filter Docker container logs.",
	Long:  cmdUsage,
	RunE:  dockerlogger,
}

func init() {
	f := Cmd.Flags()
	f.StringVar(&dockerloggerInputArgs.network, "network", "", "docker network name to monitor")
	f.StringVar(&dockerloggerInputArgs.filter, "filter", "", "additional keywords to filter, comma-separated")
	f.StringVar(&dockerloggerInputArgs.levels, "levels", "", "comma-separated log levels to show (error,warn,info,debug)")
	f.StringVar(&dockerloggerInputArgs.service, "service", "", "filter logs by service names (comma-separated, partial match)")
}

// Core functionality functions
func CreateClient() (*client.Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return cli, nil
}

// InspectNetwork retrieves detailed information about a Docker network.
func InspectNetwork(ctx context.Context, cli *client.Client, networkName string) (network.Inspect, error) {
	if networkName == "" {
		return network.Inspect{}, fmt.Errorf("network name cannot be empty")
	}

	if cli == nil {
		return network.Inspect{}, fmt.Errorf("docker client cannot be nil")
	}

	net, err := cli.NetworkInspect(ctx, networkName, network.InspectOptions{
		Verbose: true,
	})
	if err != nil {
		return network.Inspect{}, fmt.Errorf("failed to inspect network %s: %w", networkName, err)
	}

	return net, nil
}

func monitorLogs(ctx context.Context, cli *client.Client, networkName string, config *LogConfig) error {
	network, err := InspectNetwork(ctx, cli, networkName)
	if err != nil {
		return fmt.Errorf("error inspecting network: %v", err)
	}

	containers := network.Containers
	if len(containers) == 0 {
		return fmt.Errorf("no containers found in network '%s'", networkName)
	}

	fmt.Printf("Monitoring logs for network '%s'...\n", networkName)

	var wg sync.WaitGroup
	for containerID, containerInfo := range containers {
		wg.Add(1)
		go func(id, name string) {
			defer wg.Done()
			streamContainerLogs(ctx, cli, id, name, config)
		}(containerID, containerInfo.Name)
	}

	wg.Wait()
	return nil
}

func streamContainerLogs(ctx context.Context, cli *client.Client, containerID, containerName string, config *LogConfig) {
	containerInfo, err := cli.ContainerInspect(ctx, containerID)
	if err != nil {
		fmt.Printf("Error inspecting container %s: %v\n", containerName, err)
		return
	}

	serviceName := containerName
	if labels, exists := containerInfo.Config.Labels["com.docker.compose.service"]; exists {
		serviceName = labels
	}

	if !matchesServiceName(containerName, config.serviceNames) {
		return
	}

	// Update to use container.LogsOptions instead of types.ContainerLogsOptions
	opts := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Timestamps: false,
	}

	logs, err := cli.ContainerLogs(ctx, containerID, opts)
	if err != nil {
		fmt.Printf("Error streaming logs for %s: %v\n", serviceName, err)
		return
	}
	defer logs.Close()

	fmt.Printf("Started monitoring container: %s\n", serviceName)

	// Create a pipe to demultiplex Docker's stdout/stderr stream
	reader, writer := io.Pipe()
	defer reader.Close()
	defer writer.Close()

	// Demultiplex the Docker log stream in a goroutine
	go func() {
		// stdcopy.StdCopy properly handles Docker's 8-byte header format
		_, err := stdcopy.StdCopy(writer, writer, logs)
		if err != nil && err != io.EOF {
			fmt.Printf("Error demultiplexing logs for %s: %v\n", serviceName, err)
		}
		writer.Close()
	}()

	// Read demultiplexed logs line by line
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		logLine := scanner.Text()
		if logLine == "" {
			continue
		}

		// Sanitize the log line
		logLine = sanitizeLogLine(logLine)
		if logLine == "" {
			continue
		}

		logLineLower := strings.ToLower(logLine)
		if !shouldLogMessage(logLineLower, config) {
			continue
		}

		// Format timestamp and print log with colored components
		timestamp := time.Now().UTC().Format("2006-01-02 15:04:05")

		// Print with different colors for each component
		timestampColor.Printf("[%s] ", timestamp)
		serviceNameColor.Printf("[%s] ", serviceName)

		// Colorize the log level within the message and print the rest normally
		printColorizedLogLine(logLine, logLineLower)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading logs for %s: %v\n", serviceName, err)
	}
}

// printColorizedLogLine prints a log line with appropriate colors for log levels
func printColorizedLogLine(logLine, logLineLower string) {
	// Define log level patterns to search for
	logLevelPatterns := []struct {
		pattern string
		color   *color.Color
	}{
		{"ERROR", errorLevelColor},
		{"ERRO", errorLevelColor},
		{"EROR", errorLevelColor},
		{"ERR", errorLevelColor},
		{"WARNING", warnLevelColor},
		{"WARN", warnLevelColor},
		{"WRN", warnLevelColor},
		{"INFO", infoLevelColor},
		{"INF", infoLevelColor},
		{"DEBUG", debugLevelColor},
		{"DBG", debugLevelColor},
	}

	// Find the log level in the message
	foundLevel := false
	for _, lp := range logLevelPatterns {
		// Case-insensitive search for the pattern
		patternLower := strings.ToLower(lp.pattern)
		if idx := strings.Index(logLineLower, patternLower); idx != -1 {
			// Print everything before the log level
			if idx > 0 {
				messageColor.Print(logLine[:idx])
			}

			// Print the log level with its color (preserve original case)
			levelEnd := idx + len(lp.pattern)
			lp.color.Print(logLine[idx:levelEnd])

			// Print everything after the log level
			if levelEnd < len(logLine) {
				messageColor.Print(logLine[levelEnd:])
			}
			fmt.Println()
			foundLevel = true
			break
		}
	}

	// If no log level found, print the entire line normally
	if !foundLevel {
		messageColor.Println(logLine)
	}
}

// Log filtering and processing functions
func shouldLogMessage(logLine string, config *LogConfig) bool {
	// Parse configured log levels
	var allowedLevels map[string]bool
	if config.logLevels != "" {
		allowedLevels = make(map[string]bool)
		for level := range strings.SplitSeq(config.logLevels, ",") {
			allowedLevels[strings.TrimSpace(strings.ToLower(level))] = true
		}
	}

	// Check if message matches any allowed log level
	if len(allowedLevels) > 0 {
		if allowedLevels["error"] && isErrorMessage(logLine) {
			return true
		}
		if allowedLevels["warn"] && isWarningMessage(logLine) {
			return true
		}
		if allowedLevels["info"] && isInfoMessage(logLine) {
			return true
		}
		if allowedLevels["debug"] && isDebugMessage(logLine) {
			return true
		}
	}

	// Check custom keywords if no level match
	if config.customWords != "" {
		customKeywords := strings.SplitSeq(config.customWords, ",")
		for keyword := range customKeywords {
			if strings.Contains(logLine, strings.TrimSpace(strings.ToLower(keyword))) {
				return true
			}
		}
	}

	// If no levels or keywords specified, show all logs
	if config.logLevels == "" && config.customWords == "" {
		return true
	}

	return false
}

func matchesServiceName(containerName string, serviceFilters []string) bool {
	if len(serviceFilters) == 0 {
		return true
	}
	for _, filter := range serviceFilters {
		if strings.Contains(strings.ToLower(containerName), strings.ToLower(filter)) {
			return true
		}
	}
	return false
}

// Log level detection functions
func isErrorMessage(logLine string) bool {
	return strings.Contains(strings.ToLower(logLine), "error")
}

func isWarningMessage(logLine string) bool {
	return strings.Contains(strings.ToLower(logLine), "warn")
}

func isInfoMessage(logLine string) bool {
	return strings.Contains(strings.ToLower(logLine), "info")
}

func isDebugMessage(logLine string) bool {
	return strings.Contains(strings.ToLower(logLine), "debug")
}

func sanitizeLogLine(logLine string) string {
	// Remove ANSI color codes
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	logLine = ansiRegex.ReplaceAllString(logLine, "")

	// Remove any remaining non-printable characters except common whitespace
	var result strings.Builder
	for _, r := range logLine {
		// Keep printable characters and common whitespace (space, tab, newline)
		if r >= 32 || r == '\t' || r == '\n' || r == '\r' {
			result.WriteRune(r)
		}
	}

	return strings.TrimSpace(result.String())
}
