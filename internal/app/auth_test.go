package app

import (
	"encoding/base64"
	"testing"
)

func TestCheckAuth(t *testing.T) {
	key := "secret"
	header := "Basic " + base64.StdEncoding.EncodeToString([]byte(key))
	if !checkAuth(header, key) {
		t.Fatalf("expected auth to succeed")
	}
	if checkAuth(header, "wrong") {
		t.Fatalf("expected auth to fail with wrong key")
	}
	if checkAuth("", key) {
		t.Fatalf("expected auth to fail with empty header")
	}
}
