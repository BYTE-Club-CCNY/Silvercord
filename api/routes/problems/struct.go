package problems

import "github.com/google/uuid"

type GetProblemsResponse struct {
	UserID  int    `json:"user_id" validate:"required"`
	Link    string `json:"link" validate:"required"`
	Problem string `json:"problem" validate:"required"`
}

type AddProblemRequest struct {
	ID       uuid.UUID `json:"id" validate:"required"`
	ServerID int       `json:"server_id" validate:"required"`
	UserID   int       `json:"user_id" validate:"required"`
	Link     string    `json:"link" validate:"required"`
	Problem  string    `json:"problem" validate:"required"`
}
