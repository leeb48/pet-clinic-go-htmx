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
			want: customErrors.CheckConstraintError,
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
