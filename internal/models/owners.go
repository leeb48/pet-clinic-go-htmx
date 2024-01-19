package models

import (
	"database/sql"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
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
		var mySqlError *mysql.MySQLError
		if errors.As(err, &mySqlError) {
			if mySqlError.Number == customErrors.MY_SQL_CONSTRAINT_CODE {
				return 0, customErrors.CheckConstraintError
			}
		}
		return 0, err
	}

	ownerId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(ownerId), nil
}

func (model *OwnerModel) GetOwnerPageLen(pageSize int) (int, error) {

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

func (model *OwnerModel) GetOwnersPage(page, pageSize int) ([]Owner, error) {

	owners := []Owner{}

	return owners, nil
}
