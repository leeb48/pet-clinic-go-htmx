package models

import (
	"database/sql"

	"pet-clinic.bonglee.com/internal/models/customErrors"
)

type Vet struct {
	Id        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type VetModelInterface interface {
	Insert(firstName, lastName string) (int, error)
	GetVetsPageLen(pageSize int) (int, error)
	GetVets(page, pageSize int) ([]Vet, error)
	GetById(id int) (Vet, error)
	Update(id int, firstName, lastName string) error
	Remove(id int) error
	GetVetsByLastName(lastName string, page, pageSize int) ([]Vet, error)
}

type VetModel struct {
	DB *sql.DB
}

func (model *VetModel) Insert(firstName, lastName string) (int, error) {
	stmt := `
		INSERT INTO vets (firstName, lastName)
		VALUES (?, ?)
	`
	result, err := model.DB.Exec(stmt, firstName, lastName)
	if err != nil {
		return 0, customErrors.HandleMySqlError(err)
	}

	vetId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(vetId), nil
}

func (model *VetModel) GetVetsPageLen(pageSize int) (int, error) {

	stmt := `
		SELECT COUNT(*) FROM vets
	`
	var rowCount int

	err := model.DB.QueryRow(stmt).Scan(&rowCount)

	if err != nil {
		return 0, err
	}

	return rowCount / pageSize, nil
}

func (model *VetModel) GetVets(page, pageSize int) ([]Vet, error) {

	vets := []Vet{}

	offset := (page - 1) * pageSize

	stmt := `
		SELECT id, firstName, lastName
		FROM vets
		LIMIT ?
		OFFSET ?
	`

	rows, err := model.DB.Query(stmt, pageSize, offset)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var vet Vet
		if err := rows.Scan(&vet.Id, &vet.FirstName, &vet.LastName); err != nil {
			return nil, err
		}

		vets = append(vets, vet)
	}

	return vets, nil
}

func (model *VetModel) GetById(id int) (Vet, error) {
	vet := Vet{}

	stmt := `
		SELECT id, firstName, lastName
		FROM vets
		WHERE id = ?
	`

	err := model.DB.QueryRow(stmt, id).Scan(&vet.Id, &vet.FirstName, &vet.LastName)
	if err != nil {
		return vet, err
	}

	return vet, nil
}

func (model *VetModel) Update(id int, firstName, lastName string) error {

	stmt := `
		UPDATE
			vets
		SET
			firstName = COALESCE(NULLIF(?, ''), firstName),
			lastName = COALESCE(NULLIF(?, ''), lastName)
		WHERE
			id = ?
	`

	_, err := model.DB.Exec(stmt, firstName, lastName, id)
	if err != nil {
		return err
	}

	return nil
}

func (model *VetModel) Remove(id int) error {

	stmt := `
		DELETE FROM vets WHERE id = ?
	`

	_, err := model.DB.Exec(stmt, id)
	if err != nil {
		return err
	}

	return nil
}

func (model *VetModel) GetVetsByLastName(lastName string, page, pageSize int) ([]Vet, error) {
	vets := []Vet{}

	stmt := `
		SELECT id, firstName, lastName
		FROM vets
		WHERE lastName = ?
		LIMIT ?
		OFFSET ?
	`

	rows, err := model.DB.Query(stmt, lastName, pageSize, page)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var vet Vet
		if err := rows.Scan(&vet.Id, &vet.FirstName, &vet.LastName); err != nil {
			return nil, err
		}

		vets = append(vets, vet)
	}

	return vets, nil
}
