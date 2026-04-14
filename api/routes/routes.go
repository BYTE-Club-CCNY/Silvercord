package routes

import (
	"database/sql"
	leetcode_graphql "main/routes/leetcode_graphql"
	problems "main/routes/problems"
	scores "main/routes/scores"
	users "main/routes/users"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRoutes(db *sql.DB) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Heartbeat("/ping"))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Chi Go Server Running"))
	})

	problemsHandler := problems.NewHandler(db)
	r.Route("/problems", func(r chi.Router) {
		r.Get("/", problemsHandler.GetProblems)
		r.Post("/", problemsHandler.AddProblem)
	})

	scoresHandler := scores.NewHandler(db)
	r.Route("/scores", func(r chi.Router) {
		r.Get("/", scoresHandler.GetScore)
		r.Post("/", scoresHandler.UpdateScore)
		r.Get("/leaderboard", scoresHandler.GetLeaderboard)
	})

	usersHandler := users.NewHandler(db)
	r.Route("/users", func(r chi.Router) {
		r.Get("/username", usersHandler.GetUsername)
		r.Post("/register", usersHandler.RegisterLCUser)
	})

	leetcodeHandler := leetcode_graphql.NewHandler()
	r.Route("/leetcode", func(r chi.Router) {
		r.Get("/user", leetcodeHandler.GetOnlineUsername)
		r.Get("/difficulty", leetcodeHandler.GetDifficulty)
		r.Get("/extract-problem", leetcodeHandler.ExtractProblem)
	})

	return r
}
