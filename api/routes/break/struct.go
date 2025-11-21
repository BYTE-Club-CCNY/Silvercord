package break

type BreakRequest struct {
	Query  string `json:"query"`
	UserID string `json:"user_id,omitempty"`
}

type BreakResponse struct {
	Name     string `json:"name"`
	Link     string `json:"link"`
	Response string `json:"response"`
}
