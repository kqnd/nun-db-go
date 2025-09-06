package nundb

import (
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/nundb/v1/nundb/response"
)

type Client struct {
	Username string
	Database string

	// private
	responseHandler response.Handler
	conn            *websocket.Conn
	watchers        map[string][]func(interface{})
	pendings        map[string]chan interface{}
}

func NewClient(serverUrl, username, password string) (*Client, error) {
	client := &Client{
		Username: username,
		watchers: make(map[string][]func(interface{})),
		pendings: make(map[string]chan interface{}),
	}

	u, err := url.Parse(serverUrl)
	if err != nil {
		return nil, err
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}

	client.conn = conn

	fmt.Println("connected to: ", u.String())
	authMsg := fmt.Sprintf("auth %s %s", username, password)
	client.SendCommand(authMsg)

	go client.listen()

	return client, nil
}

func (c *Client) listen() {
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("conn closed: ", err)
			return
		}
		fmt.Println(msg)
		// response.Handler(string(msg))
	}
}

func (c *Client) CreateDatabase(name, pwd string) {
	c.SendCommand(fmt.Sprintf("create-db %s %s", name, pwd))
}

func (c *Client) UseDatabase(name, pwd string) {
	c.SendCommand(fmt.Sprintf("use-db %s %s", name, pwd))
}

func (c *Client) SendCommand(command string) {
	err := c.conn.WriteMessage(websocket.TextMessage, []byte(command))
	if err != nil {
		log.Fatalf("error during command: %s - error: %s", command, err)
	}
	fmt.Printf("[command] %s executed\n", command)
}
