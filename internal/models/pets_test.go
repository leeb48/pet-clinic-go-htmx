package models

import (
	"testing"
	"time"

	"pet-clinic.bonglee.com/internal/assert"
	"pet-clinic.bonglee.com/internal/models/customErrors"
)

func TestPetModelInsert(t *testing.T) {
	tests := []struct {
		name string
		pet  Pet
		want error
	}{
		{
			name: "Valid pet insert",
			pet: Pet{
				Name:      "Mango",
				Birthdate: time.Now(),
				PetTypeId: 1,
				OwnerId:   1,
			},
			want: nil,
		},
		{
			name: "Empty pet name",
			pet: Pet{
				Name:      "",
				Birthdate: time.Now(),
				PetTypeId: 1,
				OwnerId:   1,
			},
			want: customErrors.ErrConstraintFail,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db := newTestDB(t)
			model := PetModel{db}

			pet := test.pet
			err := model.Insert(pet.Name, pet.Birthdate, pet.PetTypeId, pet.OwnerId)

			assert.Equal(t, err, test.want)
		})
	}
}

func TestPetModelGetPetsByOwnerId(t *testing.T) {
	tests := []struct {
		name      string
		ownerId   int
		wantCount int
	}{
		{
			name:      "Owner with two pets",
			ownerId:   1,
			wantCount: 2,
		},
		{
			name:      "Owner with no pets",
			ownerId:   2,
			wantCount: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db := newTestDB(t)
			model := PetModel{db}

			pets, err := model.GetPetsByOwnerId(test.ownerId)

			assert.Equal(t, len(pets), test.wantCount)
			assert.NilError(t, err)
		})
	}
}

func TestPetRemove(t *testing.T) {
	tests := []struct {
		name  string
		petId int
		want  error
	}{
		{
			name:  "Valid pet removal",
			petId: 1,
			want:  nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db := newTestDB(t)
			model := PetModel{db}

			err := model.Remove(test.petId)

			assert.NilError(t, err)
		})
	}
}

func TestPetUpdate(t *testing.T) {
	tests := []struct {
		name         string
		petId        int
		newName      string
		newBirthdate string
		newPetTypeId int

		wantName      string
		wantBirthdate string
		wantPetTypeId int
	}{
		{
			name:         "Valid update",
			petId:        1,
			newName:      "Mangos",
			newBirthdate: "2023-01-02",
			newPetTypeId: 2,

			wantName:      "Mangos",
			wantBirthdate: "2023-01-02",
			wantPetTypeId: 2,
		},
		{
			name:         "Valid update with empty values",
			petId:        1,
			newName:      "",
			newBirthdate: "2020-01-01",
			newPetTypeId: 0,

			wantName:      "Mango",
			wantBirthdate: "2020-01-01",
			wantPetTypeId: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db := newTestDB(t)
			model := PetModel{db}

			newBirthdate, _ := time.Parse(time.DateOnly, test.newBirthdate)

			err := model.Update(test.petId, test.newName, newBirthdate, test.newPetTypeId)
			assert.NilError(t, err)

			updatedPet, err := model.GetById(test.petId)
			assert.NilError(t, err)
			assert.Equal(t, updatedPet.Name, test.wantName)
			assert.Equal(t, updatedPet.Birthdate.Format(time.DateOnly), test.wantBirthdate)
			assert.Equal(t, updatedPet.PetTypeId, test.wantPetTypeId)
		})
	}
}
