package controller

import (
	"encoding/json"
	"net/http"
	"note-service/internal/models"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// @Summary Create note
// @Security ApiKeyAuth
// @Tags notes
// @Description create note
// @ID create-note
// @Accept  json
// @Produce  json
// @Param input body models.Note true "note info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /api/notes [post]
func (h *Handler) createNote(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserId(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var input models.Note
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.service.Create(userId, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id": id,
	})
}

type getAllNotesResponse struct {
	Data []models.Note `json:"data"`
}

// @Summary Get all notes
// @Security ApiKeyAuth
// @Tags notes
// @Description get all notes
// @ID get-all-notes
// @Accept  json
// @Produce  json
// @Success 200 {object} getAllNotesResponse
// @Failure 400,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /api/notes [get]
func (h *Handler) getAllNotes(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserId(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	notes, err := h.service.GetAll(userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(getAllNotesResponse{
		Data: notes,
	})
}

// @Summary Get Note By Id
// @Security ApiKeyAuth
// @Tags notes
// @Description get note by id
// @ID get-note-by-id
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Note
// @Failure 400,404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Failure default {object} ErrorResponse
// @Router /api/notes/:id [get]
func (h *Handler) getNoteById(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserId(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id param", http.StatusBadRequest)
		return
	}

	note, err := h.service.GetById(userId, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(note)
}

type statusResponse struct {
	Status string `json:"status"`
}

// @Summary Update a note
// @Security ApiKeyAuth
// @Tags notes
// @Description Update a specific note by its ID
// @ID update-note
// @Accept json
// @Produce json
// @Param id path int true "Note ID"
// @Param input body models.UpdateNoteInput true "Updated note information"
// @Success 200 {object} statusResponse "Successful operation"
// @Failure 400 {object} ErrorResponse "Invalid input"
// @Failure 404 {object} ErrorResponse "Note not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/notes/{id} [put]
func (h *Handler) updateNote(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserId(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id param", http.StatusBadRequest)
		return
	}

	var input models.UpdateNoteInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.Update(userId, id, input); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(statusResponse{"ok"})
}

// @Summary Delete a note
// @Security ApiKeyAuth
// @Tags notes
// @Description Delete a specific note by its ID
// @ID delete-note
// @Accept json
// @Produce json
// @Param id path int true "Note ID"
// @Success 200 {object} statusResponse "Successful operation"
// @Failure 400 {object} ErrorResponse "Invalid ID parameter"
// @Failure 404 {object} ErrorResponse "Note not found"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /api/notes/{id} [delete]
func (h *Handler) deleteNote(w http.ResponseWriter, r *http.Request) {
	userId, err := getUserId(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id param", http.StatusBadRequest)
		return
	}

	err = h.service.Delete(userId, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(statusResponse{"ok"})
}
