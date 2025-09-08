package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/supabase-community/supabase-go"
	"net/http"
)

func SetupRoutes(client *supabase.Client) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Heartbeat("/ping"))

	//r.Mount("/api/users", UserRoutes(client))
	//r.Mount("/api/posts", PostRoutes(client))
	//r.Mount("/api/auth", AuthRoutes(client))

	// Root route
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Chi Go Server Running"))
	})

	return r
}
