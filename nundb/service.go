package nundb

import (
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

type Client struct {
	Name       string
	connection *websocket.Conn
}

func NewClient(serverUrl, name, password string) *Client {
	u, err := url.Parse(serverUrl)
	if err != nil {
		log.Fatal("error parsing url: ", err)
	}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("error during connection: ", err)
	}
	defer c.Close()

	fmt.Println("connected to: ", u.String())

	return &Client{
		Name:       name,
		connection: c,
	}
}
