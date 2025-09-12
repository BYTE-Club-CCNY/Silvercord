package problems

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/supabase-community/supabase-go"
	"main/routes/utils"
	"net/http"
	"strings"
)

type Handler struct {
	client *supabase.Client
}

var validate = validator.New()

func NewHandler(client *supabase.Client) *Handler {
	return &Handler{client: client}
}

func (h *Handler) GetProblems(w http.ResponseWriter, r *http.Request) {
	serverID := r.URL.Query().Get("server_id")
	userID := r.URL.Query().Get("user_id")
	var problems []GetProblemsResponse

	query := h.client.From(utils.LEETBOARD_PROBLEMS_TABLE).
		Select("user_id, link, problem", "", false).
		Eq("server_id", serverID).
		Eq("user_id", userID)
	_, err := query.ExecuteTo(&problems)
	if err != nil {
		utils.WriteInternalServerErrorResponse(w, "Query unsuccessful")
		return
	}
	utils.WriteJSONResponse(w, problems, http.StatusOK)
}

func (h *Handler) AddProblem(w http.ResponseWriter, r *http.Request) {
	var request AddProblemRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.WriteInternalServerErrorResponse(w, "Invalid request body")
		return
	}

	if err := validate.Struct(request); err != nil {
		utils.WriteInternalServerErrorResponse(w, "Invalid request body")
		return
	}

	insertData := map[string]interface{}{
		"id":        request.ID,
		"server_id": request.ServerID,
		"user_id":   request.UserID,
		"link":      request.Link,
		"problem":   request.Problem,
	}

	query := h.client.From(utils.LEETBOARD_PROBLEMS_TABLE).
		Insert(insertData, false, "server_id,user_id,problem", "representation", "")
	_, _, err := query.Execute()

	if err != nil {
		if strings.Contains(err.Error(), "23505") {
			utils.WriteConflictResponse(w, "User has already submitted problem in this server")
			return
		}
		utils.WriteInternalServerErrorResponse(w, "Query unsuccessful")
		return
	}

	utils.WriteSuccessResponse(w, "Success")
}
