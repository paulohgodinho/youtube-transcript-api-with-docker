package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/youtube-transcript-api/server/internal/models"
	"github.com/youtube-transcript-api/server/internal/python"
)

// Handler holds dependencies for HTTP handlers
type Handler struct {
	cli *python.CLI
}

// NewHandler creates a new HTTP handler
func NewHandler(cli *python.CLI) *Handler {
	return &Handler{
		cli: cli,
	}
}

// Health handles GET /health
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.HealthResponse{
		Status: "ok",
	})
}

// Version handles GET /version
func (h *Handler) Version(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	version, err := h.cli.GetVersion()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.VersionResponse{
		Version: version,
	})
}

// Transcripts handles POST /transcripts
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
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Success: false,
			Error:   "Failed to read request body",
		})
		return
	}

	var req models.TranscriptRequest
	if err := json.Unmarshal(body, &req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Success: false,
			Error:   "Invalid JSON in request body: " + err.Error(),
		})
		return
	}

	if len(req.VideoIDs) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
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

	transcripts := make([]models.TranscriptData, 0, len(results))
	for videoID, data := range results {
		transcripts = append(transcripts, models.TranscriptData{
			VideoID: videoID,
			Data:    data,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.TranscriptResponse{
		Success:     true,
		Transcripts: transcripts,
	})
}

// List handles POST /list
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
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Success: false,
			Error:   "Failed to read request body",
		})
		return
	}

	var req models.ListRequest
	if err := json.Unmarshal(body, &req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Success: false,
			Error:   "Invalid JSON in request body: " + err.Error(),
		})
		return
	}

	if len(req.VideoIDs) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.ErrorResponse{
			Success: false,
			Error:   "videoIds is required and must not be empty",
		})
		return
	}

	results, err := h.cli.ListTranscripts(req.VideoIDs)

	transcripts := make([]models.TranscriptData, 0, len(results))
	for videoID, data := range results {
		transcripts = append(transcripts, models.TranscriptData{
			VideoID: videoID,
			Data:    data,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.TranscriptResponse{
		Success:     true,
		Transcripts: transcripts,
	})
}
