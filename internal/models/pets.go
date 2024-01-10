package models

import (
	"database/sql"
	"time"
)

type Pet struct {
	Id        int
	Name      string
	Birthdate time.Time
	PetTypeId int
	OwnerId   int
}

type PetModelInterface interface {
	Insert(name string, birthdate time.Time, petTypeId, ownerId int) error
}

type PetModel struct {
	DB *sql.DB
}

func (model *PetModel) Insert(name string, birthdate time.Time, petTypeId, ownerId int) error {

	stmt := `
		INSERT INTO pets (name, birthdate, petTypeId, ownerId, created)
		VALUES (?, ?, ?, ?, UTC_TIMESTAMP())
	`

	_, err := model.DB.Exec(stmt, name, birthdate, petTypeId, ownerId)
	if err != nil {
		return err
	}

	return nil
}
