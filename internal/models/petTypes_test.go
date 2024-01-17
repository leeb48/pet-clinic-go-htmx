package models

import (
	"errors"
	"testing"

	"pet-clinic.bonglee.com/internal/assert"
	"pet-clinic.bonglee.com/internal/models/customErrors"
)

func TestPetTypeModelInsert(t *testing.T) {

	tests := []struct {
		name    string
		petType string
		want    error
	}{
		{
			name:    "Valid pet type",
			petType: "CAT",
			want:    nil,
		},
		{
			name:    "Empty pet type",
			petType: "",
			want:    customErrors.CheckConstraintError,
		},
		{
			name:    "Duplicate pet type",
			petType: "DOG",
			want:    customErrors.ErrDuplicatePetType,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db := newTestDB(t)
			model := PetTypeModel{db}

			err := model.Insert(test.petType)

			if !errors.Is(err, test.want) {
				t.Errorf("got: %v; want: %v", err, test.want)
			}
		})
	}
}

func TestPetTypeGetAll(t *testing.T) {
	t.Run("Get all pet types", func(t *testing.T) {
		db := newTestDB(t)
		model := PetTypeModel{db}

		petTypes, err := model.GetAll()

		assert.NilError(t, err)

		assert.Equal(t, len(petTypes), 1)
	})
}
