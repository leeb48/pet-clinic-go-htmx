package handlers

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"pet-clinic.bonglee.com/internal/app"
	"pet-clinic.bonglee.com/internal/assert"
	"pet-clinic.bonglee.com/internal/models"
)

func TestOwnerCreatePost(t *testing.T) {
	testApp := app.NewTestApp(t)
	testServer := app.NewTestServer(t, Routes(testApp))

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
			jsonData, err := json.Marshal(test.owner)
			if err != nil {
				t.Fatal(err)
			}

			statusCode, _, body := testServer.PostReq(t, test.urlPath, jsonData)

			assert.Equal(t, statusCode, test.wantCode)
			assert.StringContains(t, string(body), test.formTag)
		})
	}
}
