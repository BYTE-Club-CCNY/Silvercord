package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/supabase-community/supabase-go"
	problems "main/routes/problems"
	"net/http"
)

func SetupRoutes(client *supabase.Client) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Heartbeat("/ping"))

	// routes to add:
	// get score(serverid, userid, table)
	// get username(serverid, userid, table)
	// updatescore(serverid, userid, score, table)
	// registerlc(serverid, userid, username, table)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Chi Go Server Running"))
	})

	problemsHandler := problems.NewHandler(client)
	r.Route("/problems", func(r chi.Router) {
		r.Get("/", problemsHandler.GetProblems)
		r.Post("/", problemsHandler.AddProblem)
	})

	return r
}
