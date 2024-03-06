package mocks

import (
	"errors"
	"slices"
	"time"

	"pet-clinic.bonglee.com/internal/models"
)

type OwnerModel struct{}

var MockOwners = ResetMockOwners()

func ResetMockOwners() []models.Owner {
	return []models.Owner{
		{
			Id:        1,
			FirstName: "Mango",
			LastName:  "Lee",
			Address:   "123 Dog St",
			State:     "NV",
			City:      "Las Vegas",
			Phone:     "7024445678",
			Email:     "mangs@test.com",
			Birthdate: time.Now(),
			Created:   time.Now(),
		},
		{
			Id:        2,
			FirstName: "Mango",
			LastName:  "Lee",
			Address:   "123 Dog St",
			State:     "NV",
			City:      "Las Vegas",
			Phone:     "7024445678",
			Email:     "mangs@test.com",
			Birthdate: time.Now(),
			Created:   time.Now(),
		},
	}
}

func (model *OwnerModel) Insert(firstName, lastName, addr, state, city, phone, email, birthdate string) (int, error) {

	if firstName == "ownerModelError" {
		return 0, errors.New("DB error")
	}

	parsedBirthdate, _ := time.Parse(time.DateOnly, birthdate)

	MockOwners = append(MockOwners, models.Owner{
		FirstName: firstName, LastName: lastName, Address: addr, State: state, City: city, Phone: phone, Email: email, Birthdate: parsedBirthdate,
	})

	return 0, nil
}

func (model *OwnerModel) GetOwnersPageLen(pageSize int) (int, error) {

	return len(MockOwners) / pageSize, nil
}

func (model *OwnerModel) GetOwners(page, pageSize int) ([]models.Owner, error) {

	limit := min(pageSize, len(MockOwners))

	return MockOwners[:limit], nil
}

func (model *OwnerModel) GetOwnerById(id int) (models.Owner, error) {

	idx := slices.IndexFunc(MockOwners, func(o models.Owner) bool {
		return o.Id == id
	})

	if idx != -1 {
		return MockOwners[idx], nil
	}

	return models.Owner{}, nil
}

func (model *OwnerModel) UpdateOwner(id int, firstName, lastName, addr, state, city, phone, email, birthdate string) error {
	return nil
}

func (model *OwnerModel) RemoveOwner(id int) error {

	for idx, owner := range MockOwners {
		if owner.Id == id {
			MockOwners = append(MockOwners[:idx], MockOwners[idx+1:]...)
		}
	}

	return nil
}
