package nundb

import (
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

type Client struct {
	Username   string
	connection *websocket.Conn
}

func NewClient(serverUrl, username, password string) *Client {
	u, err := url.Parse(serverUrl)
	if err != nil {
		log.Fatal("error parsing url: ", err)
	}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("error during connection: ", err)
	}
	defer c.Close()

	client := &Client{
		Username:   username,
		connection: c,
	}

	fmt.Println("connected to: ", u.String())
	authMsg := fmt.Sprintf("auth %s %s", username, password)
	client.SendCommand(authMsg)

	return client
}

func (c *Client) SendCommand(command string) {
	err := c.connection.WriteMessage(websocket.TextMessage, []byte(command))
	if err != nil {
		log.Fatalf("error during command: %s - error: %s", command, err)
	}
	fmt.Printf("[command] %s executed\n", command)
}
