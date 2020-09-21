package web

import (
	"html/template"
	"net/http"

	goreddit "devisions.org/go-reddit"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Handler struct {
	*chi.Mux
	store goreddit.Store
}

func NewHandler(store goreddit.Store) *Handler {
	h := &Handler{
		Mux:   chi.NewMux(),
		store: store,
	}

	threads := ThreadsHandler{store}
	posts := PostsHandler{store}
	comments := CommentsHandler{store}

	h.Use(middleware.Logger)

	h.Get("/", h.HomeHandler())

	h.Route("/threads", func(r chi.Router) {

		r.Get("/", threads.List())
		r.Get("/new", threads.New())
		r.Post("/", threads.Save())
		r.Get("/{id}", threads.Show())
		r.Post("/{id}/delete", threads.Delete())

		r.Get("/{id}/new", posts.New())
		r.Post("/{id}", posts.Save())
		r.Get("/{threadID}/{postID}", posts.Show())
		r.Get("/{threadID}/{postID}/vote", posts.Vote())

		r.Post("/{threadID}/{postID}", comments.Save())
	})

	h.Get("/comments/{id}/vote", comments.Vote())

	return h
}

func (h *Handler) HomeHandler() http.HandlerFunc {
	tmpl := template.Must(template.ParseFiles("web/templates/layout.html", "web/templates/home.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		ps, err := h.store.GetPosts()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html charset=UTF-8")
		_ = tmpl.Execute(w, struct{ Posts []goreddit.Post }{ps})
	}
}
