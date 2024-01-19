package mocks

import (
	"errors"

	"pet-clinic.bonglee.com/internal/models"
)

type OwnerModel struct{}

func (model *OwnerModel) Insert(firstName, lastName, addr, state, city, phone, email, birthdate string) (int, error) {

	if firstName == "ownerModelError" {
		return 0, errors.New("DB error")
	}

	return 0, nil
}

func (model *OwnerModel) GetOwnerPageLen() (int, error) {
	return 0, nil
}

func (model *OwnerModel) GetOwnersPage(page, pageSize int) ([]models.Owner, error) {

	owners := []models.Owner{}

	return owners, nil
}
