package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/supabase-community/supabase-go"
	routes "main/routes/problems"
	"net/http"
)

func SetupRoutes(client *supabase.Client) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Heartbeat("/ping"))

	// routes to add:
	// add problem(server id, userid, link, problem, table)
	// get problems(serverid, userid, table)
	// get score(serverid, userid, table)
	// get username(serverid, userid, table)
	// updatescore(serverid, userid, score, table)
	// registerlc(serverid, userid, username, table)

	// Root route
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Chi Go Server Running"))
	})

	r.Mount("/problems", routes.NewHandler(client))

	return r
}
