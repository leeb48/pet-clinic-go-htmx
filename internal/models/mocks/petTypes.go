package mocks

import (
	"pet-clinic.bonglee.com/internal/models"
	"pet-clinic.bonglee.com/internal/models/customErrors"
)

var MockPetType = models.PetType{
	Id:   1,
	Name: "DOG",
}

type PetTypeModel struct{}

func (model *PetTypeModel) Insert(petType string) error {
	switch petType {
	case MockPetType.Name:
		return customErrors.ErrDuplicatePetType

	default:
		return nil
	}
}

func (model *PetTypeModel) GetAll() ([]string, error) {

	return []string{}, nil
}

func (model *PetTypeModel) GetIdFromPetType(string) (int, error) {
	return 0, nil
}
