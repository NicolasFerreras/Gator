package config

import (
	"testing"
)

// isolateHome apunta el directorio home (HOME/USERPROFILE) a un temp,
// para no leer/escribir ~/.gatorconfig.json real.
func isolateHome(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	t.Setenv("USERPROFILE", dir)
}

func TestSetUserThenRead_RoundTrip(t *testing.T) {
	// Aislamos el home para no tocar ~/.gatorconfig.json real.
	isolateHome(t)

	cfg := Config{DBURL: "postgres://test", CurrentUserName: ""}
	if err := cfg.SetUser("nico"); err != nil {
		t.Fatalf("SetUser: %v", err)
	}

	got, err := Read()
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if got.CurrentUserName != "nico" {
		t.Fatalf("CurrentUserName = %q, want nico", got.CurrentUserName)
	}
	if got.DBURL != "postgres://test" {
		t.Fatalf("DBURL = %q, want postgres://test", got.DBURL)
	}
}

func TestRead_MissingFile_ReturnsError(t *testing.T) {
	isolateHome(t)

	_, err := Read()
	if err == nil {
		t.Fatal("se esperaba error al leer archivo de config inexistente")
	}
}

func TestSetUser_OverwritesCurrentUser(t *testing.T) {
	isolateHome(t)

	cfg := Config{}
	if err := cfg.SetUser("primero"); err != nil {
		t.Fatalf("SetUser 1: %v", err)
	}
	if err := cfg.SetUser("segundo"); err != nil {
		t.Fatalf("SetUser 2: %v", err)
	}

	got, err := Read()
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if got.CurrentUserName != "segundo" {
		t.Fatalf("CurrentUserName = %q, want segundo", got.CurrentUserName)
	}
}
