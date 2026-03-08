package python

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"time"
)

const (
	// DefaultPythonBin is the default Python executable name
	DefaultPythonBin = "python3"
	// DefaultTimeout is the default request timeout
	DefaultTimeout = 30 * time.Second
)

// CLIConfig holds configuration for Python CLI execution
type CLIConfig struct {
	PythonBin string
	Timeout   time.Duration
}

// CLI represents a Python CLI wrapper
type CLI struct {
	config CLIConfig
}

// NewCLI creates a new Python CLI wrapper
func NewCLI(pythonBin string, timeout time.Duration) (*CLI, error) {
	if pythonBin == "" {
		pythonBin = DefaultPythonBin
	}
	if timeout == 0 {
		timeout = DefaultTimeout
	}

	cli := &CLI{
		config: CLIConfig{
			PythonBin: pythonBin,
			Timeout:   timeout,
		},
	}

	// Check if Python package is available
	if err := cli.checkPackageAvailable(); err != nil {
		return nil, err
	}

	return cli, nil
}

// checkPackageAvailable verifies that youtube_transcript_api is installed
func (c *CLI) checkPackageAvailable() error {
	cmd := exec.Command(
		c.config.PythonBin,
		"-c",
		"import youtube_transcript_api",
	)

	_, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"youtube_transcript_api must be installed. Run: pip install youtube-transcript-api\nDetails: %w",
			err,
		)
	}

	return nil
}

// GetVersion returns the version of youtube_transcript_api
func (c *CLI) GetVersion() (string, error) {
	cmd := exec.Command(
		c.config.PythonBin,
		"-c",
		"import importlib.metadata; print(importlib.metadata.version('youtube-transcript-api'))",
	)

	output, err := cmd.Output()
	if err != nil {
		// Fallback to package being installed
		return "unknown", nil
	}

	// Remove trailing newline
	version := string(output)
	if len(version) > 0 && version[len(version)-1] == '\n' {
		version = version[:len(version)-1]
	}

	return version, nil
}

// FetchTranscripts fetches transcripts for the given video IDs with options
func (c *CLI) FetchTranscripts(
	videoIDs []string,
	languages []string,
	format string,
	excludeGenerated bool,
	excludeManuallyCreated bool,
	translate string,
) (map[string]interface{}, error) {
	// Default to pretty format if not specified
	if format == "" {
		format = "pretty"
	}

	// Default to English if no languages specified
	if len(languages) == 0 {
		languages = []string{"en"}
	}

	results := make(map[string]interface{})

	// Execute Python CLI for each video ID
	for _, videoID := range videoIDs {
		args := c.buildFetchArgs(
			videoID,
			languages,
			format,
			excludeGenerated,
			excludeManuallyCreated,
			translate,
		)

		output, err := c.executeCommand(args)
		if err != nil {
			results[videoID] = map[string]interface{}{
				"error": err.Error(),
			}
			continue
		}

		// Parse output based on format
		var data interface{}
		if format == "json" {
			var jsonData []map[string]interface{}
			if err := json.Unmarshal([]byte(output), &jsonData); err != nil {
				data = output
			} else {
				data = jsonData
			}
		} else {
			data = output
		}

		results[videoID] = data
	}

	return results, nil
}

// ListTranscripts lists available transcripts for the given video IDs
func (c *CLI) ListTranscripts(videoIDs []string) (map[string]interface{}, error) {
	results := make(map[string]interface{})

	for _, videoID := range videoIDs {
		args := []string{
			"-m", "youtube_transcript_api",
			"--list-transcripts",
			videoID,
		}

		output, err := c.executeCommand(args)
		if err != nil {
			results[videoID] = map[string]interface{}{
				"error": err.Error(),
			}
			continue
		}

		results[videoID] = output
	}

	return results, nil
}

// buildFetchArgs builds the CLI arguments for fetching transcripts
func (c *CLI) buildFetchArgs(
	videoID string,
	languages []string,
	format string,
	excludeGenerated bool,
	excludeManuallyCreated bool,
	translate string,
) []string {
	args := []string{"-m", "youtube_transcript_api"}

	// Add languages first (before positional args)
	if len(languages) > 0 {
		args = append(args, "--languages")
		args = append(args, languages...)
	}

	// Add format flag
	args = append(args, "--format", format)

	// Add exclude flags
	if excludeGenerated {
		args = append(args, "--exclude-generated")
	}
	if excludeManuallyCreated {
		args = append(args, "--exclude-manually-created")
	}

	// Add translate flag if specified
	if translate != "" {
		args = append(args, "--translate", translate)
	}

	// Add video ID LAST (it's a positional argument)
	args = append(args, videoID)

	return args
}

// executeCommand executes a Python command with timeout
func (c *CLI) executeCommand(args []string) (string, error) {
	cmd := exec.Command(c.config.PythonBin, args...)

	// Set up a timeout using context
	done := make(chan error, 1)
	var output []byte

	go func() {
		var err error
		output, err = cmd.CombinedOutput()
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			return "", fmt.Errorf("command failed: %s", string(output))
		}
		return string(output), nil
	case <-time.After(c.config.Timeout):
		cmd.Process.Kill()
		return "", fmt.Errorf("command timed out after %v", c.config.Timeout)
	}
}
