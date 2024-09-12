package repository

import (
	"note-service/internal/models"

	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GetUser(username, password string) (models.User, error)
}

type Notes interface {
	Create(userId int, note models.Note) (int, error)
	GetAll(userId int) ([]models.Note, error)
	GetById(userId, noteId int) (models.Note, error)
	Delete(userId, noteId int) error
	Update(userId, noteId int, input models.UpdateNoteInput) error
}

type Repository struct {
	Authorization
	Notes
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Notes:         NewNotesPostgres(db),
	}
}
