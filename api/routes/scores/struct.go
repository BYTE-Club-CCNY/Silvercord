package scores

type GetScoreResponse struct {
	ServerID string `json:"server_id"`
	UserID   string `json:"user_id"`
	Season   int    `json:"season"`
	Score    int    `json:"score"`
}

type UpdateScoreRequest struct {
	ServerID string `json:"server_id" validate:"required"`
	UserID   string `json:"user_id" validate:"required"`
	Score    int    `json:"score" validate:"required"`
}

type LeaderboardRow struct {
	UserID string `json:"user_id"`
	Score  int    `json:"score"`
}
