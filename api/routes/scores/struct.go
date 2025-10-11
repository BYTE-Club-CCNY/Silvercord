package scores

type GetScoreResponse struct {
	ServerID string `json:"server_id"`
	UserID   string `json:"user_id"`
	Season   string `json:"season"`
	Score    int    `json:"score"`
}

type UpdateScoreRequest struct {
	ServerID string `json:"server_id" validate:"required"`
	UserID   string `json:"user_id" validate:"required"`
	Score    int    `json:"score" validate:"required"`
}

