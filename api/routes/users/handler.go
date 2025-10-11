package users

import (
	"encoding/json"
	"main/routes/utils"
	"net/http"
	"strconv"

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

func (h *Handler) GetUsername(w http.ResponseWriter, r *http.Request) {
	serverID := r.URL.Query().Get("server_id")
	userID := r.URL.Query().Get("user_id")

	var users []struct {
		LeetcodeUsername string `json:"leetcode_username"`
	}

	query := h.client.From(utils.LEETBOARD_USERNAME_TABLE).
		Select("leetcode_username", "", false).
		Eq("server_id", serverID).
		Eq("user_id", userID)

	_, err := query.ExecuteTo(&users)
	if err != nil {
		utils.WriteInternalServerErrorResponse(w, "Query unsuccessful")
		return
	}

	if len(users) == 0 {
		utils.WriteJSONResponse(w, GetUsernameResponse{
			Username: "",
		}, http.StatusOK)
		return
	}

	utils.WriteJSONResponse(w, GetUsernameResponse{
		Username: users[0].LeetcodeUsername,
	}, http.StatusOK)
}

func (h *Handler) RegisterLCUser(w http.ResponseWriter, r *http.Request) {
	var request RegisterLCUserRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.WriteInternalServerErrorResponse(w, "Invalid request body")
		return
	}

	if err := validate.Struct(request); err != nil {
		utils.WriteInternalServerErrorResponse(w, "Invalid request body")
		return
	}

	serverIDInt, _ := strconv.ParseInt(request.ServerID, 10, 64)
	userIDInt, _ := strconv.ParseInt(request.UserID, 10, 64)

	upsertData := map[string]interface{}{
		"server_id":         serverIDInt,
		"user_id":           userIDInt,
		"leetcode_username": request.Username,
	}

	query := h.client.From(utils.LEETBOARD_USERNAME_TABLE).
		Upsert(upsertData, "", "representation", "")
	_, _, err := query.Execute()

	if err != nil {
		utils.WriteInternalServerErrorResponse(w, "Query unsuccessful")
		return
	}

	utils.WriteSuccessResponse(w, "Success")
}
