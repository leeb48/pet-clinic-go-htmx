package mocks

import "time"

type PetModel struct{}

func (model *PetModel) Insert(name string, birthdate time.Time, petTypeId, ownerId int) error {

	return nil
}
