package routes

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/supabase-community/supabase-go"
	"log"
	"net/http"
)

type Handler struct {
	client *supabase.Client
}

const LEETBOARD_PROBLEMS_TABLE = "leetboard_problems"

func NewHandler(client *supabase.Client) http.Handler {
	h := &Handler{client: client}

	r := chi.NewRouter()
	r.Get("/", h.getProblems)
	return r
}

func (h *Handler) getProblems(w http.ResponseWriter, r *http.Request) {
	serverID := r.URL.Query().Get("server_id")
	userID := r.URL.Query().Get("user_id")
	var problems []GetProblemsResponse

	queryBuilder := h.client.From(LEETBOARD_PROBLEMS_TABLE)
	querySelect := queryBuilder.
		Select("user_id, link, problem", "1", false).
		Eq("server_id", serverID).
		Eq("user_id", userID)
	_, err := querySelect.ExecuteTo(&problems)
	if err != nil {
		http.Error(w, "Error getting problems", http.StatusInternalServerError)
		log.Fatal("Error: ", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(problems)
	if err != nil {
		http.Error(w, "Error encoding", http.StatusInternalServerError)
		return
	}
}
