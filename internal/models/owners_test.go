package models

import (
	"errors"
	"testing"
	"time"

	"pet-clinic.bonglee.com/internal/assert"
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

func TestGetOwnersPageLen(t *testing.T) {
	tests := []struct {
		name        string
		pageSize    int
		wantErr     error
		wantPageLen int
	}{
		{
			name:        "Page size greate than owner count",
			pageSize:    10,
			wantErr:     nil,
			wantPageLen: 0,
		},
		{
			name:        "Page size of 1",
			pageSize:    1,
			wantErr:     nil,
			wantPageLen: 3,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db := newTestDB(t)
			model := OwnerModel{db}

			pageLen, err := model.GetOwnersPageLen(test.pageSize)

			assert.Equal(t, pageLen, test.wantPageLen)
			assert.Equal(t, err, test.wantErr)
		})
	}
}

func TestGetOwners(t *testing.T) {
	tests := []struct {
		name       string
		page       int
		pageSize   int
		wantErr    error
		ownerCount int
	}{
		{
			name:       "Page size of 10",
			page:       1,
			pageSize:   10,
			wantErr:    nil,
			ownerCount: 3,
		},

		{
			name:       "Page size of 1",
			page:       1,
			pageSize:   1,
			wantErr:    nil,
			ownerCount: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db := newTestDB(t)
			model := OwnerModel{db}

			owners, err := model.GetOwners(test.page, test.pageSize)

			assert.NilError(t, err)
			assert.Equal(t, len(owners), test.ownerCount)
		})
	}
}

func TestGetOnwerById(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		wantErr error
		wantId  int
	}{
		{
			name:    "Owner exists",
			id:      1,
			wantErr: nil,
			wantId:  1,
		},

		{
			name:    "Owner does not exists",
			id:      99,
			wantErr: nil,
			wantId:  0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db := newTestDB(t)
			model := OwnerModel{db}

			owner, err := model.GetOwnerById(test.id)

			assert.Equal(t, owner.Id, test.wantId)
			assert.NilError(t, err)
		})
	}
}

func TestUpdateOwner(t *testing.T) {
	tests := []struct {
		name      string
		Id        int
		FirstName string
		LastName  string
		Email     string
		Phone     string
		Birthdate string
		Address   string
		City      string
		State     string
		wantErr   error
	}{
		{
			name:      "Update All Fields",
			Id:        1,
			FirstName: "Mango",
			LastName:  "Wu",
			Email:     "updated@test.com",
			Phone:     "9998881111",
			Birthdate: "2200-12-12",
			Address:   "113 updated st",
			City:      "new city",
			State:     "new state",
			wantErr:   nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db := newTestDB(t)
			model := OwnerModel{db}

			err := model.UpdateOwner(test.Id, test.FirstName, test.LastName, test.Address,
				test.State, test.City, test.Phone, test.Email, test.Birthdate)
			assert.NilError(t, err)

			owner, err := model.GetOwnerById(test.Id)
			assert.NilError(t, err)

			assert.Equal(t, owner.Id, test.Id)
			assert.Equal(t, owner.FirstName, test.FirstName)
			assert.Equal(t, owner.LastName, test.LastName)
			assert.Equal(t, owner.Address, test.Address)
			assert.Equal(t, owner.State, test.State)
			assert.Equal(t, owner.City, test.City)
			assert.Equal(t, owner.Phone, test.Phone)
			assert.Equal(t, owner.Email, test.Email)

			birthdateParsed, err := time.Parse(time.DateOnly, test.Birthdate)
			assert.NilError(t, err)

			assert.Equal(t, owner.Birthdate, birthdateParsed)
		})
	}
}

func TestRemoveOnwer(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		wantErr error
	}{
		{
			name:    "Remove owner",
			id:      1,
			wantErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db := newTestDB(t)
			model := OwnerModel{db}

			err := model.RemoveOwner(test.id)
			assert.NilError(t, err)

			owner, err := model.GetOwnerById(test.id)
			assert.NilError(t, err)

			assert.Equal(t, owner.Id, 0)
		})
	}
}
