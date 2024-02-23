package models

import (
	"database/sql"
	"time"

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
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Birthdate time.Time `json:"birthdate"`
	PetType   string    `json:"petType"`
	OwnerId   int
}

type PetModelInterface interface {
	Insert(name string, birthdate time.Time, petTypeId, ownerId int) error
	GetPetsByOwnerId(ownerId int) ([]PetDetail, error)
	Remove(id int) error
	Update(id int, name string, birthdate time.Time, petTypeId int) error
	GetPetsByNameAndDob(name string, birthdate time.Time) ([]PetDetail, error)
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
		return customErrors.HandleMySqlError(err)
	}

	return nil
}

func (model *PetModel) GetById(id int) (Pet, error) {
	pet := Pet{}

	stmt := `
		SELECT id, name, birthdate, petTypeId
		FROM pets
		WHERE id = ?
	`

	err := model.DB.QueryRow(stmt, id).Scan(&pet.Id, &pet.Name, &pet.Birthdate, &pet.PetTypeId)
	if err != nil {
		return pet, err
	}

	return pet, nil
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

func (model *PetModel) Remove(id int) error {
	stmt := `
		DELETE FROM pets WHERE id = ?;
	`
	_, err := model.DB.Exec(stmt, id)
	if err != nil {
		return err
	}

	return nil
}

func (model *PetModel) Update(id int, name string, birthdate time.Time, petTypeId int) error {

	stmt := `
		UPDATE
			pets
		SET
			name = COALESCE(NULLIF(?, ''), name),
			birthdate = COALESCE(NULLIF(?, ''), birthdate),
			petTypeId = COALESCE(NULLIF(?, 0), petTypeId),
			modifiedDate = UTC_TIMESTAMP()
		WHERE
			id = ?;
	`

	_, err := model.DB.Exec(stmt, name, birthdate, petTypeId, id)
	if err != nil {
		return err
	}

	return nil
}

func (model *PetModel) GetPetsByNameAndDob(name string, birthdate time.Time) ([]PetDetail, error) {
	pets := []PetDetail{}

	stmt := `
		SELECT pt.name, pt.birthdate, pty.name, pt.ownerId
		FROM pets pt
		INNER JOIN petTypes pty on pty.id = petTypeId
		WHERE pt.name = ? and pt.birthdate = ?
	`

	rows, err := model.DB.Query(stmt, name, birthdate)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var pet PetDetail
		if err := rows.Scan(&pet.Name, &pet.Birthdate, &pet.PetType, &pet.OwnerId); err != nil {
			return nil, err
		}
		pets = append(pets, pet)
	}

	return pets, nil
}
