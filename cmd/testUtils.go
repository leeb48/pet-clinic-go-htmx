package main

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"pet-clinic.bonglee.com/internal/models/mocks"
)

func newTestApp(t *testing.T) *application {
	templateCache, err := createTemplateCache()
	if err != nil {
		t.Fatal(err)
	}

	session := scs.New()
	session.Lifetime = 24 * time.Hour

	app := &application{
		logger:        slog.New(slog.NewTextHandler(io.Discard, nil)),
		templateCache: templateCache,
		session:       session,
		owners:        &mocks.OwnerModel{},
		petTypes:      &mocks.PetTypeModel{},
		pets:          &mocks.PetModel{},
	}

	return app
}

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, routes http.Handler) *testServer {
	ts := httptest.NewTLSServer(routes)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	ts.Client().Jar = jar

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

func (ts *testServer) postReq(t *testing.T, urlPath string, json []byte) (int, http.Header, string) {
	rs, err := ts.Client().Post(ts.URL+urlPath, "application/json", bytes.NewReader(json))
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	return rs.StatusCode, rs.Header, string(body)
}
