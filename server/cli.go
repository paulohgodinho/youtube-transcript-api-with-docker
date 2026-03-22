package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
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
) ([]FetchTranscriptResult, error) {
	if len(languages) == 0 {
		languages = []string{"en"}
	}

	var results []FetchTranscriptResult

	for _, videoID := range videoIDs {
		args := c.buildFetchArgs(
			videoID,
			languages,
			"json",
			excludeGenerated,
			excludeManuallyCreated,
			translate,
		)

		output, err := c.executeCommand(args)
		if err != nil {
			continue
		}

		result, err := c.parseFetchOutput(videoID, output)
		if err != nil {
			continue
		}

		results = append(results, result)
	}

	return results, nil
}

// ListTranscripts lists available transcripts for the given video IDs
func (c *CLI) ListTranscripts(videoIDs []string) ([]ListTranscriptResult, error) {
	var results []ListTranscriptResult

	for _, videoID := range videoIDs {
		args := []string{
			"-m", "youtube_transcript_api",
			"--list-transcripts",
			videoID,
		}

		output, err := c.executeCommand(args)
		if err != nil {
			continue
		}

		result, err := c.parseListOutput(videoID, output)
		if err != nil {
			continue
		}

		results = append(results, result)
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

// parseFetchOutput parses the JSON output from youtube_transcript_api fetch command
func (c *CLI) parseFetchOutput(videoID string, output string) (FetchTranscriptResult, error) {
	// Output format is nested array: [[{text, start, duration}, ...]]
	var nestedSnippets [][]TranscriptSnippet
	err := json.Unmarshal([]byte(output), &nestedSnippets)
	if err != nil {
		return FetchTranscriptResult{}, fmt.Errorf("failed to parse transcript JSON: %w", err)
	}

	// Extract the first (and should be only) inner array
	var snippets []TranscriptSnippet
	if len(nestedSnippets) > 0 {
		snippets = nestedSnippets[0]
	}

	return FetchTranscriptResult{
		VideoID:      videoID,
		Snippets:     snippets,
		Language:     "en",
		LanguageCode: "en",
		IsGenerated:  false,
	}, nil
}

// parseListOutput parses the text output from youtube_transcript_api list-transcripts command
func (c *CLI) parseListOutput(videoID string, output string) (ListTranscriptResult, error) {
	lines := strings.Split(output, "\n")

	result := ListTranscriptResult{
		VideoID:              videoID,
		ManuallyCreated:      []TranscriptMetadata{},
		Generated:            []TranscriptMetadata{},
		TranslationLanguages: []TranscriptMetadata{},
	}

	var currentSection string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.Contains(trimmed, "(MANUALLY CREATED)") {
			currentSection = "manually_created"
			continue
		}
		if strings.Contains(trimmed, "(GENERATED)") {
			currentSection = "generated"
			continue
		}
		if strings.Contains(trimmed, "(TRANSLATION LANGUAGES)") {
			currentSection = "translation"
			continue
		}

		if trimmed == "" || !strings.HasPrefix(trimmed, "- ") {
			continue
		}

		trimmed = strings.TrimPrefix(trimmed, "- ")

		metadata, err := parseTranscriptLine(trimmed)
		if err != nil {
			continue
		}

		switch currentSection {
		case "manually_created":
			result.ManuallyCreated = append(result.ManuallyCreated, metadata)
		case "generated":
			result.Generated = append(result.Generated, metadata)
		case "translation":
			result.TranslationLanguages = append(result.TranslationLanguages, metadata)
		}
	}

	return result, nil
}

// parseTranscriptLine parses a single transcript line from the list output
// Format: "en ("English")[TRANSLATABLE]" or "ar ("Arabic")"
func parseTranscriptLine(line string) (TranscriptMetadata, error) {
	metadata := TranscriptMetadata{}

	// Extract language code (first part before space or paren)
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return metadata, fmt.Errorf("empty line")
	}

	langCode := parts[0]
	metadata.LanguageCode = langCode

	// Extract language name from quotes
	startQuote := strings.Index(line, "(\"")
	endQuote := strings.Index(line, "\")")
	if startQuote != -1 && endQuote != -1 && startQuote < endQuote {
		startQuote += 2
		metadata.Language = line[startQuote:endQuote]
	}

	// Check if translatable
	if strings.Contains(line, "[TRANSLATABLE]") {
		metadata.IsTranslatable = true
	}

	// Check if auto-generated
	if strings.Contains(metadata.Language, "auto-generated") || strings.Contains(metadata.Language, "(auto-generated)") {
		metadata.IsGenerated = true
	}

	return metadata, nil
}
