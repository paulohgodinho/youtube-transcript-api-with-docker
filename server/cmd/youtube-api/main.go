package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/youtube-transcript-api/server/internal/api"
	"github.com/youtube-transcript-api/server/internal/python"
)

func main() {
	// Parse command line flags
	pythonBin := flag.String("python", python.DefaultPythonBin, "Path to Python executable")
	port := flag.String("port", "8080", "Port to listen on")
	timeout := flag.Duration("timeout", python.DefaultTimeout, "Request timeout")
	flag.Parse()

	// Override with environment variables if set
	if envPython := os.Getenv("PYTHON_BIN"); envPython != "" {
		*pythonBin = envPython
	}
	if envPort := os.Getenv("SERVER_PORT"); envPort != "" {
		*port = envPort
	}
	if envTimeout := os.Getenv("REQUEST_TIMEOUT"); envTimeout != "" {
		if t, err := time.ParseDuration(envTimeout); err == nil {
			*timeout = t
		}
	}

	// Create Python CLI wrapper
	cli, err := python.NewCLI(*pythonBin, *timeout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Create handler
	handler := api.NewHandler(cli)

	// Set up routes
	http.HandleFunc("/health", handler.Health)
	http.HandleFunc("/version", handler.Version)
	http.HandleFunc("/transcripts", handler.Transcripts)
	http.HandleFunc("/list", handler.List)

	// Start server
	addr := ":" + *port
	log.Printf("Starting YouTube Transcript API server on %s", addr)
	log.Printf("Python executable: %s", *pythonBin)
	log.Printf("Request timeout: %v", *timeout)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
