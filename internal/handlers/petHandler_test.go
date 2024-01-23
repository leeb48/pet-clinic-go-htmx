package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"pet-clinic.bonglee.com/internal/app"
	"pet-clinic.bonglee.com/internal/assert"
)

func TestNewPetTypePost(t *testing.T) {
	testApp := app.NewTestApp(t)
	testServer := app.NewTestServer(t, Routes(testApp))

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
			jsonData, err := json.Marshal(data)
			if err != nil {
				t.Fatal(err)
			}

			statusCode, _, body := testServer.PostReq(t, test.urlPath, jsonData)

			matches := app.GetFormTextDangerHtml.FindStringSubmatch(string(body))

			if len(matches) > 1 {
				assert.Equal(t, matches[1], test.errMsg)
			}

			assert.Equal(t, statusCode, test.wantCode)
		})
	}
}
