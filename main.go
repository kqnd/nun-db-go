package main

import (
	"fmt"
	"log"

	nundb "github.com/viewfromaside/nun-db-go"
)

func main() {
	client, err := nundb.NewClient("ws://localhost:3012/", "user-name", "user-pwd")
	if err != nil {
		log.Fatal("error occurred: ", err)
	}

	client.UseDatabase("oi", "oi")
	client.Set("teste", "123")
	client.Set("name", "alex fernando")

	client.Watch("teste", func(data interface{}) {
		fmt.Printf("[teste] changed to %s\n", data)
	})

	teste, err := client.Get("teste")
	if err != nil {
		fmt.Println("occurred a error: ", err)
	}

	name, err := client.Get("name")
	if err != nil {
		fmt.Println("occurred a error: ", err)
	}

	client.Set("teste", "1")
	client.Set("teste", "3")
	client.Set("teste", "2")
	client.Set("teste", "2")
	client.Set("teste", "5")

	fmt.Println("[teste]: ", teste)
	fmt.Println("[name]: ", name)
	select {}
}
