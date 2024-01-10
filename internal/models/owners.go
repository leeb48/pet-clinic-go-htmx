package models

import (
	"database/sql"
	"time"
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
		return 0, err
	}

	ownerId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(ownerId), nil
}
