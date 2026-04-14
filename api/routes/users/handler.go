package users

import (
	"database/sql"
	"encoding/json"
	"log"
	"main/routes/utils"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
)

type Handler struct {
	db *sql.DB
}

var validate = validator.New()

func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) GetUsername(w http.ResponseWriter, r *http.Request) {
	serverID := r.URL.Query().Get("server_id")
	userID := r.URL.Query().Get("user_id")

	var username string
	err := h.db.QueryRow(
		"SELECT leetcode_username FROM leetboard_username WHERE server_id=? AND user_id=?",
		serverID, userID,
	).Scan(&username)

	if err == sql.ErrNoRows {
		utils.WriteJSONResponse(w, GetUsernameResponse{
			Username: "",
		}, http.StatusOK)
		return
	}
	if err != nil {
		log.Printf("[GetUsername] Query failed - serverID: %s, userID: %s, error: %v", serverID, userID, err)
		utils.WriteInternalServerErrorResponse(w, "Query unsuccessful")
		return
	}

	utils.WriteJSONResponse(w, GetUsernameResponse{
		Username: username,
	}, http.StatusOK)
}

func (h *Handler) RegisterLCUser(w http.ResponseWriter, r *http.Request) {
	var request RegisterLCUserRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("[RegisterLCUser] Failed to decode request body - error: %v", err)
		utils.WriteInternalServerErrorResponse(w, "Invalid request body")
		return
	}

	if err := validate.Struct(request); err != nil {
		log.Printf("[RegisterLCUser] Validation failed - error: %v", err)
		utils.WriteInternalServerErrorResponse(w, "Invalid request body")
		return
	}

	serverIDInt, _ := strconv.ParseInt(request.ServerID, 10, 64)
	userIDInt, _ := strconv.ParseInt(request.UserID, 10, 64)

	_, err := h.db.Exec(
		"INSERT OR REPLACE INTO leetboard_username (server_id, user_id, leetcode_username) VALUES (?, ?, ?)",
		serverIDInt, userIDInt, request.Username,
	)
	if err != nil {
		log.Printf("[RegisterLCUser] Upsert failed - serverID: %s, userID: %s, username: %s, error: %v", request.ServerID, request.UserID, request.Username, err)
		utils.WriteInternalServerErrorResponse(w, "Query unsuccessful")
		return
	}

	utils.WriteSuccessResponse(w, "Success")
}
