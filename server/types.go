package main

// TranscriptRequest represents a request to fetch transcripts for one or more videos
type TranscriptRequest struct {
	VideoIDs               []string `json:"videoIds"`
	Languages              []string `json:"languages,omitempty"`
	Format                 string   `json:"format,omitempty"`
	ExcludeGenerated       bool     `json:"excludeGenerated,omitempty"`
	ExcludeManuallyCreated bool     `json:"excludeManuallyCreated,omitempty"`
	Translate              string   `json:"translate,omitempty"`
}

// ListRequest represents a request to list available transcripts for videos
type ListRequest struct {
	VideoIDs []string `json:"videoIds"`
}

// TranscriptResponse represents the response containing transcript data
type TranscriptResponse struct {
	Success     bool             `json:"success"`
	Transcripts []TranscriptData `json:"transcripts,omitempty"`
	Error       string           `json:"error,omitempty"`
}

// TranscriptData represents a single transcript result
type TranscriptData struct {
	VideoID string      `json:"videoId"`
	Data    interface{} `json:"data"`
}

// HealthResponse represents a health check response
type HealthResponse struct {
	Status string `json:"status"`
}

// VersionResponse represents a version check response
type VersionResponse struct {
	Version string `json:"version"`
}

// ErrorResponse represents a generic error response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}
