package main

import (
	"fmt"

	"github.com/nundb/nundb"
)

func main() {
	nundb_client := nundb.NewClient("ws://localhost:3012/", "user", "pwd")

	fmt.Println(nundb_client)
}
