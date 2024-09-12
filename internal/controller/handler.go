package controller

import (
	"note-service/internal/service"

	_ "note-service/docs"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/sign-up", h.signUp)
		r.Post("/sign-in", h.signIn)
	})

	r.Route("/api", func(r chi.Router) {

		r.With(h.userIdentity).Route("/notes", func(r chi.Router) {
			r.Post("/", h.createNote)
			r.Get("/", h.getAllNotes)
			r.Get("/{id}", h.getNoteById)
			r.Put("/{id}", h.updateNote)
			r.Delete("/{id}", h.deleteNote)
		})
	})

	return r
}
