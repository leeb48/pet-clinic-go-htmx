package models

import (
	"database/sql"
	"os"
	"testing"
)

func newTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("mysql", "test_app:1234@/test_petClinic?parseTime=true&multiStatements=true")
	if err != nil {
		t.Fatal(err)
	}

	script, err := os.ReadFile("./testdata/setup.sql")
	if err != nil {
		db.Close()
		t.Fatal(err)
	}

	_, err = db.Exec(string(script))
	if err != nil {
		executeTeardown(t, db)
		db.Close()
		t.Fatal(err)
	}

	t.Cleanup(func() {
		defer db.Close()
		executeTeardown(t, db)

	})

	return db
}

func executeTeardown(t *testing.T, db *sql.DB) {
	script, err := os.ReadFile("./testdata/teardown.sql")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
	}
}
