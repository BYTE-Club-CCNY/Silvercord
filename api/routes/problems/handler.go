package problems

import (
	"database/sql"
	"encoding/json"
	"log"
	"strconv"

	"main/routes/utils"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Handler struct {
	db *sql.DB
}

var validate = validator.New()

func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) GetProblems(w http.ResponseWriter, r *http.Request) {
	serverID := r.URL.Query().Get("server_id")
	userID := r.URL.Query().Get("user_id")

	rows, err := h.db.Query(
		"SELECT user_id, link, problem FROM leetboard_problems WHERE server_id=? AND user_id=?",
		serverID, userID,
	)
	if err != nil {
		log.Printf("[GetProblems] Query failed - serverID: %s, userID: %s, error: %v", serverID, userID, err)
		utils.WriteInternalServerErrorResponse(w, "Query unsuccessful")
		return
	}
	defer rows.Close()

	var problems []GetProblemsResponse
	for rows.Next() {
		var userIDInt int64
		var link, problem string
		if err := rows.Scan(&userIDInt, &link, &problem); err != nil {
			log.Printf("[GetProblems] Row scan failed - error: %v", err)
			utils.WriteInternalServerErrorResponse(w, "Query unsuccessful")
			return
		}
		problems = append(problems, GetProblemsResponse{
			UserID:  strconv.FormatInt(userIDInt, 10),
			Link:    link,
			Problem: problem,
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

	_, err := h.db.Exec(
		"INSERT OR IGNORE INTO leetboard_problems (server_id, user_id, link, problem) VALUES (?, ?, ?, ?)",
		serverIDInt, userIDInt, request.Link, request.Problem,
	)
	if err != nil {
		log.Printf("[AddProblem] Insert failed - serverID: %s, userID: %s, problem: %s, error: %v", request.ServerID, request.UserID, request.Problem, err)
		utils.WriteInternalServerErrorResponse(w, "Query unsuccessful")
		return
	}

	utils.WriteSuccessResponse(w, "Success")
}
