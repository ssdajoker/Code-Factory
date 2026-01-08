package web

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"time"
)

//go:embed static/*
var staticFiles embed.FS

//go:embed templates/*
var templateFiles embed.FS

// Server represents the web server
type Server struct {
	port         int
	contractsDir string
	reportsDir   string
	server       *http.Server
}

// NewServer creates a new web server
func NewServer(port int, contractsDir, reportsDir string) *Server {
	if port == 0 {
		port = 3333
	}
	if contractsDir == "" {
		contractsDir = "contracts"
	}
	if reportsDir == "" {
		reportsDir = "reports"
	}
	return &Server{
		port:         port,
		contractsDir: contractsDir,
		reportsDir:   reportsDir,
	}
}

// Start starts the web server
func (s *Server) Start() error {
	mux := http.NewServeMux()

	// API routes
	h := NewHandlers(s.contractsDir, s.reportsDir)
	mux.HandleFunc("/api/status", h.Status)
	mux.HandleFunc("/api/modes", h.Modes)
	mux.HandleFunc("/api/intake", h.Intake)
	mux.HandleFunc("/api/review", h.Review)
	mux.HandleFunc("/api/rescue", h.Rescue)
	mux.HandleFunc("/api/change-order", h.ChangeOrder)
	mux.HandleFunc("/api/contracts", h.ListContracts)
	mux.HandleFunc("/api/reports", h.ListReports)

	// Static files
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return fmt.Errorf("failed to get static files: %w", err)
	}
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	// Index page
	mux.HandleFunc("/", h.Index)

	s.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", s.port),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	log.Printf("🏭 Code Factory Web UI running at http://localhost:%d", s.port)
	return s.server.ListenAndServe()
}

// Stop gracefully stops the server
func (s *Server) Stop(ctx context.Context) error {
	if s.server != nil {
		return s.server.Shutdown(ctx)
	}
	return nil
}

// Port returns the server port
func (s *Server) Port() int {
	return s.port
}
