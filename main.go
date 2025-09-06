package main

import (
	"fmt"
	"log"

	"github.com/nundb/v1/nundb"
)

func main() {
	client, err := nundb.NewClient("ws://localhost:3012/", "user-name", "user-pwd")
	if err != nil {
		log.Fatal("error occurred: ", err)
	}

	client.UseDatabase("oi", "oi")

	fmt.Println(client)
	select {}

}
