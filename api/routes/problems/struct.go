package routes

type GetProblemsResponse struct {
	UserID  int    `json:"user_id" validate:"required"`
	Link    string `json:"link" validate:"required"`
	Problem string `json:"problem" validate:"required"`
}
