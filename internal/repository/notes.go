package repository

import (
	"fmt"
	"note-service/internal/models"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type NotesPostgres struct {
	db *sqlx.DB
}

func NewNotesPostgres(db *sqlx.DB) *NotesPostgres {
	return &NotesPostgres{db: db}
}

func (r *NotesPostgres) Create(userId int, note models.Note) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createNoteQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", notesTable)
	row := tx.QueryRow(createNoteQuery, note.Title, note.Description)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	createUsersNoteQuery := fmt.Sprintf("INSERT INTO %s (user_id, note_id) VALUES ($1, $2)", usersNotesTable)
	_, err = tx.Exec(createUsersNoteQuery, userId, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (r *NotesPostgres) GetAll(userId int) ([]models.Note, error) {
	var notes []models.Note

	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul on tl.id = ul.note_id WHERE ul.user_id = $1",
		notesTable, usersNotesTable)
	err := r.db.Select(&notes, query, userId)

	return notes, err
}

func (r *NotesPostgres) GetById(userId, noteId int) (models.Note, error) {
	var note models.Note

	query := fmt.Sprintf(`SELECT tl.id, tl.title, tl.description FROM %s tl
								INNER JOIN %s ul on tl.id = ul.note_id WHERE ul.user_id = $1 AND ul.note_id = $2`,
		notesTable, usersNotesTable)
	err := r.db.Get(&note, query, userId, noteId)

	return note, err
}

func (r *NotesPostgres) Delete(userId, noteId int) error {
	query := fmt.Sprintf("DELETE FROM %s tl USING %s ul WHERE tl.id = ul.note_id AND ul.user_id=$1 AND ul.note_id=$2",
		notesTable, usersNotesTable)
	_, err := r.db.Exec(query, userId, noteId)

	return err
}

func (r *NotesPostgres) Update(userId, noteId int, input models.UpdateNoteInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	// title=$1
	// description=$1
	// title=$1, description=$2
	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s tl SET %s FROM %s ul WHERE tl.id = ul.note_id AND ul.note_id=$%d AND ul.user_id=$%d",
		notesTable, setQuery, usersNotesTable, argId, argId+1)
	args = append(args, noteId, userId)

	logrus.Debugf("updateQuery: %s", query)
	logrus.Debugf("args: %s", args)

	_, err := r.db.Exec(query, args...)
	return err
}
