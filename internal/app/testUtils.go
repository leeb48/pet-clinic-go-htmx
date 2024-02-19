package app

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"pet-clinic.bonglee.com/internal/models/mocks"
)

func NewTestApp(t *testing.T) *App {
	templateCache, err := CreateTemplateCache()
	if err != nil {
		t.Fatal(err)
	}

	session := scs.New()
	session.Lifetime = 24 * time.Hour

	app := &App{
		Logger:        slog.New(slog.NewTextHandler(io.Discard, nil)),
		TemplateCache: templateCache,
		Session:       session,
		Owners:        &mocks.OwnerModel{},
		PetTypes:      &mocks.PetTypeModel{},
		Pets:          &mocks.PetModel{},
		Vets:          &mocks.VetModel{},
	}

	return app
}

type TestServer struct {
	*httptest.Server
}

func NewTestServer(t *testing.T, routes http.Handler) *TestServer {
	ts := httptest.NewTLSServer(routes)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	ts.Client().Jar = jar

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &TestServer{ts}
}

func (ts *TestServer) PostReq(t *testing.T, urlPath string, json []byte) (int, http.Header, string) {
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

func (ts *TestServer) GetReq(t *testing.T, urlPath string) (int, http.Header, string) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
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

func (ts *TestServer) PutReq(t *testing.T, urlPath string, json []byte) (int, http.Header, string) {
	req, err := http.NewRequest(http.MethodPut, ts.URL+urlPath, bytes.NewBuffer(json))
	if err != nil {
		t.Fatal(err)
	}

	rs, err := ts.Client().Do(req)
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

func (ts *TestServer) DeleteReq(t *testing.T, urlPath string) (int, http.Header, string) {
	req, err := http.NewRequest(http.MethodDelete, ts.URL+urlPath, nil)
	if err != nil {
		t.Fatal(err)
	}

	rs, err := ts.Client().Do(req)
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

var GetFormTextDangerHtml = regexp.MustCompile(`<div class="form-text text-danger">(.*) `)
