package mocks

import (
	"time"

	"pet-clinic.bonglee.com/internal/models"
)

type PetModel struct{}

var MockPets []models.Pet = ResetMockPets()
var PetUpdateCount int

func (model *PetModel) Insert(name string, birthdate time.Time, petTypeId, ownerId int) error {

	MockPets = append(MockPets, models.Pet{
		Id: len(MockPets) + 1, Name: name, Birthdate: birthdate, PetTypeId: petTypeId, OwnerId: ownerId,
	})

	return nil
}

func (model *PetModel) GetPetsByOwnerId(ownerId int) ([]models.PetDetail, error) {
	return nil, nil
}

func (model *PetModel) Remove(id int) error {

	return nil
}

func (model *PetModel) Update(id int, name string, birthdate time.Time, petTypeId int) error {
	PetUpdateCount++
	return nil
}

func ResetMockPets() []models.Pet {
	return []models.Pet{
		{
			Id:        1,
			Name:      "Mango",
			Birthdate: time.Now(),
			PetTypeId: 1,
			OwnerId:   1,
		},
		{
			Id:        2,
			Name:      "Acorn",
			Birthdate: time.Now(),
			PetTypeId: 1,
			OwnerId:   1,
		},
	}
}
