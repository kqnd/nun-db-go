package nundbgo

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/viewfromaside/nun-db-go/response"
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

	strictQueueWatchers bool
}

func NewClient(serverUrl, username, password string) (*Client, error) {
	watchers := make(map[string][]func(interface{}))
	pendings := make([]chan interface{}, 0)

	client := &Client{
		Username:            username,
		watchers:            watchers,
		pendings:            pendings,
		strictQueueWatchers: false,
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

func (c *Client) SetWatchersQueueMode(strict bool) {
	c.strictQueueWatchers = strict
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
		c.responseHandler.WatchingValues(c.strictQueueWatchers)

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

func (c *Client) Set(key string, value interface{}) error {
	var strValue string
	switch v := value.(type) {
	case string:
		strValue = v
	default:
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("error on parsing value to json: %w", err)
		}
		strValue = string(jsonBytes)
	}
	c.SendCommand(fmt.Sprintf("set %s %s", key, strValue))
	return nil
}

func (c *Client) Increment(key string, value int) {
	c.SendCommand(fmt.Sprintf("increment %s %d", key, value))
}

func (c *Client) Remove(key string) {
	c.SendCommand("remove " + key)
}

func (c *Client) RemoveAllWatchers() {
	c.SendCommand("unwatch-all")
	c.watchers = make(map[string][]func(interface{}))
}

func (c *Client) RemoveWatcher(key string) {
	c.SendCommand("unwatch " + key)
	delete(c.watchers, key)
}

func (c *Client) Watch(key string, cb func(interface{})) {
	c.queue.Lock()
	c.watchers[key] = append(c.watchers[key], cb)
	c.queue.Unlock()
	c.SendCommand("watch " + key)
}

func (c *Client) Get(key string) (interface{}, error) {
	ch := make(chan interface{})

	c.queue.Lock()
	c.pendings = append(c.pendings, ch)
	c.queue.Unlock()

	c.SendCommand("get " + key)
	return <-ch, nil
}
