package service

import (
	"note-service/internal/models"
	"note-service/internal/repository"

	"github.com/sirupsen/logrus"
)

type NoteService struct {
	rep           repository.Notes
	spellerClient *SpellerClient
}

func NewNoteService(rep repository.Notes, spellerClient *SpellerClient) *NoteService {
	return &NoteService{
		rep:           rep,
		spellerClient: spellerClient}
}

func (s *NoteService) Create(userId int, note models.Note) (int, error) {
	correctedText, err := s.spellerClient.CheckText(note.Description)
	if err != nil {
		logrus.Errorf("error during spell checking: %s", err.Error())
		return 0, err
	}

	note.Description = correctedText

	return s.rep.Create(userId, note)
}

func (s *NoteService) GetAll(userId int) ([]models.Note, error) {
	return s.rep.GetAll(userId)
}

func (s *NoteService) GetById(userId, noteId int) (models.Note, error) {
	return s.rep.GetById(userId, noteId)
}

func (s *NoteService) Delete(userId, noteId int) error {
	return s.rep.Delete(userId, noteId)
}

func (s *NoteService) Update(userId, noteId int, input models.UpdateNoteInput) error {
	if err := input.Validate(); err != nil {
		return err
	}

	return s.rep.Update(userId, noteId, input)
}
