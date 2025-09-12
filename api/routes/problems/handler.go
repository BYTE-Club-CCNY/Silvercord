package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/supabase-community/supabase-go"
	"net/http"
)

type Handler struct {
	client *supabase.Client
}

func NewHandler(client *supabase.Client) http.Handler {
	h := &Handler{client: client}

	r := chi.NewRouter()
	r.Get("/", h.getProblems)
	return r
}

func (h *Handler) getProblems(w http.ResponseWriter, r *http.Request) {

}
