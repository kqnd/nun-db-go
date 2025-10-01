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

func TestGettingValues(t *testing.T) {
	client, err := nundb.NewClient("ws://localhost:3012", "user-name", "user-pwd")
	if err != nil {
		t.Fatalf("client could not be nil, err: %s", err)
	}

	client.Set("foo", "bar")
	foo, err := client.Get("foo")
	if err != nil {
		t.Fatalf("foo could not be nil, err: %s", err)
	}

	if foo != "bar" {
		t.Fatalf("foo value is different that bar")
	}
}
