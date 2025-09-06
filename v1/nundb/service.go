package nundb

import (
	"fmt"
	"log"
	"net/url"
	"sync"

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
	pendings        []chan interface{}
	queue           sync.Mutex
}

func NewClient(serverUrl, username, password string) (*Client, error) {
	watchers := make(map[string][]func(interface{}))
	pendings := make([]chan interface{}, 0)

	client := &Client{
		Username: username,
		watchers: watchers,
		pendings: pendings,
	}

	u, err := url.Parse(serverUrl)
	if err != nil {
		return nil, err
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}

	handler := response.CreateHandler(&client.watchers, &client.pendings)
	client.responseHandler = handler
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

		c.responseHandler.SetPayload(string(msg))
		c.responseHandler.GettingValues()

	}
}

func (c *Client) SendCommand(command string) {
	err := c.conn.WriteMessage(websocket.TextMessage, []byte(command))
	if err != nil {
		log.Fatalf("error during command: %s - error: %s", command, err)
	}
	fmt.Printf("[command] %s executed\n", command)
}

func (c *Client) CreateDatabase(name, pwd string) {
	c.SendCommand(fmt.Sprintf("create-db %s %s", name, pwd))
}

func (c *Client) UseDatabase(name, pwd string) {
	c.SendCommand(fmt.Sprintf("use-db %s %s", name, pwd))
}

func (c *Client) Set(key, value string) {
	c.SendCommand(fmt.Sprintf("set %s %s", key, value))
}

func (c *Client) Get(key string) (interface{}, error) {
	ch := make(chan interface{})

	c.queue.Lock()
	c.pendings = append(c.pendings, ch)
	c.queue.Unlock()

	c.SendCommand("get " + key)
	return <-ch, nil
}
