package users

type GetUsernameResponse struct {
	Username string `json:"username"`
}

type RegisterLCUserRequest struct {
	ServerID int64  `json:"server_id" validate:"required"`
	UserID   int64  `json:"user_id" validate:"required"`
	Username string `json:"username" validate:"required"`
}
