package mocks

import (
	"time"

	"pet-clinic.bonglee.com/internal/models"
)

type PetModel struct{}

func (model *PetModel) Insert(name string, birthdate time.Time, petTypeId, ownerId int) error {

	return nil
}

func (model *PetModel) GetPetsByOwnerId(ownerId int) ([]models.PetDetail, error) {
	return nil, nil
}

func (model *PetModel) Remove(id int) error {
	return nil
}

func (model *PetModel) Update(id int, name string, birthdate time.Time, petTypeId int) error {

	return nil
}
