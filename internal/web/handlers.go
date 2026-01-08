package web

import (
	"context"
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/ssdajoker/Code-Factory/internal/modes"
)

// Handlers contains HTTP handlers
type Handlers struct {
	contractsDir string
	reportsDir   string
}

// NewHandlers creates new handlers
func NewHandlers(contractsDir, reportsDir string) *Handlers {
	return &Handlers{
		contractsDir: contractsDir,
		reportsDir:   reportsDir,
	}
}

// StatusResponse represents the status response
type StatusResponse struct {
	Status       string `json:"status"`
	Version      string `json:"version"`
	ContractsDir string `json:"contracts_dir"`
	ReportsDir   string `json:"reports_dir"`
}

// Status returns factory status
func (h *Handlers) Status(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	resp := StatusResponse{
		Status:       "running",
		Version:      "0.3.0",
		ContractsDir: h.contractsDir,
		ReportsDir:   h.reportsDir,
	}
	jsonResponse(w, resp)
}

// ModesResponse represents available modes
type ModesResponse struct {
	Modes []ModeInfo `json:"modes"`
}

// ModeInfo describes a mode
type ModeInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Endpoint    string `json:"endpoint"`
}

// Modes returns available modes
func (h *Handlers) Modes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	resp := ModesResponse{
		Modes: []ModeInfo{
			{Name: "INTAKE", Description: "Capture project vision and generate spec", Endpoint: "/api/intake"},
			{Name: "REVIEW", Description: "Check code compliance with spec", Endpoint: "/api/review"},
			{Name: "RESCUE", Description: "Reverse engineer spec from code", Endpoint: "/api/rescue"},
			{Name: "CHANGE_ORDER", Description: "Track and manage spec drift", Endpoint: "/api/change-order"},
		},
	}
	jsonResponse(w, resp)
}

// IntakeRequest represents intake request
type IntakeRequest struct {
	ProjectName          string `json:"project_name"`
	Description          string `json:"description"`
	TargetUsers          string `json:"target_users"`
	CoreFeatures         string `json:"core_features"`
	TechnicalConstraints string `json:"technical_constraints"`
	SuccessCriteria      string `json:"success_criteria"`
}

// Intake handles intake mode
func (h *Handlers) Intake(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req IntakeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	intake := modes.NewIntakeMode(nil, h.contractsDir)
	intake.SetStepValue(req.ProjectName)
	intake.NextStep()
	intake.SetStepValue(req.Description)
	intake.NextStep()
	intake.SetStepValue(req.TargetUsers)
	intake.NextStep()
	intake.SetStepValue(req.CoreFeatures)
	intake.NextStep()
	intake.SetStepValue(req.TechnicalConstraints)
	intake.NextStep()
	intake.SetStepValue(req.SuccessCriteria)
	intake.NextStep()

	spec, err := intake.GenerateSpec(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	path, err := intake.SaveSpec()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, map[string]interface{}{
		"success": true,
		"spec":    spec,
		"path":    path,
	})
}

// ReviewRequest represents review request
type ReviewRequest struct {
	SpecFile  string   `json:"spec_file"`
	CodePaths []string `json:"code_paths"`
}

// Review handles review mode
func (h *Handlers) Review(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	review := modes.NewReviewMode(nil, h.reportsDir)
	review.SetSpecFile(req.SpecFile)
	review.SetCodePaths(req.CodePaths)

	result, err := review.RunReview(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	path, err := review.SaveReport()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, map[string]interface{}{
		"success":          true,
		"compliance_score": result.ComplianceScore,
		"report":           result.FullReport,
		"path":             path,
	})
}

// RescueRequest represents rescue request
type RescueRequest struct {
	CodebasePath string `json:"codebase_path"`
}

// Rescue handles rescue mode
func (h *Handlers) Rescue(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RescueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	rescue := modes.NewRescueMode(nil, h.contractsDir, h.reportsDir)
	rescue.SetCodebasePath(req.CodebasePath)

	result, err := rescue.ScanCodebase(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	specPath, reportPath, err := rescue.SaveResults()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, map[string]interface{}{
		"success":       true,
		"files_scanned": result.FilesScanned,
		"spec":          result.InferredSpec,
		"spec_path":     specPath,
		"report_path":   reportPath,
	})
}

// ChangeOrderRequest represents change order request
type ChangeOrderRequest struct {
	SpecFile     string `json:"spec_file"`
	CodebasePath string `json:"codebase_path"`
}

// ChangeOrder handles change order mode
func (h *Handlers) ChangeOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ChangeOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	co := modes.NewChangeOrderMode(nil, h.contractsDir)
	co.SetSpecFile(req.SpecFile)
	co.SetCodebasePath(req.CodebasePath)

	result, err := co.DetectDrift(context.Background())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	path, err := co.SaveChangeOrder()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, map[string]interface{}{
		"success": true,
		"changes": result.Changes,
		"report":  result.FullReport,
		"path":    path,
	})
}

// ListContracts lists contract files
func (h *Handlers) ListContracts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	files := listMarkdownFiles(h.contractsDir)
	jsonResponse(w, map[string]interface{}{
		"contracts": files,
	})
}

// ListReports lists report files
func (h *Handlers) ListReports(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	files := listMarkdownFiles(h.reportsDir)
	jsonResponse(w, map[string]interface{}{
		"reports": files,
	})
}

// Index serves the main page
func (h *Handlers) Index(w http.ResponseWriter, r *http.Request) {
	tmplContent, err := templateFiles.ReadFile("templates/index.html")
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.New("index").Parse(string(tmplContent))
	if err != nil {
		http.Error(w, "Template parse error", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func listMarkdownFiles(dir string) []string {
	var files []string
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".md") {
			files = append(files, path)
		}
		return nil
	})
	return files
}
