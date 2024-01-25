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

type PetDetail struct {
	Id        int
	Name      string
	Birthdate time.Time
	PetType   string
}

type PetModelInterface interface {
	Insert(name string, birthdate time.Time, petTypeId, ownerId int) error
	GetPetsByOwnerId(ownerId int) ([]PetDetail, error)
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

func (model *PetModel) GetPetsByOwnerId(ownerId int) ([]PetDetail, error) {

	petDetails := []PetDetail{}

	stmt := `
		SELECT pt.id, pt.name, pt.birthdate, py.name
		FROM pets pt
		INNER JOIN petTypes py on py.id = pt.petTypeId
		WHERE pt.ownerId = ?
	`

	rows, err := model.DB.Query(stmt, ownerId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var petDetail PetDetail
		if err := rows.Scan(&petDetail.Id, &petDetail.Name, &petDetail.Birthdate, &petDetail.PetType); err != nil {
			return nil, err
		}
		petDetails = append(petDetails, petDetail)
	}

	return petDetails, nil
}
