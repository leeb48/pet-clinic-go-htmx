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
