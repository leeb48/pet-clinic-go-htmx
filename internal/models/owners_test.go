package models

import (
	"errors"
	"testing"
	"time"

	"pet-clinic.bonglee.com/internal/models/customErrors"
)

func TestOwnerModelInsert(t *testing.T) {

	tests := []struct {
		name  string
		owner Owner
		want  error
	}{
		{
			name: "Valid new owner request",
			owner: Owner{
				FirstName: "Bong",
				LastName:  "Lee",
				Email:     "test@test.com",
				Phone:     "2223334444",
				Birthdate: time.Now(),
				Address:   "1234 S Street",
				City:      "Las Vegas",
				State:     "NV",
			},
			want: nil,
		},
		{
			name: "Missing FirstName request",
			owner: Owner{
				FirstName: "",
				LastName:  "Lee",
				Email:     "test@test.com",
				Phone:     "2223334444",
				Birthdate: time.Now(),
				Address:   "1234 S Street",
				City:      "Las Vegas",
				State:     "NV",
			},
			want: customErrors.ErrConstraintFail,
		},
		{
			name: "Missing Phone request",
			owner: Owner{
				FirstName: "",
				LastName:  "Lee",
				Email:     "test@test.com",
				Phone:     "",
				Birthdate: time.Now(),
				Address:   "1234 S Street",
				City:      "Las Vegas",
				State:     "NV",
			},
			want: customErrors.ErrConstraintFail,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db := newTestDB(t)
			model := OwnerModel{db}

			owner := test.owner

			_, err := model.Insert(
				owner.FirstName, owner.LastName, owner.Address, owner.State, owner.City, owner.Phone, owner.Email, owner.Birthdate.Format("2006-01-02"),
			)

			if !errors.Is(err, test.want) {
				t.Errorf("got: %v; want: %v", err, test.want)
			}
		})
	}
}
