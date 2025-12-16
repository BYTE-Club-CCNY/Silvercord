package problems

import (
	"encoding/json"
	"log"
	"strconv"

	"main/routes/utils"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/supabase-community/supabase-go"
)

type Handler struct {
	client *supabase.Client
}

var validate = validator.New()

func NewHandler(client *supabase.Client) *Handler {
	return &Handler{client: client}
}

func (h *Handler) GetProblems(w http.ResponseWriter, r *http.Request) {
	serverID := r.URL.Query().Get("server_id")
	userID := r.URL.Query().Get("user_id")

	var rawProblems []struct {
		UserID  int64  `json:"user_id"`
		Link    string `json:"link"`
		Problem string `json:"problem"`
	}

	query := h.client.From(utils.LEETBOARD_PROBLEMS_TABLE).
		Select("user_id, link, problem", "", false).
		Eq("server_id", serverID).
		Eq("user_id", userID)
	_, err := query.ExecuteTo(&rawProblems)
	if err != nil {
		log.Printf("[GetProblems] Query failed - serverID: %s, userID: %s, error: %v", serverID, userID, err)
		utils.WriteInternalServerErrorResponse(w, "Query unsuccessful")
		return
	}

	var problems []GetProblemsResponse
	for _, raw := range rawProblems {
		problems = append(problems, GetProblemsResponse{
			UserID:  strconv.FormatInt(raw.UserID, 10),
			Link:    raw.Link,
			Problem: raw.Problem,
		})
	}

	utils.WriteJSONResponse(w, problems, http.StatusOK)
}

func (h *Handler) AddProblem(w http.ResponseWriter, r *http.Request) {
	var request AddProblemRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("[AddProblem] Failed to decode request body - error: %v", err)
		utils.WriteInternalServerErrorResponse(w, "Invalid request body")
		return
	}

	if err := validate.Struct(request); err != nil {
		log.Printf("[AddProblem] Validation failed - error: %v", err)
		utils.WriteInternalServerErrorResponse(w, "Invalid request body")
		return
	}

	serverIDInt, _ := strconv.ParseInt(request.ServerID, 10, 64)
	userIDInt, _ := strconv.ParseInt(request.UserID, 10, 64)

	insertData := map[string]interface{}{
		"server_id": serverIDInt,
		"user_id":   userIDInt,
		"link":      request.Link,
		"problem":   request.Problem,
	}

	query := h.client.From(utils.LEETBOARD_PROBLEMS_TABLE).
		Insert(insertData, true, "server_id,user_id,problem", "representation", "")
	_, _, err := query.Execute()

	if err != nil {
		log.Printf("[AddProblem] Insert failed - serverID: %s, userID: %s, problem: %s, error: %v", request.ServerID, request.UserID, request.Problem, err)
		utils.WriteInternalServerErrorResponse(w, "Query unsuccessful")
		return
	}

	utils.WriteSuccessResponse(w, "Success")
}
