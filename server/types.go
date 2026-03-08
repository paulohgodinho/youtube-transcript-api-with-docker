package main

// TranscriptRequest represents a request to fetch transcripts for one or more YouTube videos.
//
// VideoIDs is required and must contain at least one video ID.
// All other fields are optional and allow filtering and transformation of the results.
//
// @Description YouTube video transcript fetch request with optional filtering and translation
type TranscriptRequest struct {
	// VideoIDs is a required list of YouTube video IDs to fetch transcripts for.
	// Each ID should be a valid YouTube video ID (typically 11 characters).
	// Example: "dQw4w9WgXcQ" (Rick Astley - Never Gonna Give You Up)
	VideoIDs []string `json:"videoIds" example:"dQw4w9WgXcQ,jNQXAC9IVRw"`

	// Languages is an optional list of language codes to filter transcripts by (e.g., 'en', 'es', 'fr').
	// If not specified, the API defaults to English ('en').
	// Use standard ISO 639-1 language codes.
	Languages []string `json:"languages,omitempty" example:"en,es,fr"`

	// Format specifies the output format for the transcript content.
	// Valid values: 'json' (structured with timing information), 'text' (plain text only).
	// Defaults to 'json' if not specified.
	Format string `json:"format,omitempty" example:"json" default:"json"`

	// ExcludeGenerated, when true, excludes auto-generated transcripts from the results.
	// YouTube auto-generates transcripts using speech recognition if no manual transcript exists.
	// Defaults to false (includes auto-generated transcripts).
	ExcludeGenerated bool `json:"excludeGenerated,omitempty" example:"false" default:"false"`

	// ExcludeManuallyCreated, when true, excludes manually created transcripts from the results.
	// Defaults to false (includes manually created transcripts).
	ExcludeManuallyCreated bool `json:"excludeManuallyCreated,omitempty" example:"false" default:"false"`

	// Translate, when specified with a language code, translates the transcript to that language.
	// Only works if the transcript's isTranslatable field is true.
	// Example: "es" translates to Spanish, "fr" translates to French.
	// Leave empty or omit to skip translation.
	Translate string `json:"translate,omitempty" example:"es"`
}

// ListRequest represents a request to list available transcripts for one or more YouTube videos.
//
// This endpoint doesn't fetch the actual transcript content, only metadata about
// what transcripts are available for each video (language, source, availability of translations).
//
// @Description YouTube video transcript listing request for metadata discovery
type ListRequest struct {
	// VideoIDs is a required list of YouTube video IDs to list available transcripts for.
	// Each ID should be a valid YouTube video ID (typically 11 characters).
	// Example: "dQw4w9WgXcQ" (Rick Astley - Never Gonna Give You Up)
	VideoIDs []string `json:"videoIds" example:"dQw4w9WgXcQ,jNQXAC9IVRw"`
}

// TranscriptSnippet represents a single segment or snippet of a transcript.
//
// Each snippet contains the text content and precise timing information, allowing
// for synchronization with video playback.
//
// @Description A single segment of transcript text with timing information
type TranscriptSnippet struct {
	// Text is the actual transcript content for this snippet.
	// Contains the words spoken during the time period defined by Start and Duration.
	Text string `json:"text" example:"Welcome to this video, today we'll discuss something interesting"`

	// Start is the start time of this snippet in seconds from the beginning of the video.
	// Useful for seeking to a specific point in the video during playback.
	Start float64 `json:"start" example:"0.0"`

	// Duration is the length of this snippet in seconds.
	// The end time would be: Start + Duration.
	Duration float64 `json:"duration" example:"2.5"`
}

// FetchTranscriptResult represents the result of fetching a transcript for a single video.
//
// Contains the actual transcript content (broken into snippets with timing),
// language information, and metadata about the transcript source.
//
// @Description Complete transcript for a single video with full metadata
type FetchTranscriptResult struct {
	// VideoID is the YouTube video ID this transcript belongs to.
	// Matches one of the video IDs from the request.
	VideoID string `json:"videoId" example:"dQw4w9WgXcQ"`

	// Snippets is an array of transcript segments, each with text and timing information.
	// Ordered chronologically from the beginning to the end of the video.
	// Useful for creating timestamped searchable transcripts.
	Snippets []TranscriptSnippet `json:"snippets"`

	// Language is the human-readable language name (e.g., "English", "Spanish", "French").
	// Useful for display purposes in user interfaces.
	Language string `json:"language" example:"English"`

	// LanguageCode is the ISO 639-1 language code for this transcript (e.g., "en", "es", "fr").
	// Standard format used for language identification in APIs and databases.
	LanguageCode string `json:"languageCode" example:"en"`

	// IsGenerated indicates whether this transcript was auto-generated by YouTube.
	// True means YouTube generated it using speech recognition; false means it was manually created by the uploader.
	// Auto-generated transcripts may contain errors, especially with names, technical terms, or accents.
	IsGenerated bool `json:"isGenerated" example:"false"`
}

// TranscriptMetadata represents metadata about an available transcript without its content.
//
// Used by the /list endpoint to describe what transcripts are available for a video
// without fetching the full transcript text.
//
// @Description Metadata describing an available transcript (without content)
type TranscriptMetadata struct {
	// Language is the human-readable language name (e.g., "English", "Spanish", "French").
	// Useful for display purposes in user interfaces.
	Language string `json:"language" example:"English"`

	// LanguageCode is the ISO 639-1 language code for this transcript (e.g., "en", "es", "fr").
	// Standard format used for language identification in APIs and databases.
	LanguageCode string `json:"languageCode" example:"en"`

	// IsGenerated indicates whether this transcript was auto-generated by YouTube.
	// True means YouTube generated it using speech recognition; false means it was manually created.
	IsGenerated bool `json:"isGenerated" example:"false"`

	// IsTranslatable indicates whether this transcript can be translated to other languages.
	// True means the API can provide translations to other languages via the translate parameter.
	// Check TranslationLanguages for the list of available translation targets.
	IsTranslatable bool `json:"isTranslatable" example:"true"`

	// TranslationLanguages is a list of ISO 639-1 language codes to which this transcript can be translated.
	// Only populated if IsTranslatable is true.
	// Example: ["es", "fr", "de"] means this transcript can be translated to Spanish, French, or German.
	TranslationLanguages []string `json:"translationLanguages,omitempty" example:"es,fr,de"`
}

// ListTranscriptResult represents the available transcripts for a single video.
//
// Organized into three categories: manually created, auto-generated, and translation language options.
// Each transcript is represented as TranscriptMetadata without the actual content.
//
// @Description Available transcripts for a single video (metadata only, no content)
type ListTranscriptResult struct {
	// VideoID is the YouTube video ID this metadata belongs to.
	// Matches one of the video IDs from the request.
	VideoID string `json:"videoId" example:"dQw4w9WgXcQ"`

	// ManuallyCreated is a list of manually created transcripts available for this video.
	// These are typically more accurate and complete than auto-generated ones.
	// Created by the video uploader or YouTube community contributors.
	ManuallyCreated []TranscriptMetadata `json:"manuallyCreated"`

	// Generated is a list of auto-generated transcripts available for this video.
	// YouTube creates these automatically using speech recognition if no manual transcript exists.
	// May contain errors, especially with names, technical terms, or non-native speakers.
	Generated []TranscriptMetadata `json:"generated"`

	// TranslationLanguages lists all languages to which available transcripts can be translated.
	// This represents the combined translation capabilities of all available transcripts for this video.
	TranslationLanguages []TranscriptMetadata `json:"translationLanguages"`
}

// TranscriptResponse is the standard response wrapper for transcript operations.
//
// The Transcripts field content varies depending on the endpoint:
// - /transcripts endpoint: contains []FetchTranscriptResult with actual transcript content
// - /list endpoint: contains []ListTranscriptResult with metadata only
//
// When Success is false, the Transcripts field will be empty and Error will contain details.
//
// @Description Standard response wrapper for transcript API operations
type TranscriptResponse struct {
	// Success indicates whether the request was processed successfully.
	// True if all requested transcripts were fetched/listed without errors.
	// False if any error occurred during processing.
	Success bool `json:"success" example:"true"`

	// Transcripts contains the response data.
	// Type varies based on the endpoint called:
	// - For /transcripts: array of FetchTranscriptResult objects
	// - For /list: array of ListTranscriptResult objects
	// Only populated if Success is true. Omitted if Success is false.
	Transcripts interface{} `json:"transcripts,omitempty"`

	// Error contains the error message if Success is false.
	// Empty or omitted if Success is true.
	// Provides details about what went wrong (missing parameters, API errors, etc.).
	Error string `json:"error,omitempty" example:"videoIds is required and must not be empty"`
}

// HealthResponse represents a health check response from the server.
//
// @Description Server health status indicator
type HealthResponse struct {
	// Status indicates the server health status.
	// Always "ok" if the server is running and responding.
	// Used for monitoring, load balancer health checks, and uptime verification.
	Status string `json:"status" example:"ok"`
}

// VersionResponse represents the version of the youtube-transcript-api CLI.
//
// @Description Version information of the youtube-transcript-api Python CLI package
type VersionResponse struct {
	// Version is the semantic version string of the youtube-transcript-api package
	// being used by the server (e.g., "0.6.2", "0.7.0").
	// Helps verify that the correct CLI tool version is installed and accessible.
	// Returns "unknown" if the version cannot be determined.
	Version string `json:"version" example:"0.6.2"`
}

// ErrorResponse represents a standard error response from the API.
//
// All errors from the API follow this standard format for consistency.
//
// @Description Standard API error response with success flag and error message
type ErrorResponse struct {
	// Success is always false for error responses.
	// Allows clients to check response status consistently across all endpoints.
	Success bool `json:"success" example:"false"`

	// Error contains a human-readable error message describing what went wrong.
	// Provides sufficient detail for clients to understand and potentially fix the issue.
	// Examples: "videoIds is required and must not be empty", "Invalid JSON in request body", etc.
	Error string `json:"error" example:"videoIds is required and must not be empty"`
}
