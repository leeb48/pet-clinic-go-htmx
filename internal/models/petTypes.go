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
