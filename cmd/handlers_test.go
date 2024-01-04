package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"pet-clinic.bonglee.com/internal/models/mocks"
)

// todo: see if there is a way to extract the error message from the template
var errorHtmlRX = regexp.MustCompile(`<div class="form-text text-danger">(.*)`)

func TestNewPetTypePost(t *testing.T) {

	templateCache, err := createTemplateCache()
	if err != nil {
		t.Fatal(err)
	}

	app := &application{
		logger:        slog.New(slog.NewTextHandler(io.Discard, nil)),
		templateCache: templateCache,
		petTypes:      &mocks.PetTypeModel{},
	}

	testServer := httptest.NewTLSServer(app.routes())

	testServer.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	tests := []struct {
		name        string
		petType     string
		urlPath     string
		wantCode    int
		wantFormTag string
	}{
		// {
		// 	name:        "Valid pet type",
		// 	petType:     "Cat",
		// 	urlPath:     "/pet/add-pet-type",
		// 	wantCode:    http.StatusSeeOther,
		// 	wantFormTag: "",
		// },
		// {
		// 	name:        "Duplicate pet type",
		// 	petType:     "Dog",
		// 	urlPath:     "/pet/add-pet-type",
		// 	wantCode:    http.StatusUnprocessableEntity,
		// 	wantFormTag: `<form hx-post="/pet/add-pet-type" hx-ext='json-enc' novalidate>`,
		// },
		{
			name:        "Empty pet type",
			petType:     "",
			urlPath:     "/pet/add-pet-type",
			wantCode:    http.StatusUnprocessableEntity,
			wantFormTag: `<form hx-post="/pet/add-pet-type" hx-ext='json-enc' novalidate>`,
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
			matches := errorHtmlRX.FindStringSubmatch(string(body))

			for _, match := range matches {
				fmt.Println(match)
			}

			if rs.StatusCode != test.wantCode {
				t.Errorf("got: %v; want: %v", rs.StatusCode, test.wantCode)
			}

			if !strings.Contains(string(body), test.wantFormTag) {
				t.Errorf("got: %q; expected to contain: %q", string(body), test.wantFormTag)
			}
		})
	}
}
