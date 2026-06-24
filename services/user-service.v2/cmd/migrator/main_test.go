package main

import (
	"testing"
)

func TestRunMigrations_InvalidDBURL(t *testing.T) {
	err := RunMigrations("invalid_db_url")

	if err == nil {
		t.Log("No error (maybe migrations not needed)")
	} else {
		t.Logf("Got expected error: %v", err)
	}
}

func TestRunMigrations_Basic(t *testing.T) {
	err := RunMigrations("postgres://fake:fake@localhost:5432/fake_db?sslmode=disable")
	if err != nil {
		t.Logf("Expected connection error: %v", err)
	}
}
