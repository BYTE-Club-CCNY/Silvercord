package routes

import (
	grpcClient "github.com/BYTE-Club-CCNY/Silvercord/api/grpc"
	"github.com/BYTE-Club-CCNY/Silvercord/api/routes/leetcode_graphql"
	"github.com/BYTE-Club-CCNY/Silvercord/api/routes/problems"
	"github.com/BYTE-Club-CCNY/Silvercord/api/routes/professor"
	"github.com/BYTE-Club-CCNY/Silvercord/api/routes/scores"
	"github.com/BYTE-Club-CCNY/Silvercord/api/routes/users"
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/supabase-community/supabase-go"
)

func SetupRoutes(client *supabase.Client, llmClient *grpcClient.LLMClient) *chi.Mux {
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
		r.Get("/leaderboard", scoresHandler.GetLeaderboard)
	})

	usersHandler := users.NewHandler(client)
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

	professorHandler := professor.NewHandler(llmClient)
	r.Route("/professor", func(r chi.Router) {
		r.Post("/", professorHandler.GetProfessorInfo)
	})

	return r
}
