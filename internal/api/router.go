package api

import (
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	chi.Router
}

func NewRouter() *Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	return &Router{r}
}
