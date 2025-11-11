package problems

type GetProblemsResponse struct {
	UserID  string `json:"user_id" validate:"required"`
	Link    string `json:"link" validate:"required"`
	Problem string `json:"problem" validate:"required"`
}

type AddProblemRequest struct {
	ServerID string `json:"server_id" validate:"required"`
	UserID   string `json:"user_id" validate:"required"`
	Link     string `json:"link" validate:"required"`
	Problem  string `json:"problem" validate:"required"`
}
