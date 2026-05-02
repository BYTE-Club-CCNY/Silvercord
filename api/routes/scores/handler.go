package scores

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

func (h *Handler) GetScore(w http.ResponseWriter, r *http.Request) {
	serverID := r.URL.Query().Get("server_id")
	userID := r.URL.Query().Get("user_id")
	season := r.URL.Query().Get("season")

	if season == "" {
		season = utils.CURRENT_SEASON
	}

	var srvID int64
	var uID int64
	var seasonVal int
	var score int

	err := h.db.QueryRow(
		"SELECT server_id, user_id, season, score FROM leetboard_scores WHERE server_id=? AND user_id=? AND season=?",
		serverID, userID, season,
	).Scan(&srvID, &uID, &seasonVal, &score)

	if err == sql.ErrNoRows {
		seasonInt, _ := strconv.Atoi(season)
		utils.WriteJSONResponse(w, GetScoreResponse{
			ServerID: serverID,
			UserID:   userID,
			Season:   seasonInt,
			Score:    0,
		}, http.StatusOK)
		return
	}
	if err != nil {
		log.Printf("[GetScore] Query failed - serverID: %s, userID: %s, season: %s, error: %v", serverID, userID, season, err)
		utils.WriteInternalServerErrorResponse(w, "Query unsuccessful")
		return
	}

	utils.WriteJSONResponse(w, GetScoreResponse{
		ServerID: strconv.FormatInt(srvID, 10),
		UserID:   strconv.FormatInt(uID, 10),
		Season:   seasonVal,
		Score:    score,
	}, http.StatusOK)
}

func (h *Handler) UpdateScore(w http.ResponseWriter, r *http.Request) {
	var request UpdateScoreRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("[UpdateScore] Failed to decode request body - error: %v", err)
		utils.WriteInternalServerErrorResponse(w, "Invalid request body")
		return
	}

	if err := validate.Struct(request); err != nil {
		log.Printf("[UpdateScore] Validation failed - error: %v", err)
		utils.WriteInternalServerErrorResponse(w, "Invalid request body")
		return
	}

	currentSeason, _ := strconv.Atoi(utils.CURRENT_SEASON)
	serverIDInt, _ := strconv.ParseInt(request.ServerID, 10, 64)
	userIDInt, _ := strconv.ParseInt(request.UserID, 10, 64)

	_, err := h.db.Exec(
		"INSERT INTO leetboard_scores (server_id, user_id, season, score) VALUES (?, ?, ?, ?) ON CONFLICT(server_id, user_id, season) DO UPDATE SET score=?",
		serverIDInt, userIDInt, currentSeason, request.Score, request.Score,
	)
	if err != nil {
		log.Printf("[UpdateScore] Upsert failed - serverID: %s, userID: %s, score: %d, error: %v", request.ServerID, request.UserID, request.Score, err)
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

	rows, err := h.db.Query(
		"SELECT user_id, score FROM leetboard_scores WHERE server_id=? AND season=? ORDER BY score DESC",
		serverID, season,
	)
	if err != nil {
		log.Printf("[GetLeaderboard] Query failed - serverID: %s, season: %s, error: %v", serverID, season, err)
		utils.WriteInternalServerErrorResponse(w, "Query unsuccessful: "+err.Error())
		return
	}
	defer rows.Close()

	var scores []LeaderboardRow
	for rows.Next() {
		var userID int64
		var score int
		if err := rows.Scan(&userID, &score); err != nil {
			log.Printf("[GetLeaderboard] Row scan failed - error: %v", err)
			utils.WriteInternalServerErrorResponse(w, "Query unsuccessful")
			return
		}
		scores = append(scores, LeaderboardRow{
			UserID: strconv.FormatInt(userID, 10),
			Score:  score,
		})
	}

	utils.WriteJSONResponse(w, scores, http.StatusOK)
}
