package models

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
	"pet-clinic.bonglee.com/internal/models/customErrors"
)

type PetType struct {
	Id   int
	Name string
}

type PetTypeModelInterface interface {
	Insert(petType string) error
	GetAll() ([]string, error)
	GetIdFromPetType(string) (int, error)
}

type PetTypeModel struct {
	DB *sql.DB
}

func (model *PetTypeModel) Insert(petType string) error {
	stmt := `
		INSERT INTO petTypes (name)
		VALUES(?)
	`

	_, err := model.DB.Exec(stmt, petType)
	if err != nil {
		var mySqlError *mysql.MySQLError
		if errors.As(err, &mySqlError) {
			if mySqlError.Number == customErrors.MY_SQL_DUPLICATE_CODE && strings.Contains(mySqlError.Message, customErrors.DUPLICATE_PET_TYPE_KEY) {
				return customErrors.ErrDuplicatePetType
			}
		}
		return err
	}

	return nil
}

func (model *PetTypeModel) GetAll() ([]string, error) {
	stmt := `
		SELECT name
		FROM petTypes
	`

	rows, err := model.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	petTypes := []string{}

	for rows.Next() {
		var name string

		err = rows.Scan(&name)
		if err != nil {
			return nil, err
		}

		petTypes = append(petTypes, name)
	}

	return petTypes, nil
}

func (model *PetTypeModel) GetIdFromPetType(petType string) (int, error) {
	stmt := `
		SELECT id
		FROM petTypes
		WHERE name = ?
	`

	var id int

	err := model.DB.QueryRow(stmt, petType).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
