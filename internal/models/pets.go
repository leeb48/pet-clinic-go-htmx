package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"pet-clinic.bonglee.com/internal/models/customErrors"
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
		var mySqlError *mysql.MySQLError
		if errors.As(err, &mySqlError) {
			fmt.Println(mySqlError.Number)
			if mySqlError.Number == customErrors.MY_SQL_CONSTRAINT_CODE {
				return customErrors.CheckConstraintError
			}
		}
		fmt.Println(err)
		return err
	}

	return nil
}
