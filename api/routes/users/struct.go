package users

type GetUsernameResponse struct {
	Username string `json:"username"`
}

type RegisterLCUserRequest struct {
	ServerID string `json:"server_id" validate:"required"`
	UserID   string `json:"user_id" validate:"required"`
	Username string `json:"username" validate:"required"`
}

