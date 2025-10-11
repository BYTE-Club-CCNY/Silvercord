package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/supabase-community/supabase-go"
	problems "main/routes/problems"
	scores "main/routes/scores"
	users "main/routes/users"
	"net/http"
)

func SetupRoutes(client *supabase.Client) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Heartbeat("/ping"))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Chi Go Server Running"))
	})

	problemsHandler := problems.NewHandler(client)
	r.Route("/problems", func(r chi.Router) {
		r.Get("/", problemsHandler.GetProblems)
		r.Post("/", problemsHandler.AddProblem)
	})

	scoresHandler := scores.NewHandler(client)
	r.Route("/scores", func(r chi.Router) {
		r.Get("/", scoresHandler.GetScore)
		r.Post("/", scoresHandler.UpdateScore)
	})

	usersHandler := users.NewHandler(client)
	r.Route("/users", func(r chi.Router) {
		r.Get("/username", usersHandler.GetUsername)
		r.Post("/register", usersHandler.RegisterLCUser)
	})

	return r
}
