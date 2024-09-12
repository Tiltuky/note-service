package service

import (
	"note-service/internal/models"
	"note-service/internal/repository"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type Notes interface {
	Create(userId int, note models.Note) (int, error)
	GetAll(userId int) ([]models.Note, error)
	GetById(userId, noteId int) (models.Note, error)
	Delete(userId, noteId int) error
	Update(userId, noteId int, input models.UpdateNoteInput) error
}

type Service struct {
	Authorization
	Notes
}

func NewService(rep *repository.Repository, spellerClient *SpellerClient) *Service {
	return &Service{
		Authorization: NewAuthService(rep.Authorization),
		Notes:         NewNoteService(rep.Notes, spellerClient),
	}
}
