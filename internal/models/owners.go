package models

import (
	"database/sql"
	"time"

	"pet-clinic.bonglee.com/internal/models/customErrors"
)

type Owner struct {
	Id        int
	FirstName string
	LastName  string
	Address   string
	State     string
	City      string
	Phone     string
	Email     string
	Birthdate time.Time
	Created   time.Time
}

type OwnerModelInterface interface {
	Insert(firstName, lastName, addr, state, city, phone, email, birthdate string) (int, error)
	GetOwnersPageLen(pageSize int) (int, error)
	GetOwners(page, pageSize int) ([]Owner, error)
	GetOwnerById(id int) (Owner, error)
	UpdateOwner(id int, firstName, lastName, addr, state, city, phone, email, birthdate string) error
	Remove(id int) error
}

type OwnerModel struct {
	DB *sql.DB
}

func (model *OwnerModel) Insert(firstName, lastName, addr, state, city, phone, email, birthdate string) (int, error) {
	stmt := `
		INSERT INTO owners (firstName, lastName, address, state, city, phone, email, birthdate, created)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, UTC_TIMESTAMP())
		`

	result, err := model.DB.Exec(stmt, firstName, lastName, addr, state, city, phone, email, birthdate)
	if err != nil {
		return 0, customErrors.HandleMySqlError(err)
	}

	ownerId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(ownerId), nil
}

func (model *OwnerModel) GetOwnersPageLen(pageSize int) (int, error) {

	stmt := `
		SELECT COUNT(*) FROM owners
	`
	var rowCount int

	err := model.DB.QueryRow(stmt).Scan(&rowCount)

	if err != nil {
		return 0, err
	}

	return rowCount / pageSize, nil
}

func (model *OwnerModel) GetOwners(page, pageSize int) ([]Owner, error) {

	owners := []Owner{}

	offset := (page - 1) * pageSize

	stmt := `
		SELECT id, firstName, lastName, address, state, city, phone, email, birthdate
		FROM owners
		LIMIT ?
		OFFSET ?
	`

	rows, err := model.DB.Query(stmt, pageSize, offset)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var owner Owner
		if err := rows.Scan(&owner.Id, &owner.FirstName, &owner.LastName, &owner.Address,
			&owner.State, &owner.City, &owner.Phone, &owner.Email, &owner.Birthdate); err != nil {
			return nil, err
		}

		owners = append(owners, owner)
	}

	return owners, nil
}

func (model *OwnerModel) GetOwnerById(id int) (Owner, error) {

	owner := Owner{}

	getOwnerStmt := `
		SELECT id, firstName, lastName, address, state, city, phone, email, birthdate
		FROM owners
		WHERE id = ?
	`

	rows, err := model.DB.Query(getOwnerStmt, id)
	if err != nil {
		return owner, err
	}

	for rows.Next() {
		if err := rows.Scan(&owner.Id, &owner.FirstName, &owner.LastName, &owner.Address, &owner.State, &owner.City, &owner.Phone, &owner.Email, &owner.Birthdate); err != nil {
			return owner, err
		}
	}

	return owner, nil
}

func (model *OwnerModel) UpdateOwner(id int, firstName, lastName, addr, state, city, phone, email, birthdate string) error {

	stmt := `
		UPDATE
			owners
		SET
			firstName = COALESCE(NULLIF(?, ''), firstName),
			lastName = COALESCE(NULLIF(?, ''), lastName),
			email = COALESCE(NULLIF(?, ''), email),
			phone = COALESCE(NULLIF(?, ''), phone),
			birthdate = COALESCE(NULLIF(?, ''), birthdate),
			address = COALESCE(NULLIF(?, ''), address),
			city = COALESCE(NULLIF(?, ''), city),
			state = COALESCE(NULLIF(?, ''), state),
			modifiedDate = UTC_TIMESTAMP()
		WHERE
			id = ?;
	`

	_, err := model.DB.Exec(stmt, firstName, lastName, email, phone, birthdate, addr, city, state, id)
	if err != nil {
		return err
	}

	return nil
}

func (model *OwnerModel) Remove(id int) error {

	stmt := `
		DELETE FROM owners WHERE id = ?
	`

	_, err := model.DB.Exec(stmt, id)
	if err != nil {
		return err
	}

	return nil
}
