package break

import (
	"encoding/json"
	"net/http"

	"github.com/BYTE-Club-CCNY/Silvercord/api/routes/utils"
	grpcClient "github.com/BYTE-Club-CCNY/Silvercord/api/grpc"
)

type Handler struct {
	llmClient *grpcClient.LLMClient
}

func NewHandler(llmClient *grpcClient.LLMClient) *Handler {
	return &Handler{llmClient: llmClient}
}

func (h *Handler) GetBreakInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJSONResponse(w, utils.NewConflictMessage("Method not allowed"), http.StatusMethodNotAllowed)
		return
	}

	var req BreakRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSONResponse(w, utils.NewConflictMessage("Invalid request body"), http.StatusBadRequest)
		return
	}

	if req.Query == "" {
		utils.WriteJSONResponse(w, utils.NewConflictMessage("query is required"), http.StatusBadRequest)
		return
	}

	if req.UserID == "" {
		req.UserID = "anonymous"
	}

	responseText, err := h.llmClient.ProcessQuery(r.Context(), req.UserID, "break", req.Query)
	if err != nil {
		utils.WriteInternalServerErrorResponse(w, "Failed to get break information: "+err.Error())
		return
	}

	var breakResp BreakResponse
	if err := json.Unmarshal([]byte(responseText), &breakResp); err != nil {
		utils.WriteInternalServerErrorResponse(w, "Failed to parse LLM response: "+err.Error())
		return
	}

	utils.WriteJSONResponse(w, breakResp, http.StatusOK)
}

