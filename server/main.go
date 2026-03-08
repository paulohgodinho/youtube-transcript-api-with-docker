package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// @title           YouTube Transcript API
// @version         1.0.0
// @description     A REST API wrapper for fetching and listing YouTube video transcripts.
// @description     This API provides programmatic access to YouTube transcripts through a simple REST interface,
// @description     leveraging the youtube-transcript-api CLI tool under the hood.
// @description     It supports multiple languages, transcript filtering, and translation capabilities.

// @contact.name    Project Repository
// @contact.url     https://github.com/jdepoix/youtube-transcript-api

// @license.name    MIT License
// @license.url     https://github.com/jdepoix/youtube-transcript-api/blob/master/LICENSE

// @host            localhost:8080
// @BasePath        /

// @schemes         http https

// @externalDocs.description  YouTube Transcript API GitHub Repository
// @externalDocs.url          https://github.com/jdepoix/youtube-transcript-api

// @tag.name        System
// @tag.description System and health check endpoints

// @tag.name        Transcripts
// @tag.description Transcript fetching and listing endpoints

func main() {
	// Parse command line flags
	pythonBin := flag.String("python", DefaultPythonBin, "Path to Python executable")
	port := flag.String("port", "8080", "Port to listen on")
	timeout := flag.Duration("timeout", DefaultTimeout, "Request timeout")
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
	cli, err := NewCLI(*pythonBin, *timeout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Create handler
	handler := NewHandler(cli)

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
