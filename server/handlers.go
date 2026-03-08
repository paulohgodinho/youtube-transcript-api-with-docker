package main

import (
	"encoding/json"
	"io"
	"net/http"
)

// Handler holds dependencies for HTTP handlers
type Handler struct {
	cli *CLI
}

// NewHandler creates a new HTTP handler
func NewHandler(cli *CLI) *Handler {
	return &Handler{
		cli: cli,
	}
}

// Health checks if the server is running and responsive.
//
// @Summary      Health Check
// @Description  Checks if the server is running and responsive. Returns a simple status response indicating server health.
// @Description  This endpoint can be used for monitoring and load balancer health checks.
//
// @Tags         System
// @Accept       json
// @Produce      json
//
// @Success      200  {object}  HealthResponse  "Server is healthy"
// @Failure      405  {object}  ErrorResponse   "Method not allowed"
//
// @Router       /health [get]
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(HealthResponse{
		Status: "ok",
	})
}

// Version returns the version of the youtube-transcript-api CLI being used.
//
// @Summary      Get API Version
// @Description  Returns the version of the youtube-transcript-api Python package being used by the server.
// @Description  This helps verify that the correct CLI tool version is installed and accessible.
// @Description  If the version cannot be determined, returns "unknown".
//
// @Tags         System
// @Accept       json
// @Produce      json
//
// @Success      200  {object}  VersionResponse  "Version retrieved successfully"
// @Failure      500  {object}  ErrorResponse    "Failed to get version from CLI"
// @Failure      405  {object}  ErrorResponse    "Method not allowed"
//
// @Router       /version [get]
func (h *Handler) Version(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	version, err := h.cli.GetVersion()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(VersionResponse{
		Version: version,
	})
}

// Transcripts fetches transcripts for one or more YouTube videos.
//
// @Summary      Fetch Video Transcripts
// @Description  Fetches transcripts for one or more YouTube videos with optional filtering and translation.
// @Description  Supports language filtering, transcript format selection, and exclusion of auto-generated or manually-created transcripts.
// @Description  Can also translate transcripts to other languages if the transcript is marked as translatable.
// @Description  Returns an array of transcript results, one for each requested video.
//
// @Tags         Transcripts
// @Accept       json
// @Produce      json
//
// @Param        request  body      TranscriptRequest  true   "Request containing video IDs and optional filters"
//
// @Success      200      {object}  TranscriptResponse           "Transcripts retrieved successfully"
// @Failure      400      {object}  ErrorResponse                "Invalid request (missing videoIds, invalid JSON, empty body, etc.)"
// @Failure      405      {object}  ErrorResponse                "Method not allowed"
// @Failure      500      {object}  ErrorResponse                "Server error during transcript fetching"
//
// @Router       /transcripts [post]
func (h *Handler) Transcripts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Error:   "Failed to read request body",
		})
		return
	}

	var req TranscriptRequest
	if err := json.Unmarshal(body, &req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Error:   "Invalid JSON in request body: " + err.Error(),
		})
		return
	}

	if len(req.VideoIDs) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Error:   "videoIds is required and must not be empty",
		})
		return
	}

	results, err := h.cli.FetchTranscripts(
		req.VideoIDs,
		req.Languages,
		req.Format,
		req.ExcludeGenerated,
		req.ExcludeManuallyCreated,
		req.Translate,
	)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(TranscriptResponse{
		Success:     true,
		Transcripts: results,
	})
}

// List returns all available transcripts for one or more YouTube videos.
//
// @Summary      List Available Transcripts
// @Description  Lists all available transcripts for one or more YouTube videos, including manually created and auto-generated transcripts with their metadata.
// @Description  Does not fetch the actual transcript content, only information about what transcripts are available.
// @Description  Includes language information, whether transcripts are translatable, and available translation languages.
// @Description  Returns an array of list results, one for each requested video.
//
// @Tags         Transcripts
// @Accept       json
// @Produce      json
//
// @Param        request  body      ListRequest        true   "Request containing video IDs to list transcripts for"
//
// @Success      200      {object}  TranscriptResponse           "Available transcripts listed successfully"
// @Failure      400      {object}  ErrorResponse                "Invalid request (missing videoIds, invalid JSON, empty body, etc.)"
// @Failure      405      {object}  ErrorResponse                "Method not allowed"
// @Failure      500      {object}  ErrorResponse                "Server error during transcript listing"
//
// @Router       /list [post]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Error:   "Failed to read request body",
		})
		return
	}

	var req ListRequest
	if err := json.Unmarshal(body, &req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Error:   "Invalid JSON in request body: " + err.Error(),
		})
		return
	}

	if len(req.VideoIDs) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			Success: false,
			Error:   "videoIds is required and must not be empty",
		})
		return
	}

	results, err := h.cli.ListTranscripts(req.VideoIDs)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(TranscriptResponse{
		Success:     true,
		Transcripts: results,
	})
}
