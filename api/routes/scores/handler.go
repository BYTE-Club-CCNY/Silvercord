package scores

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/supabase-community/supabase-go"
	"main/routes/utils"
	"net/http"
)

type Handler struct {
	client *supabase.Client
}

var validate = validator.New()

func NewHandler(client *supabase.Client) *Handler {
	return &Handler{client: client}
}

func (h *Handler) GetScore(w http.ResponseWriter, r *http.Request) {
	serverID := r.URL.Query().Get("server_id")
	userID := r.URL.Query().Get("user_id")
	
	var scores []GetScoreResponse

	query := h.client.From(utils.LEETBOARD_SCORES_TABLE).
		Select("server_id, user_id, season, score", "", false).
		Eq("server_id", serverID).
		Eq("user_id", userID).
		Eq("season", utils.CURRENT_SEASON)
	
	_, err := query.ExecuteTo(&scores)
	if err != nil {
		utils.WriteInternalServerErrorResponse(w, "Query unsuccessful")
		return
	}

	if len(scores) == 0 {
		utils.WriteJSONResponse(w, GetScoreResponse{
			ServerID: serverID,
			UserID:   userID,
			Season:   utils.CURRENT_SEASON,
			Score:    0,
		}, http.StatusOK)
		return
	}

	utils.WriteJSONResponse(w, scores[0], http.StatusOK)
}

func (h *Handler) UpdateScore(w http.ResponseWriter, r *http.Request) {
	var request UpdateScoreRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.WriteInternalServerErrorResponse(w, "Invalid request body")
		return
	}

	if err := validate.Struct(request); err != nil {
		utils.WriteInternalServerErrorResponse(w, "Invalid request body")
		return
	}

	upsertData := map[string]interface{}{
		"server_id": request.ServerID,
		"user_id":   request.UserID,
		"season":    utils.CURRENT_SEASON,
		"score":     request.Score,
	}

	query := h.client.From(utils.LEETBOARD_SCORES_TABLE).
		Upsert(upsertData, "", "", "representation", "")
	_, _, err := query.Execute()

	if err != nil {
		utils.WriteInternalServerErrorResponse(w, "Query unsuccessful")
		return
	}

	utils.WriteSuccessResponse(w, "Success")
}
