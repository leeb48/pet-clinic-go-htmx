package mocks

import "errors"

type OwnerModel struct{}

func (model *OwnerModel) Insert(firstName, lastName, addr, state, city, phone, email, birthdate string) (int, error) {

	if firstName == "ownerModelError" {
		return 0, errors.New("DB error")
	}

	return 0, nil
}
