package dockerlogger

import (
	"bufio"
	"context"
	_ "embed"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

//go:embed dockerlogger.md
var cmdUsage string

type inputArgs struct {
	network      *string
	showAll      *bool
	showErrors   *bool
	showWarnings *bool
	showInfo     *bool
	showDebug    *bool
	filter       *string
	levels       *string
	service      *string
}

var dockerloggerInputArgs = inputArgs{}

var (
	// Colors for log output
	normalColor  = color.New(color.FgGreen)
	warningColor = color.New(color.FgYellow, color.Bold)
	errorColor   = color.New(color.FgRed, color.Bold)
)

// Types
type LogConfig struct {
	showAll      bool
	showErrors   bool
	showWarns    bool
	showInfo     bool
	showDebug    bool
	customWords  string
	logLevels    string
	serviceNames []string
}

// Main function
func dockerlogger(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	if *dockerloggerInputArgs.network == "" {
		return fmt.Errorf("--network flag is required")
	}

	config := LogConfig{
		showAll:     *dockerloggerInputArgs.showAll,
		showErrors:  *dockerloggerInputArgs.showErrors,
		showWarns:   *dockerloggerInputArgs.showWarnings,
		showInfo:    *dockerloggerInputArgs.showInfo,
		showDebug:   *dockerloggerInputArgs.showDebug,
		customWords: *dockerloggerInputArgs.filter,
		logLevels:   *dockerloggerInputArgs.levels,
	}

	if *dockerloggerInputArgs.service != "" {
		config.serviceNames = strings.Split(*dockerloggerInputArgs.service, ",")
	}

	// Set up Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("error initializing Docker client: %v", err)
	}
	defer cli.Close()

	// Monitor logs
	return monitorLogs(ctx, cli, *dockerloggerInputArgs.network, &config)
}

// Cobra command for Docker logger
var Cmd = &cobra.Command{
	Use:   "dockerlogger",
	Short: "Monitor and filter Docker container logs",
	Long:  cmdUsage,
	RunE:  dockerlogger,
}

func init() {
	dockerloggerInputArgs.network = Cmd.Flags().String("network", "", "Docker network name to monitor")
	dockerloggerInputArgs.showAll = Cmd.Flags().Bool("all", false, "Show all logs")
	dockerloggerInputArgs.showErrors = Cmd.Flags().Bool("errors", false, "Show error logs")
	dockerloggerInputArgs.showWarnings = Cmd.Flags().Bool("warnings", false, "Show warning logs")
	dockerloggerInputArgs.showInfo = Cmd.Flags().Bool("info", false, "Show info logs")
	dockerloggerInputArgs.showDebug = Cmd.Flags().Bool("debug", false, "Show debug logs")
	dockerloggerInputArgs.filter = Cmd.Flags().String("filter", "", "Additional keywords to filter, comma-separated")
	dockerloggerInputArgs.levels = Cmd.Flags().String("levels", "", "Comma-separated log levels to show (error,warn,info,debug)")
	dockerloggerInputArgs.service = Cmd.Flags().String("service", "", "Filter logs by service names (comma-separated, partial match)")
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
func InspectNetwork(ctx context.Context, cli *client.Client, networkName string) (types.NetworkResource, error) {
	if networkName == "" {
		return types.NetworkResource{}, fmt.Errorf("network name cannot be empty")
	}

	if cli == nil {
		return types.NetworkResource{}, fmt.Errorf("docker client cannot be nil")
	}

	network, err := cli.NetworkInspect(ctx, networkName, types.NetworkInspectOptions{
		Verbose: true,
	})
	if err != nil {
		return types.NetworkResource{}, fmt.Errorf("failed to inspect network %s: %w", networkName, err)
	}

	return network, nil
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

	// Read logs line by line
	scanner := bufio.NewScanner(logs)
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

		// Format timestamp and print log
		timestamp := time.Now().UTC().Format("2006-01-02 15:04:05")
		var logColor *color.Color
		if isErrorMessage(logLineLower) {
			logColor = errorColor
		} else if isWarningMessage(logLineLower) {
			logColor = warningColor
		} else {
			logColor = normalColor
		}

		logColor.Printf("[%s] [%s] %s\n", timestamp, serviceName, logLine)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading logs for %s: %v\n", serviceName, err)
	}
}

// Log filtering and processing functions
func shouldLogMessage(logLine string, config *LogConfig) bool {
	// If showAll is true, skip other checks
	if config.showAll {
		return true
	}

	// Parse configured log levels
	var allowedLevels map[string]bool
	if config.logLevels != "" {
		allowedLevels = make(map[string]bool)
		for _, level := range strings.Split(config.logLevels, ",") {
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
		customKeywords := strings.Split(config.customWords, ",")
		for _, keyword := range customKeywords {
			if strings.Contains(logLine, strings.TrimSpace(strings.ToLower(keyword))) {
				return true
			}
		}
	}

	// Check individual flag settings if no level match
	if config.showErrors && isErrorMessage(logLine) {
		return true
	}
	if config.showWarns && isWarningMessage(logLine) {
		return true
	}
	if config.showInfo && isInfoMessage(logLine) {
		return true
	}
	if config.showDebug && isDebugMessage(logLine) {
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
	return ansiRegex.ReplaceAllString(logLine, "")
}
