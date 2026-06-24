package db

import (
	"testing"
)

func TestConnect_ValidDriver(t *testing.T) {
	db, err := Connect("postgres://fake:fake@localhost:5432/fake_db?sslmode=disable")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if db == nil {
		t.Error("expected non-nil db object")
	} else {
		_ = db.Close()
	}
}
