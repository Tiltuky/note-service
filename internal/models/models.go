package models

import "errors"

type User struct {
	Id       int    `json:"-" db:"id"`
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Note struct {
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title" binding:"required"`
	Description string `json:"description" db:"description"`
}

type UsersNote struct {
	Id     int
	UserId int
	NoteId int
}

type UpdateNoteInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
}

func (i UpdateNoteInput) Validate() error {
	if i.Title == nil && i.Description == nil {
		return errors.New("update structure has no values")
	}

	return nil
}
