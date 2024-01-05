package models

import "testing"

func TestPetTypeModelInsert(t *testing.T) {
	db := newTestDB(t)
	model := PetTypeModel{db}

	model.Insert("Test Type")
}
