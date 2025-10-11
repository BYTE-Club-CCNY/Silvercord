package scores

type GetScoreResponse struct {
	ServerID int64 `json:"server_id"`
	UserID   int64 `json:"user_id"`
	Season   int   `json:"season"`
	Score    int   `json:"score"`
}

type UpdateScoreRequest struct {
	ServerID int64 `json:"server_id" validate:"required"`
	UserID   int64 `json:"user_id" validate:"required"`
	Score    int   `json:"score" validate:"required"`
}
