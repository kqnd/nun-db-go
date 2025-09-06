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
	client.Set("teste", "123")
	client.Set("name", "alex fernando")

	teste, err := client.Get("teste")
	if err != nil {
		fmt.Println("occurred a error: ", err)
	}

	name, err := client.Get("name")
	if err != nil {
		fmt.Println("occurred a error: ", err)
	}

	fmt.Println("[teste]: ", teste)
	fmt.Println("[name]: ", name)
	select {}
}
