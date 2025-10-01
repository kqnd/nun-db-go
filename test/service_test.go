package test

import (
	"testing"

	nundb "github.com/kqnd/nun-db-go"
)

func TestNewClient(t *testing.T) {
	_, err := nundb.NewClient("ws://localhost:3012", "user-name", "user-pwd")
	if err != nil {
		t.Fatalf("client could not be nil, err: %s", err)
	}
}
