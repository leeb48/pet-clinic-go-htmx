package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"testing"
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
			rs, err := testServer.Client().Post(testServer.URL+test.urlPath, "application/json", bytes.NewReader(jsonData))
			if err != nil {
				t.Fatal(err)
			}

			defer rs.Body.Close()
			body, err := io.ReadAll(rs.Body)
			if err != nil {
				t.Fatal(err)
			}
			matches := getFormTextDangerHtml.FindStringSubmatch(string(body))

			if len(matches) > 1 && matches[1] != test.errMsg {
				t.Errorf("got: %s; want %s", matches[1], test.errMsg)
			}

			if rs.StatusCode != test.wantCode {
				t.Errorf("got: %v; want: %v", rs.StatusCode, test.wantCode)
			}
		})
	}
}

func TestOwnerCreatePost(t *testing.T) {
	// app := newTestApp(t)
	// testServer := newTestServer(t, app.routes())

	// tests := []struct {
	// 	name  string
	// 	owner models.Owner
	// 	urlPath string
	// 	wantCode int
	// 	alertMsg string
	// 	formTag string
	// } {
	// 	{
	// 		name: "Valid new owner request",
	// 		owner: models.Owner{
	// 			fi
	// 		},
	// 	}
	// }

}
