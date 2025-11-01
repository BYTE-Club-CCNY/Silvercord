package scores

import (
	"encoding/json"
	"main/routes/utils"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/supabase-community/postgrest-go"
	"github.com/supabase-community/supabase-go"
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
	season := r.URL.Query().Get("season")

	if season == "" {
		season = utils.CURRENT_SEASON
	}

	var rawScores []struct {
		ServerID int64 `json:"server_id"`
		UserID   int64 `json:"user_id"`
		Season   int   `json:"season"`
		Score    int   `json:"score"`
	}

	query := h.client.From(utils.LEETBOARD_SCORES_TABLE).
		Select("server_id, user_id, season, score", "", false).
		Eq("server_id", serverID).
		Eq("user_id", userID).
		Eq("season", season)

	_, err := query.ExecuteTo(&rawScores)
	if err != nil {
		utils.WriteInternalServerErrorResponse(w, "Query unsuccessful")
		return
	}

	if len(rawScores) == 0 {
		seasonInt, _ := strconv.Atoi(season)

		utils.WriteJSONResponse(w, GetScoreResponse{
			ServerID: serverID,
			UserID:   userID,
			Season:   seasonInt,
			Score:    0,
		}, http.StatusOK)
		return
	}

	utils.WriteJSONResponse(w, GetScoreResponse{
		ServerID: strconv.FormatInt(rawScores[0].ServerID, 10),
		UserID:   strconv.FormatInt(rawScores[0].UserID, 10),
		Season:   rawScores[0].Season,
		Score:    rawScores[0].Score,
	}, http.StatusOK)
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

	currentSeason, _ := strconv.Atoi(utils.CURRENT_SEASON)
	serverIDInt, _ := strconv.ParseInt(request.ServerID, 10, 64)
	userIDInt, _ := strconv.ParseInt(request.UserID, 10, 64)

	upsertData := map[string]interface{}{
		"server_id": serverIDInt,
		"user_id":   userIDInt,
		"season":    currentSeason,
		"score":     request.Score,
	}

	// NOTE, this must be within our constraints in the SQL table for the below to work on conflict:
	// ALTER TABLE [table name here]
	// ADD CONSTRAINT [constraint name here]
	// UNIQUE ([rows here]);
	query := h.client.From(utils.LEETBOARD_SCORES_TABLE).
		Upsert(upsertData, "server_id, user_id, season", "", "")
	_, _, err := query.Execute()

	if err != nil {
		utils.WriteInternalServerErrorResponse(w, "Query unsuccessful: "+err.Error())
		return
	}

	utils.WriteSuccessResponse(w, "Success")
}

func (h *Handler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	serverID := r.URL.Query().Get("server_id")
	season := r.URL.Query().Get("season")

	if season == "" {
		season = utils.CURRENT_SEASON
	}

	var scores []LeaderboardRow

	var rawScores []struct {
		UserID int64 `json:"user_id"`
		Score  int   `json:"score"`
	}

	query := h.client.From(utils.LEETBOARD_SCORES_TABLE).
		Select("user_id, score", "", false).
		Eq("server_id", serverID).
		Eq("season", season).
		Order("score", &postgrest.OrderOpts{Ascending: false})

	_, err := query.ExecuteTo(&rawScores)
	if err != nil {
		utils.WriteInternalServerErrorResponse(w, "Query unsuccessful: "+err.Error())
		return
	}

	for _, raw := range rawScores {
		scores = append(scores, LeaderboardRow{
			UserID: strconv.FormatInt(raw.UserID, 10),
			Score:  raw.Score,
		})
	}

	utils.WriteJSONResponse(w, scores, http.StatusOK)
}
