package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"testing"
	"time"

	"pet-clinic.bonglee.com/internal/app"
	"pet-clinic.bonglee.com/internal/assert"
	"pet-clinic.bonglee.com/internal/models"
	"pet-clinic.bonglee.com/internal/models/mocks"
)

func TestOwnerCreatePost(t *testing.T) {
	testApp := app.NewTestApp(t)
	testServer := app.NewTestServer(t, Routes(testApp))

	tests := []struct {
		name            string
		createOwnerForm CreateOwnerForm
		urlPath         string
		ownerCount      int
		wantCode        int
		formTag         string
	}{
		{
			name:    "Valid new owner request",
			urlPath: "/owner/create",
			createOwnerForm: CreateOwnerForm{
				Owner: models.OwnerCreateDto{
					FirstName: "Bong",
					LastName:  "Lee",
					Email:     "test@test.com",
					Phone:     "2223334444",
					Birthdate: "03-03-1991",
					Address:   "1234 S Street",
					City:      "Las Vegas",
					State:     "NV",
				},
			},

			ownerCount: len(mocks.MockOwners) + 1,
			wantCode:   http.StatusOK,
			formTag:    "",
		},
		{
			name:    "Missing FirstName",
			urlPath: "/owner/create",
			createOwnerForm: CreateOwnerForm{
				Owner: models.OwnerCreateDto{
					FirstName: "",
					LastName:  "Lee",
					Email:     "test@test.com",
					Phone:     "2223334444",
					Birthdate: "03-03-1991",
					Address:   "1234 S Street",
					City:      "Las Vegas",
					State:     "NV",
				},
			},
			ownerCount: len(mocks.MockOwners),
			wantCode:   http.StatusUnprocessableEntity,
			formTag:    `<form hx-post='/owner/create' novalidate hx-target="body" hx-target-4*="body">`,
		},
		{
			name:    "Owner model error",
			urlPath: "/owner/create",

			createOwnerForm: CreateOwnerForm{
				Owner: models.OwnerCreateDto{
					FirstName: "ownerModelError",
					LastName:  "Lee",
					Email:     "test@test.com",
					Phone:     "2223334444",
					Birthdate: "03-03-1991",
					Address:   "1234 S Street",
					City:      "Las Vegas",
					State:     "NV",
				},
			},
			ownerCount: len(mocks.MockOwners),
			wantCode:   http.StatusUnprocessableEntity,
			formTag:    `<form hx-post='/owner/create' novalidate hx-target="body" hx-target-4*="body">`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			mocks.MockOwners = mocks.ResetMockOwners()

			jsonData, err := json.Marshal(test.createOwnerForm)
			if err != nil {
				t.Fatal(err)
			}

			statusCode, _, body := testServer.PostReq(t, test.urlPath, jsonData)

			assert.Equal(t, statusCode, test.wantCode)
			assert.StringContains(t, string(body), test.formTag)
			assert.Equal(t, len(mocks.MockOwners), test.ownerCount)
		})
	}
}

func TestOwnerList(t *testing.T) {
	testApp := app.NewTestApp(t)
	testServer := app.NewTestServer(t, Routes(testApp))

	tests := []struct {
		name           string
		urlPath        string
		wantCode       int
		wantPageLen    string
		wantOwnerCount int
	}{
		{
			name:           "Default owner list",
			urlPath:        "/owner",
			wantCode:       http.StatusOK,
			wantPageLen:    "const pageLen = '0'",
			wantOwnerCount: len(mocks.MockOwners),
		},
		{
			name:           "Default owner list",
			urlPath:        "/owner?pageSize=1",
			wantCode:       http.StatusOK,
			wantPageLen:    fmt.Sprintf("const pageLen = '%v'", len(mocks.MockOwners)),
			wantOwnerCount: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			statusCode, _, body := testServer.GetReq(t, test.urlPath)

			assert.Equal(t, statusCode, test.wantCode)

			assert.StringContains(t, string(body), test.wantPageLen)

			matchOwnerRow := regexp.MustCompile("<td>mangs@test.com</td>")
			matches := matchOwnerRow.FindAllStringIndex(string(body), -1)
			assert.Equal(t, len(matches), test.wantOwnerCount)
		})
	}
}

func TestOwnerDetail(t *testing.T) {
	testApp := app.NewTestApp(t)
	testServer := app.NewTestServer(t, Routes(testApp))

	tests := []struct {
		name      string
		urlPath   string
		wantCode  int
		wantOwner string
	}{
		{
			name:      "Valid Owner",
			urlPath:   "/owner/detail/1",
			wantCode:  http.StatusOK,
			wantOwner: `<h5 class="my-3">Lee, Mango</h5>`,
		},
		{
			name:      "Owner does not exist",
			urlPath:   "/owner/detail/99",
			wantCode:  http.StatusOK,
			wantOwner: `<h1>Owner does not exist</h1>`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			statusCode, _, body := testServer.GetReq(t, test.urlPath)

			assert.Equal(t, statusCode, test.wantCode)
			assert.StringContains(t, body, test.wantOwner)
		})
	}
}

func TestOwnerEditPut(t *testing.T) {
	testApp := app.NewTestApp(t)
	testServer := app.NewTestServer(t, Routes(testApp))

	tests := []struct {
		name           string
		urlPath        string
		wantCode       int
		editOwnerForm  EditOwnerForm
		petCount       int
		petUpdateCount int
	}{
		{
			name:     "Valid owner edit (with additional pet)",
			urlPath:  "/owner/edit/1",
			wantCode: http.StatusOK,
			editOwnerForm: EditOwnerForm{
				Id: 1,
				Owner: models.OwnerCreateDto{
					FirstName: "Bong",
					LastName:  "Lee",
					Email:     "test@test.com",
					Phone:     "2223334444",
					Birthdate: "2018-05-05",
					Address:   "1234 S Street",
					City:      "Las Vegas",
					State:     "NV",
				},
				Pets: []models.PetDetail{
					{
						Id:        0,
						Name:      "Mangoose",
						Birthdate: time.Now(),
						PetType:   "DOG",
					},
				},
			},
			petCount:       len(mocks.MockPets) + 1,
			petUpdateCount: 0,
		},

		{
			name:     "Valid owner edit",
			urlPath:  "/owner/edit/1",
			wantCode: http.StatusOK,
			editOwnerForm: EditOwnerForm{
				Id: 1,
				Owner: models.OwnerCreateDto{
					FirstName: "Bong",
					LastName:  "Lee",
					Email:     "test@test.com",
					Phone:     "2223334444",
					Birthdate: "2018-05-05",
					Address:   "1234 S Street",
					City:      "Las Vegas",
					State:     "NV",
				},
				Pets: []models.PetDetail{
					{
						Id:        1,
						Name:      "Mangoose",
						Birthdate: time.Now(),
						PetType:   "DOG",
					},
				},
			},
			petCount:       len(mocks.MockPets),
			petUpdateCount: 1,
		},

		{
			name:     "Valid owner edit remove pets",
			urlPath:  "/owner/edit/1",
			wantCode: http.StatusOK,
			editOwnerForm: EditOwnerForm{
				Id: 1,
				Owner: models.OwnerCreateDto{
					FirstName: "Bong",
					LastName:  "Lee",
					Email:     "test@test.com",
					Phone:     "2223334444",
					Birthdate: "2018-05-05",
					Address:   "1234 S Street",
					City:      "Las Vegas",
					State:     "NV",
				},
				Pets:         []models.PetDetail{},
				DeletePetIds: []int{1},
			},
			petCount: len(mocks.MockPets) - 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			mocks.MockPets = mocks.ResetMockPets()
			mocks.PetUpdateCount = 0

			jsonData, err := json.Marshal(test.editOwnerForm)
			if err != nil {
				t.Fatal(err)
			}

			_, header, _ := testServer.PutReq(t, test.urlPath, jsonData)

			assert.Equal(t, header.Get("HX-Redirect"), fmt.Sprintf("/owner/detail/%v", test.editOwnerForm.Id))
			assert.Equal(t, len(mocks.MockPets), test.petCount)
			assert.Equal(t, mocks.PetUpdateCount, test.petUpdateCount)

		})
	}
}

func TestOwnerRemove(t *testing.T) {
	testApp := app.NewTestApp(t)
	testServer := app.NewTestServer(t, Routes(testApp))

	tests := []struct {
		name       string
		urlPath    string
		wantCode   int
		ownerCount int
	}{
		{
			name:       "Valid owner remove",
			urlPath:    "/owner/1",
			wantCode:   http.StatusOK,
			ownerCount: len(mocks.MockOwners) - 1,
		},

		{
			name:       "Non-existent owner remove",
			urlPath:    "/owner/3",
			wantCode:   http.StatusOK,
			ownerCount: len(mocks.MockOwners),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mocks.MockOwners = mocks.ResetMockOwners()

			statusCode, _, _ := testServer.DeleteReq(t, test.urlPath)

			assert.Equal(t, statusCode, test.wantCode)
			assert.Equal(t, len(mocks.MockOwners), test.ownerCount)
		})
	}
}
