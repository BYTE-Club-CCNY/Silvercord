package professor

import (
	"encoding/json"
	"net/http"

	"github.com/BYTE-Club-CCNY/Silvercord/api/routes/utils"
	grpcClient "github.com/BYTE-Club-CCNY/Silvercord/api/grpc"
)

type Handler struct {
	llmClient *grpcClient.LLMClient
}

// NewHandler creates a new professor handler with a gRPC LLM client
func NewHandler(llmClient *grpcClient.LLMClient) *Handler {
	return &Handler{llmClient: llmClient}
}

// ProfessorRequest represents the incoming request for professor information
type ProfessorRequest struct {
	ProfessorName string `json:"professor_name"`
	UserID        string `json:"user_id,omitempty"`
}

// ProfessorResponse represents the response with professor information
type ProfessorResponse struct {
	Name     string `json:"name"`
	Link     string `json:"link"`
	Response string `json:"response"`
}

// GetProfessorInfo handles POST requests for professor information
func (h *Handler) GetProfessorInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJSONResponse(w, utils.NewConflictMessage("Method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	var req ProfessorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONResponse(w, utils.NewConflictMessage("Invalid request body"), http.StatusBadRequest)
		return
	}

	if req.ProfessorName == "" {
		utils.WriteJSONResponse(w, utils.NewConflictMessage("professor_name is required"), http.StatusBadRequest)
		return
	}

	// Default user ID if not provided
	if req.UserID == "" {
		req.UserID = "anonymous"
	}

	// Call the LLM service via gRPC
	responseText, err := h.llmClient.ProcessQuery(r.Context(), req.UserID, "professor", req.ProfessorName)
	if err != nil {
		utils.WriteInternalServerErrorResponse(w, "Failed to get professor information: "+err.Error())
		return
	}

	// Parse the JSON response from the LLM service
	var profResp ProfessorResponse
	if err := json.Unmarshal([]byte(responseText), &profResp); err != nil {
		utils.WriteInternalServerErrorResponse(w, "Failed to parse LLM response: "+err.Error())
		return
	}

	// Return the professor information
	utils.WriteJSONResponse(w, profResp, http.StatusOK)
}
