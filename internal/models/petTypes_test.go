package models

import (
	"errors"
	"testing"

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
