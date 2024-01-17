package main

import (
	"encoding/json"
	"net/http"
	"regexp"
	"testing"
	"time"

	"pet-clinic.bonglee.com/internal/assert"
	"pet-clinic.bonglee.com/internal/models"
)

var getFormTextDangerHtml = regexp.MustCompile(`<div class="form-text text-danger">(.*) `)

func TestNewPetTypePost(t *testing.T) {

	app := newTestApp(t)
	testServer := newTestServer(t, app.routes())

	tests := []struct {
		name     string
		petType  string
		urlPath  string
		wantCode int
		errMsg   string
	}{
		{
			name:     "Valid pet type",
			petType:  "Cat",
			urlPath:  "/pet/add-pet-type",
			wantCode: http.StatusSeeOther,
			errMsg:   "",
		},
		{
			name:     "Duplicate pet type",
			petType:  "Dog",
			urlPath:  "/pet/add-pet-type",
			wantCode: http.StatusUnprocessableEntity,
			errMsg:   "Duplicate pet type",
		},
		{
			name:     "Empty pet type",
			petType:  "",
			urlPath:  "/pet/add-pet-type",
			wantCode: http.StatusUnprocessableEntity,
			errMsg:   "Pet type cannot be blank",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			data := map[string]string{"newPetType": test.petType}
			jsonData, _ := json.Marshal(data)

			statusCode, _, body := testServer.postReq(t, test.urlPath, jsonData)

			matches := getFormTextDangerHtml.FindStringSubmatch(string(body))

			if len(matches) > 1 {
				assert.Equal(t, matches[1], test.errMsg)
			}

			assert.Equal(t, statusCode, test.wantCode)
		})
	}
}

func TestOwnerCreatePost(t *testing.T) {
	app := newTestApp(t)
	testServer := newTestServer(t, app.routes())

	tests := []struct {
		name     string
		owner    models.Owner
		urlPath  string
		wantCode int
		formTag  string
	}{
		{
			name:    "Valid new owner request",
			urlPath: "/owner/create",
			owner: models.Owner{
				FirstName: "Bong",
				LastName:  "Lee",
				Email:     "test@test.com",
				Phone:     "2223334444",
				Birthdate: time.Now(),
				Address:   "1234 S Street",
				City:      "Las Vegas",
				State:     "NV",
			},
			wantCode: http.StatusSeeOther,
			formTag:  "",
		},
		{
			name:    "Missing FirstName",
			urlPath: "/owner/create",
			owner: models.Owner{
				FirstName: "",
				LastName:  "Lee",
				Email:     "test@test.com",
				Phone:     "2223334444",
				Birthdate: time.Now(),
				Address:   "1234 S Street",
				City:      "Las Vegas",
				State:     "NV",
			},
			wantCode: http.StatusUnprocessableEntity,
			formTag:  `<form hx-post='/owner/create' novalidate hx-target="body" hx-target-4*="body">`,
		},
		{
			name:    "Owner model error",
			urlPath: "/owner/create",
			owner: models.Owner{
				FirstName: "ownerModelError",
				LastName:  "Lee",
				Email:     "test@test.com",
				Phone:     "2223334444",
				Birthdate: time.Now(),
				Address:   "1234 S Street",
				City:      "Las Vegas",
				State:     "NV",
			},
			wantCode: http.StatusUnprocessableEntity,
			formTag:  `<form hx-post='/owner/create' novalidate hx-target="body" hx-target-4*="body">`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			jsonData, _ := json.Marshal(test.owner)

			statusCode, _, body := testServer.postReq(t, test.urlPath, jsonData)

			assert.Equal(t, statusCode, test.wantCode)
			assert.StringContains(t, string(body), test.formTag)
		})
	}

}
